package apply

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/jmt-labs/lodestone/internal/lodestone/planning"
	"github.com/jmt-labs/lodestone/internal/lodestone/schema"
)

func TestCheckRecommendationGates(t *testing.T) {
	good := schema.Recommendation{
		Risk:          schema.RiskLow,
		Effort:        schema.EffortXS,
		Compatibility: 0.91,
	}
	if r := CheckRecommendation(good); !r.Passed {
		t.Errorf("good rec failed gates: %v", r.Violations)
	}

	cases := []schema.Recommendation{
		{Risk: schema.RiskMed, Effort: schema.EffortXS, Compatibility: 0.9},
		{Risk: schema.RiskLow, Effort: schema.EffortS, Compatibility: 0.9},
		{Risk: schema.RiskLow, Effort: schema.EffortXS, Compatibility: 0.5},
	}
	for i, rec := range cases {
		if r := CheckRecommendation(rec); r.Passed {
			t.Errorf("case %d should fail, got pass", i)
		}
	}
}

func TestCheckRateLimit(t *testing.T) {
	now := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	applies := []Apply{{RecID: "x", AppliedAt: now.Add(-1 * time.Hour)}}
	if r := CheckRateLimit(applies, now); r.Passed {
		t.Errorf("expected rate-limit violation, got pass")
	}
	applies = []Apply{{RecID: "x", AppliedAt: now.Add(-48 * time.Hour)}}
	if r := CheckRateLimit(applies, now); !r.Passed {
		t.Errorf("expected pass (>24h ago), got: %v", r.Violations)
	}
	if r := CheckRateLimit(nil, now); !r.Passed {
		t.Errorf("empty list should pass")
	}
}

func TestCheckCleanGit(t *testing.T) {
	if r := CheckCleanGit(""); !r.Passed {
		t.Errorf("empty status should pass")
	}
	if r := CheckCleanGit(" M file.go"); r.Passed {
		t.Errorf("dirty status should fail")
	}
}

func TestStateAppendAndList(t *testing.T) {
	dir := t.TempDir()
	s, err := NewState(dir)
	if err != nil {
		t.Fatalf("NewState: %v", err)
	}
	now := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	a := Apply{RecID: "rec-1", Branch: "lodestone/apply-x-2026-05-20", Status: "draft_open", AppliedAt: now}
	if err := s.Append(a); err != nil {
		t.Fatalf("Append: %v", err)
	}
	got, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(got) != 1 || got[0].RecID != "rec-1" {
		t.Errorf("List = %+v", got)
	}
}

func TestStateFindBy(t *testing.T) {
	dir := t.TempDir()
	s, _ := NewState(dir)
	_ = s.Append(Apply{RecID: "r", Branch: "b", Status: "draft_open"})
	if _, ok, _ := s.FindBy("b"); !ok {
		t.Errorf("FindBy by branch failed")
	}
	if _, ok, _ := s.FindBy("r"); !ok {
		t.Errorf("FindBy by rec id failed")
	}
	if _, ok, _ := s.FindBy("missing"); ok {
		t.Errorf("FindBy unknown should return false")
	}
}

type fakeGit struct {
	statusOut string
	calls     []string
}

func (g *fakeGit) Status(_ context.Context, _ string) (string, error) {
	g.calls = append(g.calls, "status")
	return g.statusOut, nil
}
func (g *fakeGit) CreateBranch(_ context.Context, _, b string) error {
	g.calls = append(g.calls, "branch:"+b)
	return nil
}
func (g *fakeGit) Add(_ context.Context, _ string, paths ...string) error {
	for _, p := range paths {
		g.calls = append(g.calls, "add:"+p)
	}
	return nil
}
func (g *fakeGit) Commit(_ context.Context, _, msg string) error {
	g.calls = append(g.calls, "commit:"+msg[:10])
	return nil
}
func (g *fakeGit) Push(_ context.Context, _, b string) error {
	g.calls = append(g.calls, "push:"+b)
	return nil
}
func (g *fakeGit) DeleteBranchLocal(_ context.Context, _, b string) error {
	g.calls = append(g.calls, "del-local:"+b)
	return nil
}
func (g *fakeGit) DeleteBranchRemote(_ context.Context, _, b string) error {
	g.calls = append(g.calls, "del-remote:"+b)
	return nil
}

type fakePR struct {
	number int
	url    string
	closed []int
}

func (p *fakePR) CreateDraftPR(_ context.Context, _, _, _, _ string) (int, string, error) {
	if p.number == 0 {
		p.number = 42
		p.url = "https://github.com/x/y/pull/42"
	}
	return p.number, p.url, nil
}
func (p *fakePR) ClosePR(_ context.Context, _ string, n int) error {
	p.closed = append(p.closed, n)
	return nil
}

const fakePlannerOutput = `===SPEC===
# Spec für Auto-Apply

Test-Spec.
===PLAN===
# Plan

- [ ] T1: noop
===END===`

func TestApplyHappyPath(t *testing.T) {
	repo := t.TempDir()
	state, _ := NewState(filepath.Join(repo, ".lodestone"))
	git := &fakeGit{}
	pr := &fakePR{}
	now := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	planner := planning.New(
		planning.WithRunner(&planning.FakeRunner{Output: fakePlannerOutput}),
		planning.WithNow(func() time.Time { return now }),
	)
	eng := New(repo, state,
		WithGit(git),
		WithPR(pr),
		WithPlanning(planner),
		WithNow(func() time.Time { return now }),
	)

	rec := schema.Recommendation{
		ID:            "sha256:abcdef1234",
		SignalID:      "sha256:sig-1",
		Risk:          schema.RiskLow,
		Effort:        schema.EffortXS,
		Compatibility: 0.9,
	}
	fp := schema.Fingerprint{SchemaVersion: schema.FingerprintSchemaVersion}

	res, err := eng.Apply(context.Background(), fp, rec)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if res.Apply.PRNumber != 42 {
		t.Errorf("PRNumber = %d", res.Apply.PRNumber)
	}
	if res.Apply.Status != "draft_open" {
		t.Errorf("Status = %q", res.Apply.Status)
	}
	if !contains(git.calls, "branch:"+res.Apply.Branch) {
		t.Errorf("branch creation not in git calls: %v", git.calls)
	}
}

func TestApplyFailsGate(t *testing.T) {
	repo := t.TempDir()
	state, _ := NewState(filepath.Join(repo, ".lodestone"))
	eng := New(repo, state,
		WithGit(&fakeGit{}),
		WithPR(&fakePR{}),
		WithPlanning(planning.New(planning.WithRunner(&planning.FakeRunner{Output: fakePlannerOutput}))),
	)
	rec := schema.Recommendation{Risk: schema.RiskHigh, Effort: schema.EffortXS, Compatibility: 0.9}
	if _, err := eng.Apply(context.Background(), schema.Fingerprint{}, rec); err == nil {
		t.Fatal("expected gate failure")
	}
}

func TestApplyFailsRateLimit(t *testing.T) {
	repo := t.TempDir()
	state, _ := NewState(filepath.Join(repo, ".lodestone"))
	now := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	_ = state.Append(Apply{RecID: "earlier", Branch: "earlier", AppliedAt: now.Add(-1 * time.Hour)})
	eng := New(repo, state,
		WithGit(&fakeGit{}),
		WithPR(&fakePR{}),
		WithPlanning(planning.New(planning.WithRunner(&planning.FakeRunner{Output: fakePlannerOutput}))),
		WithNow(func() time.Time { return now }),
	)
	rec := schema.Recommendation{Risk: schema.RiskLow, Effort: schema.EffortXS, Compatibility: 0.95, ID: "sha256:later"}
	if _, err := eng.Apply(context.Background(), schema.Fingerprint{}, rec); err == nil {
		t.Fatal("expected rate-limit failure")
	}
}

func TestUndoClosesPRAndDeletesBranch(t *testing.T) {
	repo := t.TempDir()
	state, _ := NewState(filepath.Join(repo, ".lodestone"))
	_ = state.Append(Apply{RecID: "r1", Branch: "lodestone/apply-r1-2026", PRNumber: 7, Status: "draft_open"})
	git := &fakeGit{}
	pr := &fakePR{}
	eng := New(repo, state, WithGit(git), WithPR(pr))
	a, err := eng.Undo(context.Background(), "lodestone/apply-r1-2026")
	if err != nil {
		t.Fatalf("Undo: %v", err)
	}
	if a.Status != "undone" {
		t.Errorf("status = %q", a.Status)
	}
	if len(pr.closed) != 1 || pr.closed[0] != 7 {
		t.Errorf("expected PR 7 closed, got %v", pr.closed)
	}
	if !contains(git.calls, "del-local:lodestone/apply-r1-2026") {
		t.Errorf("local branch not deleted: %v", git.calls)
	}
}

func TestBranchNameDeterministic(t *testing.T) {
	now := time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC)
	rec := schema.Recommendation{ID: "sha256:0123456789abcdef"}
	got := branchName(rec, now)
	want := "lodestone/apply-01234567-2026-05-20"
	if got != want {
		t.Errorf("branchName = %q, want %q", got, want)
	}
}

func contains(list []string, item string) bool {
	for _, s := range list {
		if s == item {
			return true
		}
	}
	return false
}
