package ingest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

const arxivAtomFixture = `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>ArXiv Query</title>
  <entry>
    <id>http://arxiv.org/abs/2606.01234v1</id>
    <updated>2026-05-19T12:00:00Z</updated>
    <published>2026-05-18T12:00:00Z</published>
    <title>Multi-Agent Reasoning with LLMs</title>
    <summary>We present a framework for multi-agent reasoning…</summary>
    <author><name>Alice Researcher</name></author>
    <category term="cs.AI"/>
    <category term="cs.CL"/>
  </entry>
  <entry>
    <id>http://arxiv.org/abs/2606.05678v1</id>
    <updated>2026-05-15T08:00:00Z</updated>
    <published>2026-05-15T08:00:00Z</published>
    <title>Efficient MCP Servers</title>
    <summary>An optimized runtime for the Model Context Protocol…</summary>
    <author><name>Bob Engineer</name></author>
    <category term="cs.DC"/>
  </entry>
</feed>`

func TestArXivName(t *testing.T) {
	if got := NewArXiv().Name(); got != "arxiv" {
		t.Errorf("Name() = %q", got)
	}
}

func TestArXivFetchParsesAtom(t *testing.T) {
	now := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/query" {
			http.Error(w, "wrong path "+r.URL.Path, http.StatusBadRequest)
			return
		}
		q := r.URL.Query()
		if q.Get("search_query") == "" || q.Get("sortBy") != "submittedDate" {
			http.Error(w, "bad query", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/atom+xml")
		_, _ = w.Write([]byte(arxivAtomFixture))
	}))
	defer srv.Close()

	src := NewArXiv(
		WithArXivBaseURL(srv.URL),
		WithArXivNow(func() time.Time { return now }),
		WithArXivSleep(func(time.Duration) {}),
	)
	sigs, err := src.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(sigs) != 2 {
		t.Fatalf("got %d signals, want 2", len(sigs))
	}

	first := sigs[0]
	if !strings.Contains(first.Title, "Multi-Agent") {
		t.Errorf("Title = %q", first.Title)
	}
	if first.URL != "http://arxiv.org/abs/2606.01234v1" {
		t.Errorf("URL = %q", first.URL)
	}
	if first.Source != "arxiv" {
		t.Errorf("Source = %q", first.Source)
	}
	if !first.CapturedAt.Equal(now) {
		t.Errorf("CapturedAt = %v", first.CapturedAt)
	}
	if len(first.TopicTags) != 2 || first.TopicTags[0] != "cs.AI" {
		t.Errorf("TopicTags = %v", first.TopicTags)
	}
	wantPublished := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	if !first.LastCommit.Equal(wantPublished) {
		t.Errorf("LastCommit = %v, want %v", first.LastCommit, wantPublished)
	}
	if first.ID == "" {
		t.Errorf("ID empty")
	}
}

func TestArXivFetchEmpty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<?xml version="1.0"?><feed xmlns="http://www.w3.org/2005/Atom"></feed>`))
	}))
	defer srv.Close()
	src := NewArXiv(WithArXivBaseURL(srv.URL), WithArXivSleep(func(time.Duration) {}))
	sigs, err := src.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(sigs) != 0 {
		t.Fatalf("expected 0, got %d", len(sigs))
	}
}

func TestArXivFetchRetriesOn5xx(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()
	src := NewArXiv(WithArXivBaseURL(srv.URL), WithArXivSleep(func(time.Duration) {}))
	if _, err := src.Fetch(context.Background()); err == nil {
		t.Fatal("expected error")
	}
	if got := atomic.LoadInt32(&hits); got != defaultMaxRetries {
		t.Errorf("retries = %d, want %d", got, defaultMaxRetries)
	}
}

func TestArXivCacheRoundtrip(t *testing.T) {
	now := time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC)
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		_, _ = w.Write([]byte(arxivAtomFixture))
	}))
	defer srv.Close()
	dir := t.TempDir()
	src := NewArXiv(
		WithArXivBaseURL(srv.URL),
		WithArXivCacheDir(dir),
		WithArXivNow(func() time.Time { return now }),
		WithArXivSleep(func(time.Duration) {}),
	)
	if _, err := src.Fetch(context.Background()); err != nil {
		t.Fatalf("Fetch 1: %v", err)
	}
	if atomic.LoadInt32(&hits) != 1 {
		t.Fatalf("hits = %d, want 1", hits)
	}
	if _, err := os.Stat(filepath.Join(dir, "arxiv-2026-05-20.json")); err != nil {
		t.Fatalf("cache missing: %v", err)
	}
	if _, err := src.Fetch(context.Background()); err != nil {
		t.Fatalf("Fetch 2: %v", err)
	}
	if got := atomic.LoadInt32(&hits); got != 1 {
		t.Errorf("expected cache hit (1 server call), got %d", got)
	}
}
