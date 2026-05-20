package ingest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const changelogHTMLFixture = `<!DOCTYPE html>
<html><body>
  <main>
    <h2 id="2026-05-19-opus-4-7">2026-05-19 — Opus 4.7 released</h2>
    <p>Improved reasoning and tool use.</p>
    <h3 id="2026-05-12-mcp">2026-05-12 — MCP server registry expanded</h3>
    <p>New community servers added.</p>
    <h2 id="2026-04-30-sonnet">2026-04-30: Sonnet 4.6 update</h2>
  </main>
</body></html>`

func TestAnthropicChangelogName(t *testing.T) {
	if got := NewAnthropicChangelog().Name(); got != "anthropic_changelog" {
		t.Errorf("Name() = %q", got)
	}
}

func TestOpenAIChangelogName(t *testing.T) {
	if got := NewOpenAIChangelog().Name(); got != "openai_changelog" {
		t.Errorf("Name() = %q", got)
	}
}

func TestChangelogScrapeFixture(t *testing.T) {
	now := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(changelogHTMLFixture))
	}))
	defer srv.Close()

	src := NewAnthropicChangelog(
		WithChangelogURL(srv.URL),
		WithChangelogNow(func() time.Time { return now }),
		WithChangelogSleep(func(time.Duration) {}),
	)
	sigs, err := src.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(sigs) != 3 {
		t.Fatalf("got %d signals, want 3", len(sigs))
	}

	first := sigs[0]
	if !strings.Contains(first.Title, "Opus 4.7") {
		t.Errorf("Title = %q", first.Title)
	}
	if !strings.Contains(first.URL, "#2026-05-19-opus-4-7") {
		t.Errorf("URL should embed slug, got %q", first.URL)
	}
	if first.Source != "anthropic_changelog" {
		t.Errorf("Source = %q", first.Source)
	}
	wantDate := time.Date(2026, 5, 19, 0, 0, 0, 0, time.UTC)
	if !first.LastCommit.Equal(wantDate) {
		t.Errorf("LastCommit = %v, want %v", first.LastCommit, wantDate)
	}
}

func TestChangelogMaxEntries(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(changelogHTMLFixture))
	}))
	defer srv.Close()
	src := NewOpenAIChangelog(
		WithChangelogURL(srv.URL),
		WithChangelogMax(1),
		WithChangelogSleep(func(time.Duration) {}),
	)
	sigs, err := src.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(sigs) != 1 {
		t.Fatalf("max limit ignored, got %d", len(sigs))
	}
}

func TestChangelogEmpty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html><body>no headings here</body></html>`))
	}))
	defer srv.Close()
	src := NewAnthropicChangelog(
		WithChangelogURL(srv.URL),
		WithChangelogSleep(func(time.Duration) {}),
	)
	sigs, err := src.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(sigs) != 0 {
		t.Errorf("got %d signals, want 0", len(sigs))
	}
}

func TestParseChangelogHTMLExtractsTitleAndDate(t *testing.T) {
	entries := parseChangelogHTML(changelogHTMLFixture)
	if len(entries) != 3 {
		t.Fatalf("got %d entries, want 3", len(entries))
	}
	if entries[0].Title != "Opus 4.7 released" {
		t.Errorf("first.Title = %q", entries[0].Title)
	}
	if entries[1].Slug != "2026-05-12-mcp" {
		t.Errorf("second.Slug = %q", entries[1].Slug)
	}
	if entries[2].Date.Format("2006-01-02") != "2026-04-30" {
		t.Errorf("third.Date = %v", entries[2].Date)
	}
}
