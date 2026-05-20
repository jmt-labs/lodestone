package store

import (
	"reflect"
	"testing"
	"time"

	"github.com/jmt-labs/lodestone/internal/lodestone/schema"
)

func newTestStore(t *testing.T) *FileStore {
	t.Helper()
	s, err := New(t.TempDir())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func TestSignalAppendHasListSince(t *testing.T) {
	s := newTestStore(t)
	now := time.Now().UTC().Truncate(time.Second)

	a := schema.Signal{
		SchemaVersion: schema.SignalSchemaVersion,
		ID:            "sig-a",
		Source:        "github_trending",
		URL:           "https://example.test/a",
		Title:         "A",
		CapturedAt:    now,
	}

	if err := s.Append(a); err != nil {
		t.Fatalf("append: %v", err)
	}
	if ok, err := s.Has("sig-a"); err != nil || !ok {
		t.Fatalf("Has(sig-a)=%v,%v", ok, err)
	}
	if ok, err := s.Has("sig-missing"); err != nil || ok {
		t.Fatalf("Has(sig-missing)=%v,%v", ok, err)
	}

	if err := s.Append(a); err != nil {
		t.Fatalf("append idempotent: %v", err)
	}

	got, err := s.ListSince(time.Time{})
	if err != nil {
		t.Fatalf("ListSince: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 signal after idempotent append, got %d", len(got))
	}
	if !reflect.DeepEqual(got[0], a) {
		t.Fatalf("roundtrip mismatch:\nwant=%+v\ngot =%+v", a, got[0])
	}

	older := schema.Signal{
		SchemaVersion: schema.SignalSchemaVersion,
		ID:            "sig-old",
		Source:        "hackernews",
		URL:           "https://example.test/old",
		Title:         "Old",
		CapturedAt:    now.Add(-72 * time.Hour),
	}
	if err := s.Append(older); err != nil {
		t.Fatalf("append older: %v", err)
	}

	got, err = s.ListSince(now.Add(-1 * time.Hour))
	if err != nil {
		t.Fatalf("ListSince filtered: %v", err)
	}
	if len(got) != 1 || got[0].ID != "sig-a" {
		t.Fatalf("expected only sig-a after filter, got %+v", got)
	}
}

func TestFingerprintWriteRead(t *testing.T) {
	s := newTestStore(t)
	now := time.Now().UTC().Truncate(time.Second)

	in := schema.Fingerprint{
		SchemaVersion: schema.FingerprintSchemaVersion,
		GeneratedAt:   now,
		Languages:     []string{"Go"},
		Frameworks:    []string{"cobra"},
		Deps:          map[string]string{"github.com/spf13/cobra": "v1.10.2"},
		Goals:         []string{"reliability"},
	}

	if err := s.Write(in); err != nil {
		t.Fatalf("Write: %v", err)
	}
	out, err := s.Read()
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if !reflect.DeepEqual(in, out) {
		t.Fatalf("fingerprint mismatch:\nwant=%+v\ngot =%+v", in, out)
	}
}

func TestRecommendationReplaceList(t *testing.T) {
	s := newTestStore(t)

	r1 := schema.Recommendation{
		SchemaVersion: schema.RecommendationSchemaVersion,
		ID:            "rec-1",
		SignalID:      "sig-a",
		Compatibility: 0.9,
		Effort:        schema.EffortS,
		Risk:          schema.RiskLow,
	}
	r2 := schema.Recommendation{
		SchemaVersion: schema.RecommendationSchemaVersion,
		ID:            "rec-2",
		SignalID:      "sig-b",
		Compatibility: 0.5,
		Effort:        schema.EffortM,
		Risk:          schema.RiskMed,
	}

	if err := s.Replace([]schema.Recommendation{r1, r2}); err != nil {
		t.Fatalf("replace: %v", err)
	}
	got, err := s.List()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if !reflect.DeepEqual(got, []schema.Recommendation{r1, r2}) {
		t.Fatalf("list mismatch:\nwant=%+v\ngot =%+v", []schema.Recommendation{r1, r2}, got)
	}

	r3 := schema.Recommendation{
		SchemaVersion: schema.RecommendationSchemaVersion,
		ID:            "rec-3",
		SignalID:      "sig-c",
		Compatibility: 0.7,
		Effort:        schema.EffortL,
		Risk:          schema.RiskHigh,
	}
	if err := s.Replace([]schema.Recommendation{r3}); err != nil {
		t.Fatalf("replace truncate: %v", err)
	}
	got, err = s.List()
	if err != nil {
		t.Fatalf("list after truncate: %v", err)
	}
	if !reflect.DeepEqual(got, []schema.Recommendation{r3}) {
		t.Fatalf("truncate mismatch:\nwant=%+v\ngot =%+v", []schema.Recommendation{r3}, got)
	}
}

func TestIndexRebuildAcrossInstances(t *testing.T) {
	dir := t.TempDir()
	first, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	sig := schema.Signal{
		SchemaVersion: schema.SignalSchemaVersion,
		ID:            "persisted",
		Source:        "github_trending",
		URL:           "https://example.test/p",
		Title:         "P",
		CapturedAt:    time.Now().UTC().Truncate(time.Second),
	}
	if err := first.Append(sig); err != nil {
		t.Fatalf("first append: %v", err)
	}

	second, err := New(dir)
	if err != nil {
		t.Fatalf("re-open: %v", err)
	}
	ok, err := second.Has("persisted")
	if err != nil || !ok {
		t.Fatalf("Has after re-open=%v,%v", ok, err)
	}
}
