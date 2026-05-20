package planning

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/jmt-labs/lodestone/internal/lodestone/schema"
)

const (
	DefaultModel = "claude-opus-4-7"
)

type Engine struct {
	runner Runner
	model  string
	now    func() time.Time
}

type Option func(*Engine)

func WithRunner(r Runner) Option         { return func(e *Engine) { e.runner = r } }
func WithModel(m string) Option          { return func(e *Engine) { e.model = m } }
func WithNow(fn func() time.Time) Option { return func(e *Engine) { e.now = fn } }

func New(opts ...Option) *Engine {
	e := &Engine{
		runner: NewClaudeRunner(),
		model:  DefaultModel,
		now:    func() time.Time { return time.Now().UTC() },
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

type Result struct {
	SpecPath string
	PlanPath string
	Spec     string
	Plan     string
	Prompt   string
	Model    string
}

func (e *Engine) Plan(ctx context.Context, fp schema.Fingerprint, rec schema.Recommendation) (*Result, error) {
	prompt, err := BuildPrompt(fp, rec)
	if err != nil {
		return nil, err
	}
	out, err := e.runner.Run(ctx, e.model, prompt)
	if err != nil {
		return nil, err
	}
	specMD, planMD, err := SplitResponse(out)
	if err != nil {
		return nil, fmt.Errorf("parse claude output: %w (raw output stored in error)", err)
	}
	slug := slugFromRec(rec)
	date := e.now().Format("2006-01-02")
	return &Result{
		Spec:     specMD,
		Plan:     planMD,
		Prompt:   prompt,
		Model:    e.model,
		SpecPath: filepath.Join("docs", "superpowers", "specs", fmt.Sprintf("%s-%s-design.md", date, slug)),
		PlanPath: filepath.Join("docs", "superpowers", "plans", fmt.Sprintf("%s-%s.md", date, slug)),
	}, nil
}

func (r *Result) Persist(repoRoot string) error {
	for _, item := range []struct {
		rel     string
		content string
	}{
		{r.SpecPath, r.Spec},
		{r.PlanPath, r.Plan},
	} {
		full := filepath.Join(repoRoot, item.rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", filepath.Dir(full), err)
		}
		if err := os.WriteFile(full, []byte(item.content+"\n"), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", full, err)
		}
	}
	return nil
}

var slugifyRE = regexp.MustCompile(`[^a-z0-9]+`)

func slugFromRec(rec schema.Recommendation) string {
	base := rec.SignalID
	if base == "" {
		base = rec.ID
	}
	base = strings.ToLower(base)
	base = strings.TrimPrefix(base, "sha256:")
	if len(base) > 16 {
		base = base[:16]
	}
	base = slugifyRE.ReplaceAllString(base, "-")
	base = strings.Trim(base, "-")
	if base == "" {
		base = "lodestone-plan"
	}
	return "lodestone-" + base
}
