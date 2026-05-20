# Lodestone — AI-Trend-Intelligence-System (Design-Spec)

**Datum:** 2026-05-20
**Status:** Aktiv, Phase 1 in Umsetzung
**Quelle:** Übernommen aus forgecrate-Spec
`2026-05-19-lodestone-design.md`, angepasst auf Standalone-Layout.

## Kernziel

Externe AI-Ökosystem-Signale (News, GitHub-Trending, MCP-Registry,
Framework-Releases, ArXiv, Anthropic/OpenAI-Changelogs, npm/PyPI,
awesome-Listen) in **kontextualisierte Implementierungspläne** für
konkrete Repositories umwandeln — als eigenständiges, opt-in nutzbares
Werkzeug, das keinerlei Annahmen über das Zielprojekt macht außer
„hat ein Git-Repo".

## Produktname & Metapher

**Lodestone** = Magnetstein/Kompass. Der Name harmoniert mit
handwerklichem Vokabular, suggeriert Richtung ohne Vorschrift, und ist
bewusst nicht-hype.

**Tagline:** „Liest das AI-Ökosystem für dein Repo."

## Hybrid-Integrations-Modell (drei Ebenen)

1. **CLI-Binary `lodestone`** (Cobra, Go) — deterministische Pipeline für
   Ingestion, Caching, Fingerprint und Scoring.
2. **Claude-Skills im `flavors/`-Stil** (Phase 2) — direkt in beliebige
   Repos installierbar; rufen das CLI via Bash auf.
3. **MCP-Server `lodestone-mcp`** (Phase 3, separates Go-Binary) — dünne
   Lese-Schicht mit Tools `list_signals`, `query_trends`, `score_repo`,
   `generate_plan`, `record_decision`.

**Begründung:** Skills brauchen Einstiegspunkte für Claude; Batch-
Processing gehört in Go (deterministisch, debuggbar, kein LLM); MCP
bleibt dünn.

## Daten-Pipeline & Datei-Struktur

**Stufen:**

```
ingest      → .lodestone/signals.jsonl
fingerprint → .lodestone/fingerprint.json
score       → .lodestone/recommendations.jsonl
plan <id>   → docs/superpowers/specs/YYYY-MM-DD-<slug>-design.md
              + docs/superpowers/plans/YYYY-MM-DD-<slug>.md
```

**Lokale Artefakte unter `.lodestone/`:**

- `cache/` — Rohfetches mit TTL-Datum (`<source>-YYYY-MM-DD.json`)
- `signals.jsonl` — normalisiert, dedupliziert
- `signals.idx` — Offset-Index (lazy rebuild)
- `fingerprint.json` — letzter Repo-Scan
- `recommendations.jsonl` — gescorte Vorschläge
- `decisions.log` — Audit-Trail (einzige Datei, die **nicht** in
  `.gitignore` aufgenommen wird; Skill-Installer fügt entsprechenden
  Snippet in `.gitignore` des Zielprojekts ein)

Jede Stufe ist isoliert re-runnbar und debuggbar.

## Trend-Discovery: Quellen & Anti-Hype

**Sieben Quellen** (Phase 1: zwei davon; Phase 2: vier weitere; Phase 3:
MCP-Registry-Pull):

- GitHub Trending (Phase 1)
- HackerNews (Phase 1)
- ArXiv (Phase 2)
- Anthropic / OpenAI Changelogs (Phase 2)
- npm / PyPI Trending (Phase 2)
- Awesome-Listen (Phase 2)
- MCP Registry (Phase 3)

**Source-Interface:**

```go
type Source interface {
    Name() string
    Fetch(ctx context.Context) ([]RawSignal, error)
}
```

**Anti-Hype-Defaults** (in `.lodestone.yaml` konfigurierbar):

```yaml
lodestone:
  min_stars: 50
  min_age_days: 30
  max_last_commit_age_days: 180
  require_license: true
```

**Scheduling:** Default manuell. Optional ein GitHub-Action-Template mit
wöchentlichem Cron (Phase 3). Kein Daemon.

## Repo-Fingerprint (Go, deterministisch)

Erfasst:

- **Sprachen** mit LOC-Anteilen
- **Frameworks** (heuristisch erkannt: react, vue, next, anthropic-sdk,
  cobra, fastapi, …)
- **Dependencies** mit Versionen (aus `go.mod`, `package.json`,
  `requirements.txt`, `pyproject.toml`, …)
- **Metriken:** Test-Quote (Test-LOC / Non-Test-LOC), CI-Provider
  (GitHub Actions, GitLab CI, CircleCI, …)
- **MCP-Server** (aus `.mcp.json`, falls vorhanden)
- **Goals[]** und **TechInterests[]** aus `.lodestone.yaml`

Fehlen `goals`: interaktive Skill-Abfrage beim ersten Lauf (Phase 2).

## Scoring-Dimensionen

| Dimension | Werte | Berechnung |
|---|---|---|
| `compatibility` | 0.0 – 1.0 | Gewichtete Jaccard-Ähnlichkeit Signal-Tags ∩ Frameworks/Sprachen; Sprach-Match 1.5×, Framework-Match 1.0× |
| `effort` | XS – XL | Heuristik: neue Dependencies, neue Dateien, Stars |
| `roi` | low / med / high | Abbildung aus `goals[]`-Treffern + Kompatibilität (ab Phase 2) |
| `risk` | low / med / high | Stars, Wartung (LastCommit-Alter), Lizenz, (Phase 2) CVE-Listen |

**Explanation-Layer (ab Phase 2):** Jede Empfehlung enthält `rationale`
(3 Sätze) und `counter_evidence` (1 Satz) vom `mechanical`-Modell.

**Anzeige-Schwelle:** Compatibility ≥ 0.4 sichtbar; ≥ 0.7 als
„Empfehlung" markiert.

**Deterministische Sortierung:** `compatibility DESC, stars DESC,
id ASC`. Zwei Läufe mit identischem Input müssen byte-identische
Outputs erzeugen.

## Planning-Engine (Phase 2)

Hierarchie `Recommendation → Epic → Story → Task → Subtask`:

1. `lodestone plan <rec-id>` triggert Skill `lodestone-plan`.
2. Go lädt Fingerprint + Recommendation, übergibt Kontext an Claude.
3. Skill nutzt `superpowers:writing-plans` und erzeugt Spec/Plan/
   WorkPackages-YAML.

**WorkPackage-Schema (Beispiel):**

```yaml
- id: WP-001
  type: task
  title: "Add fetch adapter for HackerNews"
  files_affected:
    - internal/lodestone/ingest/hackernews.go
  executor: developer
  estimated_minutes: 45
  acceptance_criteria:
    - "Unit-Test deckt Fetch & Timeout ab"
```

## Agent-System

| Agent | Typ | Modell | Output |
|---|---|---|---|
| Trend Scout | Go | — | `signals.jsonl` |
| Repo Analyzer | Go | — | `fingerprint.json` |
| Compatibility Scorer | Go (+ LLM ab Phase 2) | `mechanical` | `recommendations.jsonl` |
| Recommendation Skill | Skill | `default` | aufbereitete Liste |
| Planning Engine | Go + LLM | `planning` | Spec/Plan/Tasks |

**Memory:** Nur Entscheidungen + Architektur-Begründungen →
`.claude/memory.json`. Rohdaten bleiben in `.lodestone/`.

## Phase-Roadmap

- **Phase 1 (MVP, ~2 Wochen):** 2 Quellen, Fingerprint Go+Node,
  deterministisches Scoring, JSONL-Store, CLI-Subkommandos
  `ingest/fingerprint/score/signals`, E2E-Smoke-Test.
- **Phase 2 (~4 Wochen):** 4 weitere Quellen, 4 Skills
  (`lodestone-scout`, `-recommend`, `-plan`, `-review-trends`),
  Plan-Generator mit `superpowers:writing-plans`, Goals-Block in
  Config, `.gitignore`-Snippet bei Skill-Install.
- **Phase 3 (~6 Wochen):** `lodestone-mcp` Binary, GitHub-Action-
  Template für scheduled Ingest (wöchentlich, Branch-PR),
  Memory-Persistierung von Decisions.
- **Phase 4 (TBD):** Auto-PR-Engine mit harten Schranken
  (`risk: low` ∧ `effort: XS` ∧ `compatibility ≥ 0.85`, max 1
  Auto-PR/Tag/Repo, niemals auf `main`, `--draft` als Default),
  Success-Tracker, Quartals-Re-Scoring als Vorschlag-PR,
  Cross-Repo-Sharing (erfordert eigenen Privacy-Spec).

**Phase 4 explizit nicht blockierend:** Auto-PRs nur mit strengen
Kriterien, Rollback via `lodestone undo <pr>`, Selbstverbesserung
über Vorschlag-PRs.

## Tradeoffs

- **Build vs. Buy:** Kein Fertig-Tool für AI/MCP-Ökosystem-Scanning →
  Build.
- **Local vs. External:** Lokale Artefakte, keine Telemetrie ohne
  Opt-In. LLM-Aufrufe erst ab Phase 2, dann ebenfalls opt-in.
- **Pull vs. Push:** Pull-Default; Push via Action-Template opt-in.
- **Storage:** JSONL + Files im MVP; SQLite/bbolt erst bei
  Skalierungs-Nachweis.

## Nicht in Phase 1

LLM-Planning, MCP-Server, Auto-PR, Cross-Repo-Sharing, Marketplace,
Awesome-List-Ingestion.

## Testbarkeit

- Unit-Tests ≥ 70 % Coverage in `internal/`.
- E2E: `e2e/lodestone_test.sh` — `fingerprint → ingest → score` mit
  Mock-Fixtures (`LODESTONE_MOCK_FIXTURES`).
- Fingerprint-Goldens in `internal/lodestone/fingerprint/testdata/`.
- Schema-Validierung für alle JSON-Outputs.
- **Determinismus-Verifikation** in T11: zwei Score-Läufe mit
  identischem Input erzeugen byte-identische sortierte Liste.

## Betroffene Dateien (Phase 1)

- `cmd/lodestone/main.go` — Root-Cobra + Subkommandos
- `cmd/lodestone/<verb>.go` — pro Verb eine Datei (T8)
- `internal/lodestone/schema/**` — T1
- `internal/lodestone/store/**` — T2
- `internal/lodestone/ingest/**` — T3, T4
- `internal/lodestone/fingerprint/**` — T5
- `internal/lodestone/scoring/**` — T6
- `internal/config/config.go` — T7 (Goals, TechInterests,
  LodestoneConfig)
- `e2e/lodestone_test.sh` — T9
- `README.md`, `CHANGELOG.md`, `docs/lodestone.md` — T10

## Migration-Snippet aus forgecrate

Diese Spec basiert auf der ursprünglichen Vorlage aus
`jmt-labs/forgecrate@claude/ai-trend-intelligence-evolution-un74G/
docs/superpowers/specs/2026-05-19-lodestone-design.md`. Wesentliche
Anpassungen für den Standalone-Kontext:

| forgecrate | lodestone (standalone) |
|---|---|
| `.forgecrate/lodestone/…` | `.lodestone/…` |
| `forgecrate lodestone <verb>` | `lodestone <verb>` |
| `cmd/forgecrate/lodestone.go` | `cmd/lodestone/main.go` (+ Verb-Dateien) |
| `.forgecrate.yaml` (`lodestone:`-Block) | `.lodestone.yaml` (Root-Level) |
| forgecrate-Flavor-System | Standalone-Skills (Phase 2) |
