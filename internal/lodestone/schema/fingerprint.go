package schema

import "time"

const FingerprintSchemaVersion = 1

type Fingerprint struct {
	SchemaVersion  int               `json:"schema_version"`
	GeneratedAt    time.Time         `json:"generated_at"`
	Languages      []string          `json:"languages,omitempty"`
	Frameworks     []string          `json:"frameworks,omitempty"`
	Deps           map[string]string `json:"deps,omitempty"`
	LOCPerLanguage map[string]int    `json:"loc_per_language,omitempty"`
	TestRatio      float64           `json:"test_ratio,omitempty"`
	HasCI          bool              `json:"has_ci,omitempty"`
	CIProvider     string            `json:"ci_provider,omitempty"`
	MCPServers     []string          `json:"mcp_servers,omitempty"`
	Goals          []string          `json:"goals,omitempty"`
	TechInterests  []string          `json:"tech_interests,omitempty"`
}
