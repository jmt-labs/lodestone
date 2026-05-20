package planning

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jmt-labs/lodestone/internal/lodestone/schema"
)

func sampleFP() schema.Fingerprint {
	return schema.Fingerprint{
		SchemaVersion: schema.FingerprintSchemaVersion,
		GeneratedAt:   time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC),
		Languages:     []string{"Go"},
		Frameworks:    []string{"cobra"},
	}
}

func sampleRec() schema.Recommendation {
	return schema.Recommendation{
		SchemaVersion: schema.RecommendationSchemaVersion,
		ID:            "sha256:abc123",
		SignalID:      "sha256:sig-42",
		Compatibility: 0.85,
		Effort:        schema.EffortS,
		Risk:          schema.RiskLow,
	}
}

const happyOutput = `Hier ist mein Vorschlag:

===SPEC===
# Spec — Anbindung an X

Tradeoff: build vs. buy.
===PLAN===
# Plan — Anbindung an X

- [ ] T1: Integrationstest schreiben
- [ ] T2: Adapter implementieren
===END===

Viel Erfolg!`

func TestEnginePlanHappyPath(t *testing.T) {
	now := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	fake := &FakeRunner{Output: happyOutput}
	e := New(
		WithRunner(fake),
		WithModel("claude-opus-test"),
		WithNow(func() time.Time { return now }),
	)

	res, err := e.Plan(context.Background(), sampleFP(), sampleRec())
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	if !strings.Contains(res.Spec, "Anbindung an X") {
		t.Errorf("Spec content unexpected: %q", res.Spec)
	}
	if !strings.Contains(res.Plan, "- [ ] T1") {
		t.Errorf("Plan content unexpected: %q", res.Plan)
	}
	if res.Model != "claude-opus-test" {
		t.Errorf("Model = %q", res.Model)
	}
	wantSpec := filepath.Join("docs", "superpowers", "specs", "2026-05-20-lodestone-sig-42-design.md")
	if res.SpecPath != wantSpec {
		t.Errorf("SpecPath = %q, want %q", res.SpecPath, wantSpec)
	}
	if len(fake.Calls) != 1 {
		t.Fatalf("runner called %d times, want 1", len(fake.Calls))
	}
	if fake.Calls[0].Model != "claude-opus-test" {
		t.Errorf("runner model = %q", fake.Calls[0].Model)
	}
	if !strings.Contains(fake.Calls[0].Prompt, "===SPEC===") {
		t.Errorf("prompt missing marker instructions")
	}
}

func TestEnginePlanPersists(t *testing.T) {
	root := t.TempDir()
	now := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	fake := &FakeRunner{Output: happyOutput}
	e := New(
		WithRunner(fake),
		WithNow(func() time.Time { return now }),
	)
	res, err := e.Plan(context.Background(), sampleFP(), sampleRec())
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	if err := res.Persist(root); err != nil {
		t.Fatalf("Persist: %v", err)
	}
	for _, rel := range []string{res.SpecPath, res.PlanPath} {
		full := filepath.Join(root, rel)
		if _, err := os.Stat(full); err != nil {
			t.Errorf("file missing: %s (%v)", full, err)
		}
	}
}

func TestEnginePlanRunnerError(t *testing.T) {
	fake := &FakeRunner{Err: errors.New("boom")}
	e := New(WithRunner(fake))
	if _, err := e.Plan(context.Background(), sampleFP(), sampleRec()); err == nil {
		t.Fatal("expected error")
	}
}

func TestEnginePlanMissingMarkers(t *testing.T) {
	fake := &FakeRunner{Output: "no markers here"}
	e := New(WithRunner(fake))
	if _, err := e.Plan(context.Background(), sampleFP(), sampleRec()); err == nil {
		t.Fatal("expected error on missing markers")
	}
}

func TestSplitResponse(t *testing.T) {
	spec, plan, err := SplitResponse(happyOutput)
	if err != nil {
		t.Fatalf("SplitResponse: %v", err)
	}
	if !strings.HasPrefix(spec, "# Spec") {
		t.Errorf("Spec = %q", spec)
	}
	if !strings.HasPrefix(plan, "# Plan") {
		t.Errorf("Plan = %q", plan)
	}
}

func TestBuildPromptEncodesJSON(t *testing.T) {
	p, err := BuildPrompt(sampleFP(), sampleRec())
	if err != nil {
		t.Fatalf("BuildPrompt: %v", err)
	}
	if !strings.Contains(p, "cobra") {
		t.Errorf("prompt missing fingerprint content")
	}
	if !strings.Contains(p, "sha256:sig-42") {
		t.Errorf("prompt missing recommendation content")
	}
}

func TestSlugFromRec(t *testing.T) {
	got := slugFromRec(sampleRec())
	if !strings.HasPrefix(got, "lodestone-") {
		t.Errorf("slug = %q", got)
	}
	if strings.Contains(got, ":") {
		t.Errorf("slug contains colon: %q", got)
	}
}
