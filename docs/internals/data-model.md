# Datenmodell

Vier zentrale Schemas leben unter `internal/lodestone/schema/`. Alle
Schemas tragen ein `schema_version`-Feld (heute überall `1`) und sind
über JSON-Roundtrip-Tests in `schema_test.go` abgesichert.

## Signal

Ein atomares Trend-Ereignis aus einer externen Quelle. Eine Zeile in
`.lodestone/signals.jsonl`.

```go
type Signal struct {
    SchemaVersion    int       `json:"schema_version"`
    ID               string    `json:"id"`                // sha256:<hex>
    Source           string    `json:"source"`            // github_trending, hackernews, ...
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
```

**Deduplikation:** Über `ID` (sha256-Hash aus Source + URL). Beim
zweiten `ingest` desselben Signals gewinnt der frühere Eintrag.

## Fingerprint

Strukturiertes Profil des analysierten Repositories. Einzeln in
`.lodestone/fingerprint.json`, atomar via `tmp + rename` geschrieben.

```go
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
```

**Goals und TechInterests** stammen aus `.lodestone.yaml` und sind
das Bindeglied zur Planning-Engine ab Phase 2.

## Recommendation

Eine Empfehlung — Resultat von `lodestone score`. Eine Zeile in
`.lodestone/recommendations.jsonl`.

```go
type Recommendation struct {
    SchemaVersion   int         `json:"schema_version"`
    ID              string      `json:"id"`                // sha256:<hex>
    SignalID        string      `json:"signal_id"`
    Compatibility   float64     `json:"compatibility"`     // 0.0–1.0
    Effort          EffortLevel `json:"effort"`            // XS, S, M, L, XL
    ROI             ROILevel    `json:"roi,omitempty"`     // low, med, high
    Risk            RiskLevel   `json:"risk"`              // low, med, high
    Rationale       string      `json:"rationale,omitempty"`
    CounterEvidence string      `json:"counter_evidence,omitempty"`
    SuggestedNext   []string    `json:"suggested_next,omitempty"`
}
```

**ID-Schema:** `sha256:hex(signal_id + "|" + json(fingerprint))`.
Deterministisch über identische Inputs (siehe
[Determinismus](determinism.md)).

**Sortierung in der JSONL-Datei:**
`compatibility DESC, stars DESC, id ASC`. Stars werden über `SignalID`
aufgelöst.

## WorkPackage

Strukturierter Output der Planning-Engine. Wird heute noch nicht
persistiert — Spec/Plan-Markdown ist der primäre Output. Das Schema
ist für zukünftige LLM-Ausgaben (Phase 5+) vorbereitet.

```go
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
```

## Schema-Versionierung

Heute liegt jedes Schema auf `1`. Breaking Changes erhöhen die
Schema-Version und müssen einen Migrations-Pfad mitliefern (`old → new`
auf Read-Time, neuer Schreib-Path nur in der neuen Version). Pre-1.0.0
gilt: `-alpha`-Releases dürfen Schemas brechen, solange das CHANGELOG
es protokolliert. Ab `v1.0.0` ist die Schema-Version Teil der
öffentlichen API.
