package store

import (
	"time"

	"github.com/jmt-labs/lodestone/internal/lodestone/schema"
)

type SignalStore interface {
	Append(s schema.Signal) error
	ListSince(t time.Time) ([]schema.Signal, error)
	Has(id string) (bool, error)
}

type FingerprintStore interface {
	Write(fp schema.Fingerprint) error
	Read() (schema.Fingerprint, error)
}

type RecommendationStore interface {
	Replace(recs []schema.Recommendation) error
	List() ([]schema.Recommendation, error)
}
