package ingest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestGithubTrendingName(t *testing.T) {
	src := NewGithubTrending()
	if got := src.Name(); got != "github_trending" {
		t.Errorf("Name() = %q, want %q", got, "github_trending")
	}
}

func TestGithubTrendingFetchSuccess(t *testing.T) {
	now := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	payload := map[string]any{
		"items": []map[string]any{
			{
				"full_name":        "example/repo",
				"html_url":         "https://github.com/example/repo",
				"description":      "An example",
				"language":         "Go",
				"stargazers_count": 421,
				"license":          map[string]any{"key": "mit"},
				"pushed_at":        "2026-05-19T10:00:00Z",
				"topics":           []string{"ai", "cli"},
			},
			{
				"full_name":        "another/thing",
				"html_url":         "https://github.com/another/thing",
				"description":      "",
				"language":         "JavaScript",
				"stargazers_count": 75,
				"license":           nil,
				"pushed_at":        "2026-05-18T08:30:00Z",
				"topics":           []string{},
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/search/repositories" {
			http.Error(w, "wrong path: "+r.URL.Path, http.StatusBadRequest)
			return
		}
		q := r.URL.Query().Get("q")
		if q == "" || !strings.Contains(q, "stars:") {
			http.Error(w, "bad query: "+q, http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	}))
	defer srv.Close()

	src := NewGithubTrending(
		WithGithubBaseURL(srv.URL),
		WithGithubNow(func() time.Time { return now }),
		WithGithubSleep(func(time.Duration) {}),
	)

	sigs, err := src.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(sigs) != 2 {
		t.Fatalf("len(sigs) = %d, want 2", len(sigs))
	}

	first := sigs[0]
	if first.URL != "https://github.com/example/repo" {
		t.Errorf("URL = %q", first.URL)
	}
	if first.Title != "example/repo" {
		t.Errorf("Title = %q", first.Title)
	}
	if first.Stars != 421 {
		t.Errorf("Stars = %d", first.Stars)
	}
	if first.Language != "Go" {
		t.Errorf("Language = %q", first.Language)
	}
	if first.License != "mit" {
		t.Errorf("License = %q", first.License)
	}
	if first.Source != "github_trending" {
		t.Errorf("Source = %q", first.Source)
	}
	if len(first.TopicTags) != 2 || first.TopicTags[0] != "ai" {
		t.Errorf("TopicTags = %v", first.TopicTags)
	}
	if !first.CapturedAt.Equal(now) {
		t.Errorf("CapturedAt = %v, want %v", first.CapturedAt, now)
	}
	wantLastCommit := time.Date(2026, 5, 19, 10, 0, 0, 0, time.UTC)
	if !first.LastCommit.Equal(wantLastCommit) {
		t.Errorf("LastCommit = %v, want %v", first.LastCommit, wantLastCommit)
	}

	if first.ID == "" {
		t.Errorf("ID is empty")
	}
	if signalID("github_trending", first.URL) != first.ID {
		t.Errorf("ID is not deterministic for URL")
	}

	if sigs[1].License != "" {
		t.Errorf("expected empty License when API returns null, got %q", sigs[1].License)
	}
}

func TestGithubTrendingFetchEmpty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"items":[]}`))
	}))
	defer srv.Close()

	src := NewGithubTrending(
		WithGithubBaseURL(srv.URL),
		WithGithubSleep(func(time.Duration) {}),
	)

	sigs, err := src.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(sigs) != 0 {
		t.Fatalf("expected 0 signals, got %d", len(sigs))
	}
}

func TestGithubTrendingFetchTimeoutRetries(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		time.Sleep(150 * time.Millisecond)
		_, _ = w.Write([]byte(`{"items":[]}`))
	}))
	defer srv.Close()

	src := NewGithubTrending(
		WithGithubBaseURL(srv.URL),
		WithGithubTimeout(25*time.Millisecond),
		WithGithubSleep(func(time.Duration) {}),
	)

	_, err := src.Fetch(context.Background())
	if err == nil {
		t.Fatalf("expected timeout error, got nil")
	}
	if got := atomic.LoadInt32(&hits); got != githubMaxRetries {
		t.Errorf("retried %d times, want %d", got, githubMaxRetries)
	}
}

func TestGithubTrendingCacheRoundtrip(t *testing.T) {
	now := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	var requestCount int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requestCount, 1)
		_, _ = w.Write([]byte(`{"items":[{"full_name":"a/b","html_url":"https://github.com/a/b","language":"Go","stargazers_count":100,"pushed_at":"2026-05-19T00:00:00Z"}]}`))
	}))
	defer srv.Close()

	cacheDir := t.TempDir()
	src := NewGithubTrending(
		WithGithubBaseURL(srv.URL),
		WithGithubCacheDir(cacheDir),
		WithGithubNow(func() time.Time { return now }),
		WithGithubSleep(func(time.Duration) {}),
	)

	first, err := src.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch 1: %v", err)
	}
	if atomic.LoadInt32(&requestCount) != 1 {
		t.Fatalf("expected 1 request, got %d", requestCount)
	}

	cachePath := filepath.Join(cacheDir, "github_trending-2026-05-20.json")
	if _, err := os.Stat(cachePath); err != nil {
		t.Fatalf("cache file missing: %v", err)
	}

	second, err := src.Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch 2: %v", err)
	}
	if got := atomic.LoadInt32(&requestCount); got != 1 {
		t.Errorf("expected 1 request after cache hit, got %d", got)
	}
	if len(first) != len(second) || first[0].URL != second[0].URL {
		t.Errorf("cache mismatch: %+v vs %+v", first, second)
	}
}

func TestGithubTrendingTokenHeader(t *testing.T) {
	var seenAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenAuth = r.Header.Get("Authorization")
		_, _ = w.Write([]byte(`{"items":[]}`))
	}))
	defer srv.Close()

	t.Setenv("GITHUB_TOKEN", "tok-xyz")
	src := NewGithubTrending(
		WithGithubBaseURL(srv.URL),
		WithGithubSleep(func(time.Duration) {}),
	)
	if _, err := src.Fetch(context.Background()); err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if seenAuth != "Bearer tok-xyz" {
		t.Errorf("Authorization = %q, want Bearer tok-xyz", seenAuth)
	}
}

func TestGithubTrendingServerErrorRetries(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()

	src := NewGithubTrending(
		WithGithubBaseURL(srv.URL),
		WithGithubSleep(func(time.Duration) {}),
	)
	if _, err := src.Fetch(context.Background()); err == nil {
		t.Fatal("expected error on 500 response")
	}
	if got := atomic.LoadInt32(&hits); got != githubMaxRetries {
		t.Errorf("retried %d times, want %d", got, githubMaxRetries)
	}
}
