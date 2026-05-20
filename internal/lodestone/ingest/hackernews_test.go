package ingest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

type hnFakeItem struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Score int    `json:"score"`
	By    string `json:"by"`
	Time  int64  `json:"time"`
}

func newHackerNewsTestServer(top []int, items map[int]hnFakeItem) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/v0/topstories.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(top)
	})
	mux.HandleFunc("/v0/item/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		base := strings.TrimPrefix(r.URL.Path, "/v0/item/")
		base = strings.TrimSuffix(base, ".json")
		var id int
		if _, err := fmt.Sscanf(base, "%d", &id); err != nil {
			http.Error(w, "bad id", http.StatusBadRequest)
			return
		}
		item, ok := items[id]
		if !ok {
			_, _ = w.Write([]byte("null"))
			return
		}
		_ = json.NewEncoder(w).Encode(item)
	})
	return httptest.NewServer(mux)
}

func TestHackerNewsName(t *testing.T) {
	if got := NewHackerNews().Name(); got != "hackernews" {
		t.Errorf("Name() = %q, want %q", got, "hackernews")
	}
}

func TestHackerNewsFetchFiltersAndMaps(t *testing.T) {
	now := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	top := []int{1, 2, 3, 4, 5}
	items := map[int]hnFakeItem{
		1: {ID: 1, Type: "story", Title: "Claude 4.7 lands with improved reasoning", URL: "https://example.test/claude47", Score: 412, By: "alice", Time: 1_716_000_000},
		2: {ID: 2, Type: "job", Title: "Hiring senior LLM engineer", URL: "https://example.test/job", Score: 5},
		3: {ID: 3, Type: "story", Title: "New CSS layout primitives", URL: "https://example.test/css", Score: 200},
		4: {ID: 4, Type: "story", Title: "MCP servers for terminal agents", URL: "", Score: 88, By: "bob"},
		5: {ID: 5, Type: "comment", Title: "Re: previous post", URL: "https://example.test/c", Score: 0},
	}

	srv := newHackerNewsTestServer(top, items)
	defer srv.Close()

	src := NewHackerNews(
		WithHackerNewsBaseURL(srv.URL),
		WithHackerNewsNow(func() time.Time { return now }),
		WithHackerNewsSleep(func(time.Duration) {}),
	)

	sigs, err := src.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(sigs) != 2 {
		t.Fatalf("len(sigs) = %d, want 2 (matching stories); got: %+v", len(sigs), sigs)
	}

	byTitle := map[string]bool{}
	for _, s := range sigs {
		byTitle[s.Title] = true
		if s.Source != "hackernews" {
			t.Errorf("Source = %q", s.Source)
		}
		if !s.CapturedAt.Equal(now) {
			t.Errorf("CapturedAt mismatch: %v vs %v", s.CapturedAt, now)
		}
	}
	if !byTitle["Claude 4.7 lands with improved reasoning"] || !byTitle["MCP servers for terminal agents"] {
		t.Errorf("missing expected stories: got titles %+v", byTitle)
	}

	for _, s := range sigs {
		if strings.Contains(s.Title, "MCP servers") {
			wantURL := fmt.Sprintf("https://news.ycombinator.com/item?id=%d", 4)
			if s.URL != wantURL {
				t.Errorf("expected fallback HN URL %q, got %q", wantURL, s.URL)
			}
			gotTags := map[string]bool{}
			for _, tag := range s.TopicTags {
				gotTags[tag] = true
			}
			if !gotTags["mcp"] {
				t.Errorf("expected mcp tag, got %v", s.TopicTags)
			}
		}
		if strings.Contains(s.Title, "Claude") {
			if s.Stars != 412 {
				t.Errorf("score → Stars mapping wrong: got %d", s.Stars)
			}
			gotTags := map[string]bool{}
			for _, tag := range s.TopicTags {
				gotTags[tag] = true
			}
			if !gotTags["claude"] {
				t.Errorf("expected claude tag, got %v", s.TopicTags)
			}
		}
	}
}

func TestHackerNewsFetchHonorsFinalLimit(t *testing.T) {
	now := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)

	top := make([]int, 10)
	items := map[int]hnFakeItem{}
	for i := 0; i < 10; i++ {
		id := i + 1
		top[i] = id
		items[id] = hnFakeItem{
			ID:    id,
			Type:  "story",
			Title: fmt.Sprintf("agent post #%d", id),
			URL:   fmt.Sprintf("https://example.test/%d", id),
			Score: 100 - i,
		}
	}

	srv := newHackerNewsTestServer(top, items)
	defer srv.Close()

	src := NewHackerNews(
		WithHackerNewsBaseURL(srv.URL),
		WithHackerNewsNow(func() time.Time { return now }),
		WithHackerNewsFinalLimit(3),
		WithHackerNewsSleep(func(time.Duration) {}),
	)

	sigs, err := src.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(sigs) != 3 {
		t.Fatalf("len(sigs) = %d, want 3 (final limit)", len(sigs))
	}
}

func TestHackerNewsFetchEmptyTop(t *testing.T) {
	srv := newHackerNewsTestServer([]int{}, nil)
	defer srv.Close()

	src := NewHackerNews(
		WithHackerNewsBaseURL(srv.URL),
		WithHackerNewsSleep(func(time.Duration) {}),
	)
	sigs, err := src.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(sigs) != 0 {
		t.Fatalf("expected 0 signals, got %d", len(sigs))
	}
}

func TestHackerNewsCacheRoundtrip(t *testing.T) {
	now := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	var hits int32
	mux := http.NewServeMux()
	mux.HandleFunc("/v0/topstories.json", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		_, _ = w.Write([]byte(`[1]`))
	})
	mux.HandleFunc("/v0/item/", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		_, _ = w.Write([]byte(`{"id":1,"type":"story","title":"ai is fun","url":"https://example.test/x","score":42}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	cacheDir := t.TempDir()
	src := NewHackerNews(
		WithHackerNewsBaseURL(srv.URL),
		WithHackerNewsCacheDir(cacheDir),
		WithHackerNewsNow(func() time.Time { return now }),
		WithHackerNewsSleep(func(time.Duration) {}),
	)

	first, err := src.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch 1: %v", err)
	}
	if len(first) != 1 {
		t.Fatalf("expected 1 signal, got %d", len(first))
	}
	firstHits := atomic.LoadInt32(&hits)

	cachePath := filepath.Join(cacheDir, "hackernews-2026-05-20.json")
	if _, err := os.Stat(cachePath); err != nil {
		t.Fatalf("cache file missing: %v", err)
	}

	second, err := src.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch 2: %v", err)
	}
	if atomic.LoadInt32(&hits) != firstHits {
		t.Errorf("expected no new hits on cache read, got %d -> %d", firstHits, atomic.LoadInt32(&hits))
	}
	if len(second) != 1 || second[0].URL != first[0].URL {
		t.Errorf("cache roundtrip mismatch: %+v vs %+v", first, second)
	}
}

func TestHackerNewsCustomKeywords(t *testing.T) {
	now := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	top := []int{1, 2}
	items := map[int]hnFakeItem{
		1: {ID: 1, Type: "story", Title: "Rust async runtime improvements", URL: "https://example.test/rust", Score: 88},
		2: {ID: 2, Type: "story", Title: "AI agents in production", URL: "https://example.test/ai", Score: 200},
	}
	srv := newHackerNewsTestServer(top, items)
	defer srv.Close()

	src := NewHackerNews(
		WithHackerNewsBaseURL(srv.URL),
		WithHackerNewsKeywords([]string{"rust"}),
		WithHackerNewsNow(func() time.Time { return now }),
		WithHackerNewsSleep(func(time.Duration) {}),
	)
	sigs, err := src.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(sigs) != 1 || !strings.Contains(sigs[0].Title, "Rust") {
		t.Errorf("expected only Rust story, got %+v", sigs)
	}
}
