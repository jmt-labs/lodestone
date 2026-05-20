package scoring

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/jmt-labs/lodestone/internal/lodestone/schema"
)

func goRepoFingerprint() schema.Fingerprint {
	return schema.Fingerprint{
		SchemaVersion: schema.FingerprintSchemaVersion,
		GeneratedAt:   time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC),
		Languages:     []string{"Go"},
		Frameworks:    []string{"cobra"},
		Deps:          map[string]string{"github.com/spf13/cobra": "v1.10.2"},
	}
}

func mkSignal(id, lang string, stars int, tags []string, license string, lastCommit time.Time) schema.Signal {
	return schema.Signal{
		SchemaVersion: schema.SignalSchemaVersion,
		ID:            id,
		Source:        "test",
		URL:           "https://example.test/" + id,
		Title:         id,
		CapturedAt:    time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC),
		Language:      lang,
		Stars:         stars,
		TopicTags:     tags,
		License:       license,
		LastCommit:    lastCommit,
	}
}

func TestCompatibilityPerfectMatch(t *testing.T) {
	fp := goRepoFingerprint()
	sig := mkSignal("a", "Go", 1000, []string{"cobra"}, "mit", time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC))
	got := Compatibility(sig, fp)
	if got < 0.95 {
		t.Errorf("expected compatibility close to 1.0, got %v", got)
	}
}

func TestCompatibilityNoMatch(t *testing.T) {
	fp := goRepoFingerprint()
	sig := mkSignal("a", "Python", 1000, []string{"flask"}, "mit", time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC))
	got := Compatibility(sig, fp)
	if got != 0 {
		t.Errorf("expected 0, got %v", got)
	}
}

func TestCompatibilityMediumMatch(t *testing.T) {
	fp := goRepoFingerprint()
	sig := mkSignal("a", "Go", 200, []string{"web", "api"}, "mit", time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC))
	got := Compatibility(sig, fp)
	if got <= 0 || got >= 0.95 {
		t.Errorf("expected medium score (0,0.95), got %v", got)
	}
}

func TestEffortHeuristics(t *testing.T) {
	cases := []struct {
		name   string
		compat float64
		stars  int
		want   schema.EffortLevel
	}{
		{"no-match XL", 0.0, 1000, schema.EffortXL},
		{"match low-star S", 0.5, 50, schema.EffortS},
		{"match normal-star M", 0.5, 500, schema.EffortM},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			sig := mkSignal("s", "Go", c.stars, nil, "mit", time.Time{})
			if got := Effort(sig, c.compat); got != c.want {
				t.Errorf("Effort = %q, want %q", got, c.want)
			}
		})
	}
}

func TestRiskHeuristics(t *testing.T) {
	now := time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		name       string
		stars      int
		license    string
		lastCommit time.Time
		want       schema.RiskLevel
	}{
		{"all-good low", 800, "mit", now.AddDate(0, 0, -10), schema.RiskLow},
		{"no-license high", 800, "", now.AddDate(0, 0, -10), schema.RiskHigh},
		{"stale high", 800, "mit", now.AddDate(0, 0, -300), schema.RiskHigh},
		{"low-stars med", 100, "mit", now.AddDate(0, 0, -10), schema.RiskMed},
		{"mid-commit med", 800, "mit", now.AddDate(0, 0, -120), schema.RiskMed},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			sig := mkSignal("s", "Go", c.stars, nil, c.license, c.lastCommit)
			if got := Risk(sig, now); got != c.want {
				t.Errorf("Risk = %q, want %q", got, c.want)
			}
		})
	}
}

func TestScoreSortsByCompatibilityStarsID(t *testing.T) {
	fp := goRepoFingerprint()
	good := mkSignal("good", "Go", 800, []string{"cobra"}, "mit", time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC))
	mid := mkSignal("mid", "Go", 200, nil, "mit", time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC))
	bad := mkSignal("bad", "Python", 1000, []string{"flask"}, "mit", time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC))

	recs, err := Score(fp, []schema.Signal{bad, mid, good}, WithNow(func() time.Time { return time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC) }))
	if err != nil {
		t.Fatalf("Score: %v", err)
	}
	if len(recs) != 3 {
		t.Fatalf("want 3 recs, got %d", len(recs))
	}
	if recs[0].SignalID != "good" {
		t.Errorf("first = %s, want good", recs[0].SignalID)
	}
	if recs[1].SignalID != "mid" {
		t.Errorf("second = %s, want mid", recs[1].SignalID)
	}
	if recs[2].SignalID != "bad" {
		t.Errorf("last = %s, want bad", recs[2].SignalID)
	}
}

func TestScoreStarsTiebreak(t *testing.T) {
	fp := goRepoFingerprint()
	highStars := mkSignal("z-high", "Go", 5000, nil, "mit", time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC))
	lowStars := mkSignal("a-low", "Go", 100, nil, "mit", time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC))

	recs, err := Score(fp, []schema.Signal{lowStars, highStars})
	if err != nil {
		t.Fatalf("Score: %v", err)
	}
	if recs[0].SignalID != "z-high" {
		t.Errorf("tiebreak ignored stars: first = %s", recs[0].SignalID)
	}
}

func TestScoreDeterminism(t *testing.T) {
	fp := goRepoFingerprint()
	sigs := []schema.Signal{
		mkSignal("a", "Go", 800, []string{"cobra"}, "mit", time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)),
		mkSignal("b", "Go", 400, []string{"http"}, "apache-2.0", time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)),
		mkSignal("c", "Python", 100, []string{"flask"}, "mit", time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)),
		mkSignal("d", "JavaScript", 2000, []string{"react"}, "mit", time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)),
	}
	now := time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC)

	var first []byte
	for i := 0; i < 3; i++ {
		recs, err := Score(fp, sigs, WithNow(func() time.Time { return now }))
		if err != nil {
			t.Fatalf("Score[%d]: %v", i, err)
		}
		raw, err := json.Marshal(recs)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if i == 0 {
			first = raw
			continue
		}
		if !bytes.Equal(first, raw) {
			t.Fatalf("non-deterministic output on run %d:\nfirst=%s\ngot  =%s", i, first, raw)
		}
	}
}

func TestScoreSchemaVersionAndIDs(t *testing.T) {
	fp := goRepoFingerprint()
	sig := mkSignal("xyz", "Go", 800, []string{"cobra"}, "mit", time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC))

	recs, err := Score(fp, []schema.Signal{sig})
	if err != nil {
		t.Fatalf("Score: %v", err)
	}
	if recs[0].SchemaVersion != schema.RecommendationSchemaVersion {
		t.Errorf("SchemaVersion = %d, want %d", recs[0].SchemaVersion, schema.RecommendationSchemaVersion)
	}
	if recs[0].ID == "" {
		t.Errorf("rec ID empty")
	}
	if recs[0].ID == sig.ID {
		t.Errorf("rec ID must differ from signal ID")
	}
}

func TestScoreEmpty(t *testing.T) {
	fp := goRepoFingerprint()
	recs, err := Score(fp, nil)
	if err != nil {
		t.Fatalf("Score: %v", err)
	}
	if len(recs) != 0 {
		t.Errorf("expected empty recs, got %d", len(recs))
	}
}
