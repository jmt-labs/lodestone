package scoring

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/jmt-labs/lodestone/internal/lodestone/schema"
)

type scorer struct {
	now func() time.Time
}

type Option func(*scorer)

func WithNow(fn func() time.Time) Option {
	return func(s *scorer) { s.now = fn }
}

func Score(fp schema.Fingerprint, sigs []schema.Signal, opts ...Option) ([]schema.Recommendation, error) {
	s := &scorer{now: func() time.Time { return time.Now().UTC() }}
	for _, o := range opts {
		o(s)
	}

	canonical, err := canonicalFingerprint(fp)
	if err != nil {
		return nil, fmt.Errorf("canonical fingerprint: %w", err)
	}

	type entry struct {
		rec   schema.Recommendation
		stars int
	}

	now := s.now()
	entries := make([]entry, 0, len(sigs))
	for _, sig := range sigs {
		compat := Compatibility(sig, fp)
		rec := schema.Recommendation{
			SchemaVersion: schema.RecommendationSchemaVersion,
			ID:            recommendationID(sig.ID, canonical),
			SignalID:      sig.ID,
			Compatibility: compat,
			Effort:        Effort(sig, compat),
			Risk:          Risk(sig, now),
		}
		entries = append(entries, entry{rec: rec, stars: sig.Stars})
	}

	sort.SliceStable(entries, func(i, j int) bool {
		a, b := entries[i], entries[j]
		if a.rec.Compatibility != b.rec.Compatibility {
			return a.rec.Compatibility > b.rec.Compatibility
		}
		if a.stars != b.stars {
			return a.stars > b.stars
		}
		return a.rec.ID < b.rec.ID
	})

	recs := make([]schema.Recommendation, len(entries))
	for i, e := range entries {
		recs[i] = e.rec
	}
	return recs, nil
}

func canonicalFingerprint(fp schema.Fingerprint) ([]byte, error) {
	return json.Marshal(fp)
}

func recommendationID(signalID string, canonical []byte) string {
	h := sha256.New()
	h.Write([]byte(signalID))
	h.Write([]byte{'|'})
	h.Write(canonical)
	return "sha256:" + hex.EncodeToString(h.Sum(nil))
}
