package scoring

import (
	"time"

	"github.com/jmt-labs/lodestone/internal/lodestone/schema"
)

const (
	riskPopularStars = 500
	riskFreshDays    = 90
	riskStaleDays    = 180
)

func Risk(sig schema.Signal, now time.Time) schema.RiskLevel {
	hasLicense := sig.License != ""
	popular := sig.Stars >= riskPopularStars

	var fresh, stale bool
	if !sig.LastCommit.IsZero() {
		days := int(now.Sub(sig.LastCommit).Hours() / 24)
		fresh = days >= 0 && days < riskFreshDays
		stale = days > riskStaleDays
	}

	if popular && fresh && hasLicense {
		return schema.RiskLow
	}
	if !hasLicense || stale {
		return schema.RiskHigh
	}
	return schema.RiskMed
}
