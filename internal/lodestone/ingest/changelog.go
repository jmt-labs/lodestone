package ingest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/jmt-labs/lodestone/internal/lodestone/schema"
)

const (
	anthropicChangelogName     = "anthropic_changelog"
	anthropicChangelogDefault  = "https://docs.anthropic.com/en/release-notes/api"
	openaiChangelogName        = "openai_changelog"
	openaiChangelogDefault     = "https://platform.openai.com/docs/changelog"
	changelogDefaultTimeout    = 15 * time.Second
	changelogDefaultMaxEntries = 30
)

type ChangelogScraper struct {
	name       string
	pageURL    string
	httpClient *http.Client
	cacheDir   string
	maxEntries int
	timeout    time.Duration
	now        func() time.Time
	sleep      func(time.Duration)
}

type ChangelogOption func(*ChangelogScraper)

func WithChangelogURL(u string) ChangelogOption { return func(c *ChangelogScraper) { c.pageURL = u } }
func WithChangelogCacheDir(d string) ChangelogOption {
	return func(c *ChangelogScraper) { c.cacheDir = d }
}
func WithChangelogMax(n int) ChangelogOption {
	return func(c *ChangelogScraper) { c.maxEntries = n }
}
func WithChangelogTimeout(d time.Duration) ChangelogOption {
	return func(c *ChangelogScraper) { c.timeout = d }
}
func WithChangelogNow(fn func() time.Time) ChangelogOption {
	return func(c *ChangelogScraper) { c.now = fn }
}
func WithChangelogSleep(fn func(time.Duration)) ChangelogOption {
	return func(c *ChangelogScraper) { c.sleep = fn }
}

func newChangelogScraper(name, defaultURL string, opts ...ChangelogOption) *ChangelogScraper {
	c := &ChangelogScraper{
		name:       name,
		pageURL:    defaultURL,
		maxEntries: changelogDefaultMaxEntries,
		timeout:    changelogDefaultTimeout,
		now:        func() time.Time { return time.Now().UTC() },
		sleep:      time.Sleep,
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.httpClient == nil {
		c.httpClient = &http.Client{Timeout: c.timeout}
	}
	return c
}

func NewAnthropicChangelog(opts ...ChangelogOption) *ChangelogScraper {
	return newChangelogScraper(anthropicChangelogName, anthropicChangelogDefault, opts...)
}

func NewOpenAIChangelog(opts ...ChangelogOption) *ChangelogScraper {
	return newChangelogScraper(openaiChangelogName, openaiChangelogDefault, opts...)
}

func (c *ChangelogScraper) Name() string { return c.name }

func (c *ChangelogScraper) Fetch(ctx context.Context) ([]schema.Signal, error) {
	now := c.now()
	cp := cachePath(c.cacheDir, c.name, now)
	if cached, ok, err := loadCache(cp); err != nil {
		return nil, fmt.Errorf("read cache: %w", err)
	} else if ok {
		return cached, nil
	}

	cfg := defaultRetryConfig(c.sleep)
	sigs, err := retryFetch(ctx, cfg, c.name, func() ([]schema.Signal, error) {
		return c.fetchOnce(ctx, now)
	})
	if err != nil {
		return nil, err
	}

	if err := saveCache(cp, sigs); err != nil {
		return nil, fmt.Errorf("write cache: %w", err)
	}
	return sigs, nil
}

func (c *ChangelogScraper) fetchOnce(ctx context.Context, now time.Time) ([]schema.Signal, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.pageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Accept", "text/html")
	req.Header.Set("User-Agent", "lodestone/0.1 (+https://github.com/jmt-labs/lodestone)")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &httpStatusError{Source: c.name, Status: resp.StatusCode, Body: string(body)}
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	entries := parseChangelogHTML(string(raw))
	if len(entries) > c.maxEntries {
		entries = entries[:c.maxEntries]
	}

	sigs := make([]schema.Signal, 0, len(entries))
	for _, e := range entries {
		sigs = append(sigs, e.toSignal(c.name, c.pageURL, now))
	}
	return sigs, nil
}

type changelogEntry struct {
	Title string
	Date  time.Time
	Slug  string
}

func (e changelogEntry) toSignal(source, basePage string, now time.Time) schema.Signal {
	u := basePage
	if e.Slug != "" {
		u = basePage + "#" + e.Slug
	}
	return schema.Signal{
		SchemaVersion: schema.SignalSchemaVersion,
		ID:            signalID(source, u),
		Source:        source,
		URL:           u,
		Title:         e.Title,
		CapturedAt:    now,
		LastCommit:    e.Date,
	}
}

var (
	// <h2>…</h2> oder <h3>…</h3> mit Datum am Anfang oder im Datums-Span davor.
	changelogHeadingRE = regexp.MustCompile(`(?is)<h([23])(?:[^>]*?\s+id=["']([^"']+)["'])?[^>]*>(.*?)</h[23]>`)
	changelogDateRE    = regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`)
	tagStripRE         = regexp.MustCompile(`<[^>]+>`)
	whitespaceRE       = regexp.MustCompile(`\s+`)
)

func parseChangelogHTML(html string) []changelogEntry {
	matches := changelogHeadingRE.FindAllStringSubmatch(html, -1)
	if len(matches) == 0 {
		return nil
	}
	var out []changelogEntry
	for _, m := range matches {
		slug := m[2]
		inner := stripTags(m[3])
		if inner == "" {
			continue
		}
		date := extractDate(inner)
		title := strings.TrimSpace(strings.TrimPrefix(inner, date.Format("2006-01-02")))
		title = strings.TrimSpace(strings.Trim(title, "·-—:"))
		if title == "" {
			title = inner
		}
		out = append(out, changelogEntry{
			Title: title,
			Date:  date,
			Slug:  slug,
		})
	}
	return out
}

func stripTags(s string) string {
	s = tagStripRE.ReplaceAllString(s, " ")
	s = whitespaceRE.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

func extractDate(s string) time.Time {
	if m := changelogDateRE.FindString(s); m != "" {
		if t, err := time.Parse("2006-01-02", m); err == nil {
			return t.UTC()
		}
	}
	return time.Time{}
}

var (
	_ Source = (*ChangelogScraper)(nil)
)
