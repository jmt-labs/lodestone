package ingest

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jmt-labs/lodestone/internal/lodestone/schema"
)

const (
	arxivName            = "arxiv"
	arxivDefaultBaseURL  = "https://export.arxiv.org"
	arxivDefaultCategory = "cat:cs.AI"
	arxivDefaultMax      = 30
	arxivDefaultTimeout  = 15 * time.Second
)

type ArXiv struct {
	baseURL    string
	httpClient *http.Client
	cacheDir   string
	query      string
	maxResults int
	timeout    time.Duration
	now        func() time.Time
	sleep      func(time.Duration)
}

type ArXivOption func(*ArXiv)

func WithArXivBaseURL(u string) ArXivOption  { return func(a *ArXiv) { a.baseURL = u } }
func WithArXivCacheDir(d string) ArXivOption { return func(a *ArXiv) { a.cacheDir = d } }
func WithArXivQuery(q string) ArXivOption    { return func(a *ArXiv) { a.query = q } }
func WithArXivMax(n int) ArXivOption         { return func(a *ArXiv) { a.maxResults = n } }
func WithArXivTimeout(d time.Duration) ArXivOption {
	return func(a *ArXiv) { a.timeout = d }
}
func WithArXivNow(fn func() time.Time) ArXivOption {
	return func(a *ArXiv) { a.now = fn }
}
func WithArXivSleep(fn func(time.Duration)) ArXivOption {
	return func(a *ArXiv) { a.sleep = fn }
}

func NewArXiv(opts ...ArXivOption) *ArXiv {
	a := &ArXiv{
		baseURL:    arxivDefaultBaseURL,
		query:      arxivDefaultCategory,
		maxResults: arxivDefaultMax,
		timeout:    arxivDefaultTimeout,
		now:        func() time.Time { return time.Now().UTC() },
		sleep:      time.Sleep,
	}
	for _, opt := range opts {
		opt(a)
	}
	if a.httpClient == nil {
		a.httpClient = &http.Client{Timeout: a.timeout}
	}
	return a
}

func (a *ArXiv) Name() string { return arxivName }

func (a *ArXiv) Fetch(ctx context.Context) ([]schema.Signal, error) {
	now := a.now()
	cp := cachePath(a.cacheDir, arxivName, now)
	if cached, ok, err := loadCache(cp); err != nil {
		return nil, fmt.Errorf("read cache: %w", err)
	} else if ok {
		return cached, nil
	}

	cfg := defaultRetryConfig(a.sleep)
	sigs, err := retryFetch(ctx, cfg, arxivName, func() ([]schema.Signal, error) {
		return a.fetchOnce(ctx, now)
	})
	if err != nil {
		return nil, err
	}

	if err := saveCache(cp, sigs); err != nil {
		return nil, fmt.Errorf("write cache: %w", err)
	}
	return sigs, nil
}

func (a *ArXiv) fetchOnce(ctx context.Context, now time.Time) ([]schema.Signal, error) {
	u, err := url.Parse(a.baseURL + "/api/query")
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}
	qs := u.Query()
	qs.Set("search_query", a.query)
	qs.Set("sortBy", "submittedDate")
	qs.Set("sortOrder", "descending")
	qs.Set("max_results", fmt.Sprintf("%d", a.maxResults))
	u.RawQuery = qs.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Accept", "application/atom+xml")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &httpStatusError{Source: arxivName, Status: resp.StatusCode, Body: string(body)}
	}

	var feed arxivFeed
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("decode atom: %w", err)
	}

	sigs := make([]schema.Signal, 0, len(feed.Entries))
	for _, e := range feed.Entries {
		sigs = append(sigs, e.toSignal(now))
	}
	return sigs, nil
}

type arxivFeed struct {
	XMLName xml.Name     `xml:"http://www.w3.org/2005/Atom feed"`
	Entries []arxivEntry `xml:"entry"`
}

type arxivEntry struct {
	ID         string          `xml:"id"`
	Title      string          `xml:"title"`
	Summary    string          `xml:"summary"`
	Published  time.Time       `xml:"published"`
	Updated    time.Time       `xml:"updated"`
	Categories []arxivCategory `xml:"category"`
	Authors    []arxivAuthor   `xml:"author"`
}

type arxivCategory struct {
	Term string `xml:"term,attr"`
}

type arxivAuthor struct {
	Name string `xml:"name"`
}

func (e arxivEntry) toSignal(now time.Time) schema.Signal {
	tags := make([]string, 0, len(e.Categories))
	for _, c := range e.Categories {
		if c.Term != "" {
			tags = append(tags, c.Term)
		}
	}
	idURL := strings.TrimSpace(e.ID)
	return schema.Signal{
		SchemaVersion: schema.SignalSchemaVersion,
		ID:            signalID(arxivName, idURL),
		Source:        arxivName,
		URL:           idURL,
		Title:         strings.TrimSpace(e.Title),
		Summary:       strings.TrimSpace(e.Summary),
		CapturedAt:    now,
		TopicTags:     tags,
		LastCommit:    e.Published.UTC(),
	}
}

var _ Source = (*ArXiv)(nil)
