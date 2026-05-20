package schema

import "time"

const SignalSchemaVersion = 1

type Signal struct {
	SchemaVersion    int       `json:"schema_version"`
	ID               string    `json:"id"`
	Source           string    `json:"source"`
	URL              string    `json:"url"`
	Title            string    `json:"title"`
	Summary          string    `json:"summary,omitempty"`
	CapturedAt       time.Time `json:"captured_at"`
	Language         string    `json:"language,omitempty"`
	Stars            int       `json:"stars,omitempty"`
	TopicTags        []string  `json:"topic_tags,omitempty"`
	MaintenanceScore float64   `json:"maintenance_score,omitempty"`
	License          string    `json:"license,omitempty"`
	LastCommit       time.Time `json:"last_commit,omitempty"`
}
