package ingest

import (
	"context"

	"github.com/jmt-labs/lodestone/internal/lodestone/schema"
)

type Source interface {
	Name() string
	Fetch(ctx context.Context) ([]schema.Signal, error)
}
