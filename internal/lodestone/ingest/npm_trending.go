package ingest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jmt-labs/lodestone/internal/lodestone/schema"
)

const (
	npmTrendingName    = "npm_trending"
	npmDefaultBaseURL  = "https://registry.npmjs.org"
	npmDefaultKeywords = "ai,llm,mcp,agent"
	npmDefaultSize     = 20
	npmDefaultTimeout  = 15 * time.Second
)

type NPMTrending struct {
	baseURL    string
	httpClient *http.Client
	cacheDir   string
	keywords   string
	size       int
	timeout    time.Duration
	now        func() time.Time
	sleep      func(time.Duration)
}

type NPMOption func(*NPMTrending)

func WithNPMBaseURL(u string) NPMOption  { return func(n *NPMTrending) { n.baseURL = u } }
func WithNPMCacheDir(d string) NPMOption { return func(n *NPMTrending) { n.cacheDir = d } }
func WithNPMKeywords(k string) NPMOption { return func(n *NPMTrending) { n.keywords = k } }
func WithNPMSize(s int) NPMOption        { return func(n *NPMTrending) { n.size = s } }
func WithNPMTimeout(d time.Duration) NPMOption {
	return func(n *NPMTrending) { n.timeout = d }
}
func WithNPMNow(fn func() time.Time) NPMOption      { return func(n *NPMTrending) { n.now = fn } }
func WithNPMSleep(fn func(time.Duration)) NPMOption { return func(n *NPMTrending) { n.sleep = fn } }

func NewNPMTrending(opts ...NPMOption) *NPMTrending {
	n := &NPMTrending{
		baseURL:  npmDefaultBaseURL,
		keywords: npmDefaultKeywords,
		size:     npmDefaultSize,
		timeout:  npmDefaultTimeout,
		now:      func() time.Time { return time.Now().UTC() },
		sleep:    time.Sleep,
	}
	for _, opt := range opts {
		opt(n)
	}
	if n.httpClient == nil {
		n.httpClient = &http.Client{Timeout: n.timeout}
	}
	return n
}

func (n *NPMTrending) Name() string { return npmTrendingName }

func (n *NPMTrending) Fetch(ctx context.Context) ([]schema.Signal, error) {
	now := n.now()
	cp := cachePath(n.cacheDir, npmTrendingName, now)
	if cached, ok, err := loadCache(cp); err != nil {
		return nil, fmt.Errorf("read cache: %w", err)
	} else if ok {
		return cached, nil
	}

	cfg := defaultRetryConfig(n.sleep)
	sigs, err := retryFetch(ctx, cfg, npmTrendingName, func() ([]schema.Signal, error) {
		return n.fetchOnce(ctx, now)
	})
	if err != nil {
		return nil, err
	}

	if err := saveCache(cp, sigs); err != nil {
		return nil, fmt.Errorf("write cache: %w", err)
	}
	return sigs, nil
}

func (n *NPMTrending) fetchOnce(ctx context.Context, now time.Time) ([]schema.Signal, error) {
	u, err := url.Parse(n.baseURL + "/-/v1/search")
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}
	qs := u.Query()
	qs.Set("text", n.buildQuery())
	qs.Set("popularity", "1.0")
	qs.Set("size", fmt.Sprintf("%d", n.size))
	u.RawQuery = qs.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, &httpStatusError{Source: npmTrendingName, Status: resp.StatusCode, Body: string(body)}
	}

	var payload npmSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	sigs := make([]schema.Signal, 0, len(payload.Objects))
	for _, obj := range payload.Objects {
		sigs = append(sigs, obj.toSignal(now))
	}
	return sigs, nil
}

func (n *NPMTrending) buildQuery() string {
	keys := strings.Split(n.keywords, ",")
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		parts = append(parts, "keywords:"+k)
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " ")
}

type npmSearchResponse struct {
	Objects []npmSearchObject `json:"objects"`
}

type npmSearchObject struct {
	Package npmPackage `json:"package"`
	Score   npmScore   `json:"score"`
}

type npmPackage struct {
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Links       npmLinks  `json:"links"`
	Keywords    []string  `json:"keywords"`
}

type npmLinks struct {
	NPM string `json:"npm"`
}

type npmScore struct {
	Final float64 `json:"final"`
}

func (o npmSearchObject) toSignal(now time.Time) schema.Signal {
	u := o.Package.Links.NPM
	if u == "" {
		u = "https://www.npmjs.com/package/" + o.Package.Name
	}
	return schema.Signal{
		SchemaVersion:    schema.SignalSchemaVersion,
		ID:               signalID(npmTrendingName, u),
		Source:           npmTrendingName,
		URL:              u,
		Title:            o.Package.Name,
		Summary:          o.Package.Description,
		CapturedAt:       now,
		Language:         "JavaScript",
		Stars:            int(math.Round(o.Score.Final * 1000)),
		TopicTags:        o.Package.Keywords,
		MaintenanceScore: o.Score.Final,
		LastCommit:       o.Package.Date.UTC(),
	}
}

var _ Source = (*NPMTrending)(nil)
