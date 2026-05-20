package ingest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNPMTrendingName(t *testing.T) {
	if got := NewNPMTrending().Name(); got != "npm_trending" {
		t.Errorf("Name() = %q", got)
	}
}

func TestNPMTrendingFetchMapsPackages(t *testing.T) {
	now := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	payload := map[string]any{
		"objects": []map[string]any{
			{
				"package": map[string]any{
					"name":        "@anthropic-ai/sdk",
					"version":     "0.30.0",
					"description": "Anthropic Claude SDK",
					"date":        "2026-05-18T12:00:00.000Z",
					"links":       map[string]any{"npm": "https://www.npmjs.com/package/@anthropic-ai/sdk"},
					"keywords":    []string{"ai", "anthropic", "claude"},
				},
				"score": map[string]any{"final": 0.92},
			},
			{
				"package": map[string]any{
					"name":        "@modelcontextprotocol/sdk",
					"description": "MCP SDK",
					"date":        "2026-05-17T08:00:00.000Z",
					"links":       map[string]any{"npm": "https://www.npmjs.com/package/@modelcontextprotocol/sdk"},
					"keywords":    []string{"mcp", "ai"},
				},
				"score": map[string]any{"final": 0.71},
			},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/-/v1/search" {
			http.Error(w, "wrong path", http.StatusBadRequest)
			return
		}
		text := r.URL.Query().Get("text")
		if !strings.Contains(text, "keywords:") {
			http.Error(w, "bad query: "+text, http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	}))
	defer srv.Close()

	src := NewNPMTrending(
		WithNPMBaseURL(srv.URL),
		WithNPMNow(func() time.Time { return now }),
		WithNPMSleep(func(time.Duration) {}),
	)
	sigs, err := src.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(sigs) != 2 {
		t.Fatalf("got %d signals, want 2", len(sigs))
	}

	first := sigs[0]
	if first.Title != "@anthropic-ai/sdk" {
		t.Errorf("Title = %q", first.Title)
	}
	if first.Language != "JavaScript" {
		t.Errorf("Language = %q", first.Language)
	}
	if first.Stars != 920 {
		t.Errorf("Stars = %d, want 920 (score 0.92 * 1000)", first.Stars)
	}
	if first.MaintenanceScore < 0.91 || first.MaintenanceScore > 0.93 {
		t.Errorf("MaintenanceScore = %v", first.MaintenanceScore)
	}
	if len(first.TopicTags) != 3 {
		t.Errorf("TopicTags = %v", first.TopicTags)
	}
	wantDate := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	if !first.LastCommit.Equal(wantDate) {
		t.Errorf("LastCommit = %v", first.LastCommit)
	}
}

func TestNPMTrendingEmpty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"objects":[]}`))
	}))
	defer srv.Close()
	src := NewNPMTrending(WithNPMBaseURL(srv.URL), WithNPMSleep(func(time.Duration) {}))
	sigs, err := src.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(sigs) != 0 {
		t.Fatalf("got %d, want 0", len(sigs))
	}
}

func TestNPMBuildQuery(t *testing.T) {
	n := NewNPMTrending(WithNPMKeywords("ai, mcp"))
	got := n.buildQuery()
	if !strings.Contains(got, "keywords:ai") || !strings.Contains(got, "keywords:mcp") {
		t.Errorf("buildQuery() = %q", got)
	}
}
