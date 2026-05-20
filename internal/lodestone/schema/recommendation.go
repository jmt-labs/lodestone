package schema

const RecommendationSchemaVersion = 1

type EffortLevel string

const (
	EffortXS EffortLevel = "XS"
	EffortS  EffortLevel = "S"
	EffortM  EffortLevel = "M"
	EffortL  EffortLevel = "L"
	EffortXL EffortLevel = "XL"
)

type ROILevel string

const (
	ROILow  ROILevel = "low"
	ROIMed  ROILevel = "med"
	ROIHigh ROILevel = "high"
)

type RiskLevel string

const (
	RiskLow  RiskLevel = "low"
	RiskMed  RiskLevel = "med"
	RiskHigh RiskLevel = "high"
)

type Recommendation struct {
	SchemaVersion   int         `json:"schema_version"`
	ID              string      `json:"id"`
	SignalID        string      `json:"signal_id"`
	Compatibility   float64     `json:"compatibility"`
	Effort          EffortLevel `json:"effort"`
	ROI             ROILevel    `json:"roi,omitempty"`
	Risk            RiskLevel   `json:"risk"`
	Rationale       string      `json:"rationale,omitempty"`
	CounterEvidence string      `json:"counter_evidence,omitempty"`
	SuggestedNext   []string    `json:"suggested_next,omitempty"`
}
