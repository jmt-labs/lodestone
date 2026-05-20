package schema

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestSignalRoundtrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	in := Signal{
		SchemaVersion:    SignalSchemaVersion,
		ID:               "sha256:abc",
		Source:           "github_trending",
		URL:              "https://github.com/example/repo",
		Title:            "example/repo",
		Summary:          "A small example",
		CapturedAt:       now,
		Language:         "Go",
		Stars:            123,
		TopicTags:        []string{"ai", "cli"},
		MaintenanceScore: 0.42,
		License:          "MIT",
		LastCommit:       now.Add(-24 * time.Hour),
	}

	raw, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var out Signal
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !reflect.DeepEqual(in, out) {
		t.Fatalf("roundtrip mismatch:\nwant=%+v\ngot =%+v", in, out)
	}
}

func TestFingerprintRoundtrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	in := Fingerprint{
		SchemaVersion:  FingerprintSchemaVersion,
		GeneratedAt:    now,
		Languages:      []string{"Go", "JavaScript"},
		Frameworks:     []string{"cobra", "react"},
		Deps:           map[string]string{"github.com/spf13/cobra": "v1.10.2"},
		LOCPerLanguage: map[string]int{"Go": 1200, "JavaScript": 450},
		TestRatio:      0.31,
		HasCI:          true,
		CIProvider:     "github_actions",
		MCPServers:     []string{"filesystem"},
		Goals:          []string{"reliability", "speed"},
		TechInterests:  []string{"mcp", "llm-tools"},
	}

	raw, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var out Fingerprint
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !reflect.DeepEqual(in, out) {
		t.Fatalf("roundtrip mismatch:\nwant=%+v\ngot =%+v", in, out)
	}
}

func TestRecommendationRoundtrip(t *testing.T) {
	in := Recommendation{
		SchemaVersion:   RecommendationSchemaVersion,
		ID:              "sha256:def",
		SignalID:        "sha256:abc",
		Compatibility:   0.83,
		Effort:          EffortS,
		ROI:             ROIHigh,
		Risk:            RiskLow,
		Rationale:       "Matches Go stack and has stable license.",
		CounterEvidence: "Project is younger than 60 days.",
		SuggestedNext:   []string{"add-dependency", "write-spec"},
	}

	raw, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var out Recommendation
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !reflect.DeepEqual(in, out) {
		t.Fatalf("roundtrip mismatch:\nwant=%+v\ngot =%+v", in, out)
	}
}

func TestWorkPackageRoundtrip(t *testing.T) {
	in := WorkPackage{
		ID:                 "WP-001",
		Type:               "task",
		Title:              "Add fetch adapter for HackerNews",
		DependsOn:          []string{"WP-000"},
		FilesAffected:      []string{"internal/lodestone/ingest/hackernews.go"},
		ExpectedArtifacts:  []string{".lodestone/cache/hackernews-2026-05-20.json"},
		Executor:           "developer",
		EstimatedMinutes:   45,
		AcceptanceCriteria: []string{"Unit-Test deckt Fetch & Timeout ab"},
	}

	raw, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var out WorkPackage
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !reflect.DeepEqual(in, out) {
		t.Fatalf("roundtrip mismatch:\nwant=%+v\ngot =%+v", in, out)
	}
}
