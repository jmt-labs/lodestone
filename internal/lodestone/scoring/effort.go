package scoring

import "github.com/jmt-labs/lodestone/internal/lodestone/schema"

const effortLowStarThreshold = 100

func Effort(sig schema.Signal, compat float64) schema.EffortLevel {
	if compat <= 0 {
		return schema.EffortXL
	}
	if sig.Stars < effortLowStarThreshold {
		return schema.EffortS
	}
	return schema.EffortM
}
