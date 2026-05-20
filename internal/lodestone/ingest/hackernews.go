package ingest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jmt-labs/lodestone/internal/lodestone/schema"
)

const (
	hackerNewsName       = "hackernews"
	hnDefaultBaseURL     = "https://hacker-news.firebaseio.com"
	hnDefaultScanLimit   = 100
	hnDefaultFinalLimit  = 50
	hnDefaultTimeout     = 15 * time.Second
	hnInitialBackoff     = 200 * time.Millisecond
	hnMaxBackoff         = 5 * time.Second
	hnMaxRetries         = 3
	hnItemURLFormat      = "https://news.ycombinator.com/item?id=%d"
)

var hnDefaultKeywords = []string{"ai", "llm", "mcp", "claude", "agent"}

type HackerNews struct {
	baseURL    string
	httpClient *http.Client
	cacheDir   string
	keywords   []string
	scanLimit  int
	finalLimit int
	timeout    time.Duration
	now        func() time.Time
	sleep      func(time.Duration)
}

type HackerNewsOption func(*HackerNews)

func WithHackerNewsBaseURL(u string) HackerNewsOption {
	return func(h *HackerNews) { h.baseURL = u }
}

func WithHackerNewsHTTPClient(c *http.Client) HackerNewsOption {
	return func(h *HackerNews) { h.httpClient = c }
}

func WithHackerNewsCacheDir(dir string) HackerNewsOption {
	return func(h *HackerNews) { h.cacheDir = dir }
}

func WithHackerNewsKeywords(kw []string) HackerNewsOption {
	return func(h *HackerNews) { h.keywords = kw }
}

func WithHackerNewsScanLimit(n int) HackerNewsOption {
	return func(h *HackerNews) { h.scanLimit = n }
}

func WithHackerNewsFinalLimit(n int) HackerNewsOption {
	return func(h *HackerNews) { h.finalLimit = n }
}

func WithHackerNewsTimeout(d time.Duration) HackerNewsOption {
	return func(h *HackerNews) { h.timeout = d }
}

func WithHackerNewsNow(fn func() time.Time) HackerNewsOption {
	return func(h *HackerNews) { h.now = fn }
}

func WithHackerNewsSleep(fn func(time.Duration)) HackerNewsOption {
	return func(h *HackerNews) { h.sleep = fn }
}

func NewHackerNews(opts ...HackerNewsOption) *HackerNews {
	h := &HackerNews{
		baseURL:    hnDefaultBaseURL,
		keywords:   append([]string{}, hnDefaultKeywords...),
		scanLimit:  hnDefaultScanLimit,
		finalLimit: hnDefaultFinalLimit,
		timeout:    hnDefaultTimeout,
		now:        func() time.Time { return time.Now().UTC() },
		sleep:      time.Sleep,
	}
	for _, opt := range opts {
		opt(h)
	}
	if h.httpClient == nil {
		h.httpClient = &http.Client{Timeout: h.timeout}
	}
	return h
}

func (h *HackerNews) Name() string { return hackerNewsName }

func (h *HackerNews) Fetch(ctx context.Context) ([]schema.Signal, error) {
	now := h.now()
	cachePath := h.cachePath(now)

	if cachePath != "" {
		cached, ok, err := loadCache(cachePath)
		if err != nil {
			return nil, fmt.Errorf("read cache: %w", err)
		}
		if ok {
			return cached, nil
		}
	}

	topIDs, err := h.fetchTopStories(ctx)
	if err != nil {
		return nil, err
	}
	if len(topIDs) > h.scanLimit {
		topIDs = topIDs[:h.scanLimit]
	}

	keywordsLower := make([]string, len(h.keywords))
	for i, k := range h.keywords {
		keywordsLower[i] = strings.ToLower(k)
	}

	var sigs []schema.Signal
	for _, id := range topIDs {
		if len(sigs) >= h.finalLimit {
			break
		}
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		item, err := h.fetchItem(ctx, id)
		if err != nil {
			return nil, err
		}
		if item == nil || item.Type != "story" || item.Title == "" {
			continue
		}
		matched := matchKeywords(item.Title, keywordsLower)
		if len(matched) == 0 {
			continue
		}
		sigs = append(sigs, hnToSignal(*item, matched, now))
	}

	if cachePath != "" {
		if err := saveCache(cachePath, sigs); err != nil {
			return nil, fmt.Errorf("write cache: %w", err)
		}
	}
	return sigs, nil
}

func (h *HackerNews) cachePath(now time.Time) string {
	if h.cacheDir == "" {
		return ""
	}
	return filepath.Join(h.cacheDir, fmt.Sprintf("%s-%s.json", hackerNewsName, now.Format("2006-01-02")))
}

func (h *HackerNews) fetchTopStories(ctx context.Context) ([]int, error) {
	var ids []int
	err := h.requestWithRetry(ctx, h.baseURL+"/v0/topstories.json", &ids)
	return ids, err
}

func (h *HackerNews) fetchItem(ctx context.Context, id int) (*hnItem, error) {
	endpoint := fmt.Sprintf("%s/v0/item/%d.json", h.baseURL, id)
	var raw json.RawMessage
	if err := h.requestWithRetry(ctx, endpoint, &raw); err != nil {
		return nil, err
	}
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}
	var item hnItem
	if err := json.Unmarshal(raw, &item); err != nil {
		return nil, fmt.Errorf("decode item %d: %w", id, err)
	}
	return &item, nil
}

func (h *HackerNews) requestWithRetry(ctx context.Context, endpoint string, out any) error {
	var lastErr error
	backoff := hnInitialBackoff
	for attempt := 0; attempt < hnMaxRetries; attempt++ {
		if attempt > 0 {
			if err := ctx.Err(); err != nil {
				return err
			}
			h.sleep(backoff)
			backoff *= 2
			if backoff > hnMaxBackoff {
				backoff = hnMaxBackoff
			}
		}
		err := h.doRequest(ctx, endpoint, out)
		if err == nil {
			return nil
		}
		lastErr = err
		if !isRetryable(err) {
			return err
		}
	}
	return fmt.Errorf("hackernews: max retries exceeded: %w", lastErr)
}

func (h *HackerNews) doRequest(ctx context.Context, endpoint string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return &httpStatusError{Status: resp.StatusCode, Body: string(body)}
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode body: %w", err)
	}
	return nil
}

type hnItem struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Title string `json:"title"`
	URL   string `json:"url"`
	By    string `json:"by"`
	Time  int64  `json:"time"`
	Score int    `json:"score"`
}

func matchKeywords(title string, keywordsLower []string) []string {
	lower := strings.ToLower(title)
	seen := map[string]bool{}
	var matched []string
	for _, k := range keywordsLower {
		if k == "" || seen[k] {
			continue
		}
		if strings.Contains(lower, k) {
			seen[k] = true
			matched = append(matched, k)
		}
	}
	sort.Strings(matched)
	return matched
}

func hnToSignal(item hnItem, matched []string, now time.Time) schema.Signal {
	u := item.URL
	if u == "" {
		u = fmt.Sprintf(hnItemURLFormat, item.ID)
	}
	return schema.Signal{
		SchemaVersion: schema.SignalSchemaVersion,
		ID:            signalID(hackerNewsName, u),
		Source:        hackerNewsName,
		URL:           u,
		Title:         item.Title,
		CapturedAt:    now,
		Stars:         item.Score,
		TopicTags:     matched,
	}
}

var _ Source = (*HackerNews)(nil)
