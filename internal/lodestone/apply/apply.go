package apply

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jmt-labs/lodestone/internal/lodestone/planning"
	"github.com/jmt-labs/lodestone/internal/lodestone/schema"
)

type GitRunner interface {
	Status(ctx context.Context, repoRoot string) (string, error)
	CreateBranch(ctx context.Context, repoRoot, branch string) error
	Add(ctx context.Context, repoRoot string, paths ...string) error
	Commit(ctx context.Context, repoRoot, message string) error
	Push(ctx context.Context, repoRoot, branch string) error
	DeleteBranchLocal(ctx context.Context, repoRoot, branch string) error
	DeleteBranchRemote(ctx context.Context, repoRoot, branch string) error
}

type PRRunner interface {
	CreateDraftPR(ctx context.Context, repoRoot, branch, title, body string) (number int, url string, err error)
	ClosePR(ctx context.Context, repoRoot string, number int) error
}

type Engine struct {
	repoRoot string
	state    *State
	planning *planning.Engine
	git      GitRunner
	pr       PRRunner
	now      func() time.Time
}

type Option func(*Engine)

func WithGit(r GitRunner) Option { return func(e *Engine) { e.git = r } }
func WithPR(r PRRunner) Option   { return func(e *Engine) { e.pr = r } }
func WithPlanning(p *planning.Engine) Option {
	return func(e *Engine) { e.planning = p }
}
func WithNow(fn func() time.Time) Option { return func(e *Engine) { e.now = fn } }

func New(repoRoot string, state *State, opts ...Option) *Engine {
	e := &Engine{
		repoRoot: repoRoot,
		state:    state,
		planning: planning.New(),
		git:      RealGit{},
		pr:       RealPR{},
		now:      func() time.Time { return time.Now().UTC() },
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

type Result struct {
	Apply Apply
	Spec  string
	Plan  string
}

func (e *Engine) Apply(ctx context.Context, fp schema.Fingerprint, rec schema.Recommendation) (*Result, error) {
	if g := CheckRecommendation(rec); !g.Passed {
		return nil, fmt.Errorf("recommendation safety gates failed: %v", g.Violations)
	}
	list, err := e.state.List()
	if err != nil {
		return nil, fmt.Errorf("read applies: %w", err)
	}
	if g := CheckRateLimit(list, e.now()); !g.Passed {
		return nil, fmt.Errorf("rate-limit gate failed: %v", g.Violations)
	}
	status, err := e.git.Status(ctx, e.repoRoot)
	if err != nil {
		return nil, fmt.Errorf("git status: %w", err)
	}
	if g := CheckCleanGit(status); !g.Passed {
		return nil, fmt.Errorf("git clean gate failed: %v", g.Violations)
	}

	planRes, err := e.planning.Plan(ctx, fp, rec)
	if err != nil {
		return nil, fmt.Errorf("plan: %w", err)
	}
	if err := planRes.Persist(e.repoRoot); err != nil {
		return nil, fmt.Errorf("persist plan: %w", err)
	}

	branch := branchName(rec, e.now())
	if err := e.git.CreateBranch(ctx, e.repoRoot, branch); err != nil {
		return nil, fmt.Errorf("create branch: %w", err)
	}
	if err := e.git.Add(ctx, e.repoRoot, planRes.SpecPath, planRes.PlanPath); err != nil {
		return nil, fmt.Errorf("git add: %w", err)
	}
	commitMsg := fmt.Sprintf("chore(lodestone): apply recommendation %s\n\nspec: %s\nplan: %s", rec.ID, planRes.SpecPath, planRes.PlanPath)
	if err := e.git.Commit(ctx, e.repoRoot, commitMsg); err != nil {
		return nil, fmt.Errorf("git commit: %w", err)
	}
	if err := e.git.Push(ctx, e.repoRoot, branch); err != nil {
		return nil, fmt.Errorf("git push: %w", err)
	}

	prTitle := fmt.Sprintf("Lodestone Apply: %s", recShortTitle(rec))
	prBody := fmt.Sprintf("Auto-PR aus lodestone apply für Recommendation `%s`.\n\n- compatibility: %.2f\n- effort: %s\n- risk: %s\n\nSpec: %s\nPlan: %s\n\n**Hinweis:** Dieser PR ist als Draft erstellt. Bitte review die Spec/Plan-Inhalte vor Merge.", rec.ID, rec.Compatibility, rec.Effort, rec.Risk, planRes.SpecPath, planRes.PlanPath)

	number, url, err := e.pr.CreateDraftPR(ctx, e.repoRoot, branch, prTitle, prBody)
	if err != nil {
		number = 0
		url = ""
	}

	apply := Apply{
		RecID:     rec.ID,
		Branch:    branch,
		PRNumber:  number,
		PRURL:     url,
		Status:    statusFromPR(number, err),
		AppliedAt: e.now(),
	}
	if err := e.state.Append(apply); err != nil {
		return nil, fmt.Errorf("append state: %w", err)
	}
	return &Result{Apply: apply, Spec: planRes.Spec, Plan: planRes.Plan}, nil
}

func (e *Engine) Undo(ctx context.Context, identifier string) (*Apply, error) {
	apply, ok, err := e.state.FindBy(identifier)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("apply %q not found", identifier)
	}
	if apply.PRNumber != 0 {
		_ = e.pr.ClosePR(ctx, e.repoRoot, apply.PRNumber)
	}
	_ = e.git.DeleteBranchLocal(ctx, e.repoRoot, apply.Branch)
	_ = e.git.DeleteBranchRemote(ctx, e.repoRoot, apply.Branch)
	apply.Status = "undone"
	if err := e.state.Replace(apply); err != nil {
		return nil, fmt.Errorf("replace state: %w", err)
	}
	return &apply, nil
}

func statusFromPR(number int, err error) string {
	if err != nil {
		return "branch_pushed_no_pr"
	}
	if number > 0 {
		return "draft_open"
	}
	return "branch_pushed_no_pr"
}

func branchName(rec schema.Recommendation, now time.Time) string {
	suffix := strings.TrimPrefix(rec.ID, "sha256:")
	if len(suffix) > 8 {
		suffix = suffix[:8]
	}
	if suffix == "" {
		sum := sha1.Sum([]byte(rec.SignalID))
		suffix = hex.EncodeToString(sum[:4])
	}
	return fmt.Sprintf("lodestone/apply-%s-%s", suffix, now.Format("2006-01-02"))
}

func recShortTitle(rec schema.Recommendation) string {
	stripped := strings.TrimPrefix(rec.SignalID, "sha256:")
	if len(stripped) > 8 {
		return stripped[:8]
	}
	if stripped == "" {
		return rec.ID
	}
	return stripped
}

func WriteFileUnderRoot(repoRoot, rel, content string) error {
	full := filepath.Join(repoRoot, rel)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return err
	}
	return os.WriteFile(full, []byte(content), 0o644)
}
