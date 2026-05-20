package apply

import (
	"fmt"
	"time"

	"github.com/jmt-labs/lodestone/internal/lodestone/schema"
)

const (
	MinCompatibility = 0.85
	RateLimitWindow  = 24 * time.Hour
)

type GateViolation struct {
	Gate    string
	Message string
}

func (v GateViolation) Error() string {
	return fmt.Sprintf("safety gate %q: %s", v.Gate, v.Message)
}

type GateResult struct {
	Passed     bool
	Violations []GateViolation
}

func CheckRecommendation(rec schema.Recommendation) GateResult {
	r := GateResult{Passed: true}
	if rec.Risk != schema.RiskLow {
		r.Passed = false
		r.Violations = append(r.Violations, GateViolation{Gate: "risk", Message: fmt.Sprintf("risk = %q, want low", rec.Risk)})
	}
	if rec.Effort != schema.EffortXS {
		r.Passed = false
		r.Violations = append(r.Violations, GateViolation{Gate: "effort", Message: fmt.Sprintf("effort = %q, want XS", rec.Effort)})
	}
	if rec.Compatibility < MinCompatibility {
		r.Passed = false
		r.Violations = append(r.Violations, GateViolation{Gate: "compatibility", Message: fmt.Sprintf("compatibility = %.3f, want >= %.2f", rec.Compatibility, MinCompatibility)})
	}
	return r
}

func CheckRateLimit(applies []Apply, now time.Time) GateResult {
	cutoff := now.Add(-RateLimitWindow)
	for _, a := range applies {
		if a.AppliedAt.After(cutoff) {
			return GateResult{
				Passed: false,
				Violations: []GateViolation{{
					Gate:    "rate_limit",
					Message: fmt.Sprintf("apply %q from %s is within %s", a.RecID, a.AppliedAt.Format(time.RFC3339), RateLimitWindow),
				}},
			}
		}
	}
	return GateResult{Passed: true}
}

func CheckCleanGit(status string) GateResult {
	if status == "" {
		return GateResult{Passed: true}
	}
	return GateResult{
		Passed: false,
		Violations: []GateViolation{{
			Gate:    "git_clean",
			Message: fmt.Sprintf("git status is dirty:\n%s", status),
		}},
	}
}
