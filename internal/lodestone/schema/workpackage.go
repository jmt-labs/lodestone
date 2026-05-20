package schema

type WorkPackage struct {
	ID                 string   `json:"id"`
	Type               string   `json:"type"`
	Title              string   `json:"title"`
	DependsOn          []string `json:"depends_on,omitempty"`
	FilesAffected      []string `json:"files_affected,omitempty"`
	ExpectedArtifacts  []string `json:"expected_artifacts,omitempty"`
	Executor           string   `json:"executor,omitempty"`
	EstimatedMinutes   int      `json:"estimated_minutes,omitempty"`
	AcceptanceCriteria []string `json:"acceptance_criteria,omitempty"`
}
