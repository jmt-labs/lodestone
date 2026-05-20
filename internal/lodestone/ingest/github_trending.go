package ingest

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/jmt-labs/lodestone/internal/lodestone/schema"
)

const (
	githubTrendingName      = "github_trending"
	githubDefaultBaseURL    = "https://api.github.com"
	githubDefaultPerPage    = 30
	githubDefaultRecentDays = 30
	githubDefaultMinStars   = 50
	githubDefaultTimeout    = 15 * time.Second
	githubMaxRetries        = 3
	githubInitialBackoff    = 200 * time.Millisecond
	githubMaxBackoff        = 5 * time.Second
)

type GithubTrending struct {
	baseURL    string
	httpClient *http.Client
	token      string
	cacheDir   string
	perPage    int
	minStars   int
	recentDays int
	timeout    time.Duration
	now        func() time.Time
	sleep      func(time.Duration)
}

type GithubTrendingOption func(*GithubTrending)

func WithGithubBaseURL(u string) GithubTrendingOption {
	return func(g *GithubTrending) { g.baseURL = u }
}

func WithGithubCacheDir(dir string) GithubTrendingOption {
	return func(g *GithubTrending) { g.cacheDir = dir }
}

func WithGithubHTTPClient(c *http.Client) GithubTrendingOption {
	return func(g *GithubTrending) { g.httpClient = c }
}

func WithGithubTimeout(d time.Duration) GithubTrendingOption {
	return func(g *GithubTrending) { g.timeout = d }
}

func WithGithubMinStars(n int) GithubTrendingOption {
	return func(g *GithubTrending) { g.minStars = n }
}

func WithGithubRecentDays(n int) GithubTrendingOption {
	return func(g *GithubTrending) { g.recentDays = n }
}

func WithGithubNow(fn func() time.Time) GithubTrendingOption {
	return func(g *GithubTrending) { g.now = fn }
}

func WithGithubSleep(fn func(time.Duration)) GithubTrendingOption {
	return func(g *GithubTrending) { g.sleep = fn }
}

func NewGithubTrending(opts ...GithubTrendingOption) *GithubTrending {
	g := &GithubTrending{
		baseURL:    githubDefaultBaseURL,
		token:      os.Getenv("GITHUB_TOKEN"),
		perPage:    githubDefaultPerPage,
		minStars:   githubDefaultMinStars,
		recentDays: githubDefaultRecentDays,
		timeout:    githubDefaultTimeout,
		now:        func() time.Time { return time.Now().UTC() },
		sleep:      time.Sleep,
	}
	for _, opt := range opts {
		opt(g)
	}
	if g.httpClient == nil {
		g.httpClient = &http.Client{Timeout: g.timeout}
	}
	return g
}

func (g *GithubTrending) Name() string { return githubTrendingName }

func (g *GithubTrending) Fetch(ctx context.Context) ([]schema.Signal, error) {
	now := g.now()
	cachePath := g.cachePath(now)

	if cachePath != "" {
		cached, ok, err := loadCache(cachePath)
		if err != nil {
			return nil, fmt.Errorf("read cache: %w", err)
		}
		if ok {
			return cached, nil
		}
	}

	sigs, err := g.fetchWithRetry(ctx, now)
	if err != nil {
		return nil, err
	}

	if cachePath != "" {
		if err := saveCache(cachePath, sigs); err != nil {
			return nil, fmt.Errorf("write cache: %w", err)
		}
	}
	return sigs, nil
}

func (g *GithubTrending) cachePath(now time.Time) string {
	if g.cacheDir == "" {
		return ""
	}
	return filepath.Join(g.cacheDir, fmt.Sprintf("%s-%s.json", githubTrendingName, now.Format("2006-01-02")))
}

func (g *GithubTrending) fetchWithRetry(ctx context.Context, now time.Time) ([]schema.Signal, error) {
	var lastErr error
	backoff := githubInitialBackoff
	for attempt := 0; attempt < githubMaxRetries; attempt++ {
		if attempt > 0 {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			g.sleep(backoff)
			backoff *= 2
			if backoff > githubMaxBackoff {
				backoff = githubMaxBackoff
			}
		}
		sigs, err := g.fetchOnce(ctx, now)
		if err == nil {
			return sigs, nil
		}
		lastErr = err
		if !isRetryable(err) {
			return nil, err
		}
	}
	return nil, fmt.Errorf("github_trending: max retries exceeded: %w", lastErr)
}

func (g *GithubTrending) fetchOnce(ctx context.Context, now time.Time) ([]schema.Signal, error) {
	cutoff := now.AddDate(0, 0, -g.recentDays).Format("2006-01-02")
	query := fmt.Sprintf("stars:>=%d pushed:>%s", g.minStars, cutoff)

	u, err := url.Parse(g.baseURL + "/search/repositories")
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}
	qs := u.Query()
	qs.Set("q", query)
	qs.Set("sort", "stars")
	qs.Set("order", "desc")
	qs.Set("per_page", fmt.Sprintf("%d", g.perPage))
	u.RawQuery = qs.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if g.token != "" {
		req.Header.Set("Authorization", "Bearer "+g.token)
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &httpStatusError{Status: resp.StatusCode, Body: string(body)}
	}

	var payload githubSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	sigs := make([]schema.Signal, 0, len(payload.Items))
	for _, item := range payload.Items {
		sigs = append(sigs, toSignal(item, now))
	}
	return sigs, nil
}

type githubSearchResponse struct {
	Items []githubRepo `json:"items"`
}

type githubRepo struct {
	FullName        string         `json:"full_name"`
	HTMLURL         string         `json:"html_url"`
	Description     string         `json:"description"`
	Language        string         `json:"language"`
	StargazersCount int            `json:"stargazers_count"`
	License         *githubLicense `json:"license"`
	PushedAt        time.Time      `json:"pushed_at"`
	Topics          []string       `json:"topics"`
}

type githubLicense struct {
	Key string `json:"key"`
}

func toSignal(r githubRepo, now time.Time) schema.Signal {
	licenseKey := ""
	if r.License != nil {
		licenseKey = r.License.Key
	}
	return schema.Signal{
		SchemaVersion: schema.SignalSchemaVersion,
		ID:            signalID(githubTrendingName, r.HTMLURL),
		Source:        githubTrendingName,
		URL:           r.HTMLURL,
		Title:         r.FullName,
		Summary:       r.Description,
		CapturedAt:    now,
		Language:      r.Language,
		Stars:         r.StargazersCount,
		TopicTags:     r.Topics,
		License:       licenseKey,
		LastCommit:    r.PushedAt.UTC(),
	}
}

func signalID(source, urlStr string) string {
	sum := sha256.Sum256([]byte(source + "|" + urlStr))
	return "sha256:" + hex.EncodeToString(sum[:])
}

type httpStatusError struct {
	Status int
	Body   string
}

func (e *httpStatusError) Error() string {
	return fmt.Sprintf("github_trending: status %d: %s", e.Status, e.Body)
}

func isRetryable(err error) bool {
	if err == nil {
		return false
	}
	var herr *httpStatusError
	if errors.As(err, &herr) {
		return herr.Status >= 500 || herr.Status == http.StatusTooManyRequests
	}
	return true
}

func loadCache(path string) ([]schema.Signal, bool, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, false, nil
		}
		return nil, false, err
	}
	var sigs []schema.Signal
	if err := json.Unmarshal(raw, &sigs); err != nil {
		return nil, false, err
	}
	return sigs, true, nil
}

func saveCache(path string, sigs []schema.Signal) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(sigs, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".tmp.*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(raw); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}
	return os.Rename(tmpName, path)
}

var _ Source = (*GithubTrending)(nil)
