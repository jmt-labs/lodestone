# Lodestone Phase 1 (MVP) — Implementierungsplan

**Datum:** 2026-05-20
**Spec:** [`2026-05-20-lodestone-design.md`](../specs/2026-05-20-lodestone-design.md)
**Quelle:** Übernommen aus
`jmt-labs/forgecrate@claude/ai-trend-intelligence-evolution-un74G/
docs/superpowers/plans/2026-05-19-lodestone-mvp.md`, an Standalone-
Layout angepasst.

## Übersicht

Deterministische Phase-1-Pipeline:

```
fingerprint  → .lodestone/fingerprint.json
ingest       → .lodestone/signals.jsonl
score        → .lodestone/recommendations.jsonl
signals      → Lese-Anfrage über signals.jsonl
```

Fünf Go-Pakete unter `internal/lodestone/`: `schema`, `store`, `ingest`,
`fingerprint`, `scoring`. Plus `internal/config` für `.lodestone.yaml`.
**Keine** externen Go-Deps außer Cobra und yaml.v3.

## Architekturentscheidungen

- **Keine LLM-Aufrufe** in Phase 1.
- **Deterministische Reihenfolge** in `score`: `compatibility DESC,
  stars DESC, id ASC`.
- **Compatibility:** gewichtete Jaccard (Topic-Tags ∪ Language gegen
  Frameworks ∪ Languages); Language-Match-Faktor **1.5**, Framework-
  Match-Faktor **1.0**.
- **Effort-Heuristik:** Default `M`; `S` bei Topic-Overlap und niedrigen
  Stars (<100); `XL` bei 0 Match.
- **Risk:** `low` bei Stars≥500 ∧ LastCommit<90d ∧ Lizenz; `high`
  bei Bedingungsverletzung; sonst `med`.

## Tasks

### T1 — Schemas

**Dateien:**
- `internal/lodestone/schema/signal.go`
- `internal/lodestone/schema/fingerprint.go`
- `internal/lodestone/schema/recommendation.go`
- `internal/lodestone/schema/workpackage.go`
- `internal/lodestone/schema/schema_test.go`

**Akzeptanz:**
- `Signal{SchemaVersion, ID, Source, URL, Title, Summary, CapturedAt,
  Language, Stars, TopicTags, MaintenanceScore, License, LastCommit}`.
- `Fingerprint{SchemaVersion, GeneratedAt, Languages, Frameworks, Deps,
  LOCPerLanguage, TestRatio, HasCI, CIProvider, MCPServers, Goals,
  TechInterests}`.
- `Recommendation{SchemaVersion, ID, SignalID, Compatibility, Effort,
  ROI, Risk, Rationale, CounterEvidence, SuggestedNext}`.
- `WorkPackage{ID, Type, Title, DependsOn, FilesAffected,
  ExpectedArtifacts, Executor, EstimatedMinutes, AcceptanceCriteria}`.
- JSON-Roundtrip-Tests für alle vier Typen.
- `go test ./internal/lodestone/schema/...` grün.
- Commit: `feat(lodestone): schemas for Signal, Fingerprint, Recommendation, WorkPackage`.

**Abhängigkeiten:** keine.

---

### T2 — Store-Interface + FileStore

**Dateien:**
- `internal/lodestone/store/store.go`
- `internal/lodestone/store/filestore.go`
- `internal/lodestone/store/filestore_test.go`

**Akzeptanz:**
- `SignalStore`, `FingerprintStore`, `RecommendationStore` Interfaces.
- `FileStore` implementiert alle drei unter `.lodestone/`.
- Signals: JSONL, In-Memory-Set für `Has()`, ListSince streamend.
- Fingerprint: einzelne JSON-Datei, atomar via tmp+rename.
- Recommendations: JSONL, atomares Replace.
- Tmpdir-Tests via `t.TempDir()`.
- `go test ./internal/lodestone/store/...` grün.
- Commit: `feat(lodestone): file-based store with JSONL signals + JSON fingerprint`.

**Abhängigkeiten:** T1.

---

### T3 — Ingest-Interface + GitHub-Trending-Adapter

**Dateien:**
- `internal/lodestone/ingest/source.go`
- `internal/lodestone/ingest/github_trending.go`
- `internal/lodestone/ingest/github_trending_test.go`

**Akzeptanz:**
- `Source` Interface mit `Name()` und `Fetch(ctx) ([]schema.Signal, error)`.
- GitHub Search-API (`api.github.com/search/repositories`) mit
  `$GITHUB_TOKEN` (optional).
- Deterministische ID: `sha256("github_trending|" + html_url)`.
- Cache: `.lodestone/cache/github_trending-<YYYY-MM-DD>.json`.
- 15-Sekunden-Timeout, exponentieller Backoff-Retry (max 3 Versuche).
- `httptest.Server`-Tests für Success/Timeout/Empty.
- Commit: `feat(lodestone): ingest interface + github trending source`.

**Abhängigkeiten:** T1.

---

### T4 — HackerNews-Adapter

**Dateien:**
- `internal/lodestone/ingest/hackernews.go`
- `internal/lodestone/ingest/hackernews_test.go`

**Akzeptanz:**
- Firebase-API (`hacker-news.firebaseio.com/v0/topstories.json` +
  `/item/<id>.json`).
- Story-Type-Filter + Keyword-Filter (Default: `ai, llm, mcp, claude,
  agent`).
- Limit: 50 Items pro Lauf.
- Cache analog zu T3.
- `httptest.Server`-Tests.
- Commit: `feat(lodestone): hackernews source`.

**Abhängigkeiten:** T3.

---

### T5 — Fingerprint (Go + Node)

**Dateien:**
- `internal/lodestone/fingerprint/fingerprint.go`
- `internal/lodestone/fingerprint/golang.go`
- `internal/lodestone/fingerprint/node.go`
- `internal/lodestone/fingerprint/fingerprint_test.go`
- `internal/lodestone/fingerprint/testdata/go_minimal/go.mod`
- `internal/lodestone/fingerprint/testdata/node_react/package.json`

**Akzeptanz:**
- Walker, der je Sprache passende Detektoren aufruft.
- LOC-Counting mit `vendor/`, `node_modules/`, `.git/` Skip.
- Test-Ratio = (Test-LOC / Non-Test-LOC).
- Go: `go.mod` per Regex parsen (**keine** neuen Deps).
- Node: `package.json` via `encoding/json`.
- Framework-Heuristik: react, vue, next, `@anthropic-ai/sdk`, cobra.
- Goldens prüfen `languages`, `frameworks`, `deps`.
- Commit: `feat(lodestone): fingerprint for Go + Node (regex go.mod parse, package.json)`.

**Abhängigkeiten:** T1.

---

### T6 — Scoring

**Dateien:**
- `internal/lodestone/scoring/compatibility.go`
- `internal/lodestone/scoring/effort.go`
- `internal/lodestone/scoring/risk.go`
- `internal/lodestone/scoring/scoring.go`
- `internal/lodestone/scoring/scoring_test.go`

**Akzeptanz:**
- Compatibility wie in Spec; Sortierung deterministisch.
- Effort, Risk wie in Spec.
- Rec-ID: `sha256(signal_id + fp.canonical())`.
- 5 Fixture-Tests + Determinismus-Verifikation (2 Läufe).
- Commit: `feat(lodestone): deterministic scoring (compatibility, effort, risk)`.

**Abhängigkeiten:** T1, T5.

---

### T7 — Config-Erweiterung

**Dateien:**
- `internal/config/config.go`
- `internal/config/config_test.go`

**Akzeptanz:**
- `Goals []string`, `TechInterests []string` (mit `omitempty`).
- `LodestoneConfig{MinStars, MinAgeDays, MaxLastCommitAgeDays,
  RequireLicense}`.
- Defaults: 50, 30, 180, true.
- yaml.v3 für Parsing.
- Commit: `feat(config): goals, tech_interests, lodestone block`.

**Abhängigkeiten:** keine (parallel zu T2 möglich).

---

### T8 — Subkommandos (ingest/fingerprint/score/signals)

**Dateien:**
- `cmd/lodestone/main.go` (Stubs durch echte Implementierungen ersetzen)
- `cmd/lodestone/ingest.go`
- `cmd/lodestone/fingerprint.go`
- `cmd/lodestone/score.go`
- `cmd/lodestone/signals.go`

**Akzeptanz:**
- `ingest --source <name>` (mehrfach erlaubt).
- `fingerprint` ohne Flags.
- `score` ohne Flags.
- `signals --since <date> --source <name> --top <N> --json`.
- Smoke-Test: `go run ./cmd/lodestone fingerprint` schreibt
  `.lodestone/fingerprint.json` im aktuellen Repo.
- Commit: `feat(cmd): lodestone subcommands (ingest, fingerprint, score, signals)`.

**Abhängigkeiten:** T2, T3, T4, T5, T6.

---

### T9 — E2E-Smoke-Test

**Dateien:**
- `e2e/lodestone_test.sh`
- `Makefile` (Target `e2e` aktivieren)

**Akzeptanz:**
- `set -euo pipefail`, `tmpdir` via `mktemp -d`.
- Schritte: `git init`, `go mod init`, `lodestone fingerprint`,
  `lodestone ingest --source github_trending --mock`,
  `lodestone score`, Existenz-Checks für `.lodestone/*.json(l)`.
- `--mock` respektiert `LODESTONE_MOCK_FIXTURES` für Offline-Lauf.
- `make e2e` grün.
- Commit: `test(lodestone): e2e smoke test (fingerprint → ingest → score)`.

**Abhängigkeiten:** T8.

---

### T10 — README + CHANGELOG + docs/lodestone.md

**Akzeptanz:**
- README-Subkommando-Tabelle aktualisiert (alle Phase-1-Verben
  „Implemented").
- `docs/lodestone.md` füllt User-Guide-Sektionen aus.
- CHANGELOG-Eintrag: `Added: lodestone Phase 1 MVP`.
- README-Coverage-Check im CI grün.
- Commit: `docs(lodestone): user-facing overview + README + CHANGELOG`.

**Abhängigkeiten:** T8.

---

### T11 — Backward-Compat- & Determinismus-Verifikation

**Akzeptanz:**
- Determinismus: zwei Score-Läufe mit identischem Input → byte-
  identische `recommendations.jsonl`.
- Alle bestehenden Tests grün: `go test ./...`, `make e2e`,
  `golangci-lint run`, `go vet ./...`, `govulncheck ./...`.
- Diff-Dokumentation bei Abweichungen.

**Abhängigkeiten:** T1 – T10.

---

### T12 — Release-Vorbereitung

**Akzeptanz:**
- Tag `v0.1.0-alpha` lokal erzeugt, nicht gepusht.
- `goreleaser release --snapshot --clean` lokal grün.
- PR-Erstellung **nur auf explizite Aufforderung**.

**Abhängigkeiten:** T1 – T11.

## Test-Strategie

- Unit: jeder Sub-Paket-Test isoliert via `t.TempDir()`.
- E2E: `e2e/lodestone_test.sh` mit `LODESTONE_MOCK_FIXTURES`.
- Fingerprint-Goldens: deterministischer Vergleich.
- Determinismus-Spezial-Test in T6 und T11.

## Exit-Kriterien

- `lodestone score --json` liefert deterministisch sortierten Output.
- `e2e/lodestone_test.sh` läuft offline mit Mock-Fixtures.
- CI grün auf `main` (`test`, `lint`, `vuln`, `readme-coverage`).
- README dokumentiert alle Phase-1-Verben.
