# Changelog

Alle nennenswerten Änderungen an diesem Projekt werden in dieser Datei
dokumentiert. Format gemäß [Keep a Changelog](https://keepachangelog.com/de/1.1.0/),
Versionierung nach [Semantic Versioning](https://semver.org/lang/de/).

Aktueller Stand: **Phasen 1–4** sind auf `main` gemerged und CI-grün.
Detaillierter Phasen-Status:
[`docs/internals/roadmap.md`](docs/internals/roadmap.md).

## [Unreleased]

### Added — Phase 4

- **Auto-PR-Engine** (`internal/lodestone/apply/`): `lodestone apply
  <rec-id>` erzeugt einen Draft-PR aus einer Recommendation. Vier
  unabhängige Safety-Gates müssen alle passen:
  `risk == low`, `effort == XS`, `compatibility >= 0.85`, kein Apply
  in den letzten 24h (`.lodestone/applies.jsonl`), `git status` sauber.
  Branch nach Schema `lodestone/apply-<rec-suffix>-<date>`, **nie auf
  `main` committed**.
- **Undo** (`lodestone undo <branch-or-rec-id>`): schließt den
  zugehörigen PR über `gh pr close --delete-branch` und entfernt
  Branch lokal + remote. Apply-State wird auf `undone` gesetzt.
- **Stats** (`lodestone stats`): aggregiert
  `.lodestone/applies.jsonl` nach Status (`draft_open`,
  `branch_pushed_no_pr`, `undone`).
- **Privacy-Spec** für Cross-Repo-Sharing unter
  `docs/superpowers/specs/2026-05-20-lodestone-sharing-privacy.md`:
  legt fest, welche Felder veröffentlichbar sind, k=5-Anonymität für
  Goals/TechInterests, Opt-In-Flow, Re-Identifikations-Schutz —
  **noch nicht implementiert**, Code-Stub wird erst nach Beantwortung
  der offenen Fragen ausgerollt.
- **Pluggable Runner**: `GitRunner` und `PRRunner` Interfaces mit
  `RealGit`/`RealPR` (shell-out zu `git`/`gh`) und Fake-Implementierungen
  für Tests. Apply ist damit vollständig mockfähig.

### Added — Phase 3

- **`lodestone-mcp` Binary** (`cmd/lodestone-mcp/`): zweites Binary,
  spricht MCP über stdio (JSON-RPC 2.0, Protocol-Version
  `2024-11-05`). Implementiert `initialize`, `tools/list`,
  `tools/call`, schweigt bei Notifications.
- **Fünf MCP-Tools** unter `internal/lodestone/mcp/tools.go`:
  `list_signals`, `query_trends`, `score_repo`, `generate_plan`,
  `record_decision` — alle rufen direkt in die Phase-1/2-Packages
  (kein Duplikat-Code).
- **Memory-Konsolidierung** (`internal/lodestone/memory/`,
  `cmd/lodestone/memory.go`): `lodestone memory` aggregiert die
  letzten N Tage (Default 90) aus `.lodestone/decisions.log` nach
  `.claude/memory.json` — idempotent über (Datum,Verb,Summary)-Dedup.
- **GitHub-Action-Template** unter
  `.github/workflows/templates/lodestone-weekly.yml`: opt-in
  reusable Workflow (Sonntag 03:00 UTC + `workflow_dispatch`), läuft
  `fingerprint → ingest → score → memory`, öffnet PR mit Top-5-
  Summary, wenn Diffs vorhanden.
- **Build-System**: `make build` + `.goreleaser.yaml` bauen beide
  Binaries (`lodestone` und `lodestone-mcp`), Archives enthalten
  beide.

### Added — Phase 2

- **Vier neue Ingest-Quellen** in `internal/lodestone/ingest/`:
  `arxiv` (Atom-Feed via `export.arxiv.org/api/query`),
  `anthropic_changelog` und `openai_changelog` (HTML-Scrape mit
  `<h2>`/`<h3>` + `YYYY-MM-DD`-Erkennung),
  `npm_trending` (`registry.npmjs.org/-/v1/search` mit
  `popularity=1.0`).
- **Shared Ingest-Helper** (`ingest/cache.go`, `ingest/retry.go`):
  `cachePath`, `loadCache`, `saveCache`, generischer `retryFetch[T]`
  mit pluggable `defaultRetryConfig`. Ersetzt die Phase-1-Duplikation
  in `github_trending.go` und `hackernews.go`.
- **Planning-Engine** (`internal/lodestone/planning/`): shell-out an
  `claude --print --model <id>` mit ablösbarem `Runner`-Interface
  (`ClaudeRunner` / `FakeRunner`). `BuildPrompt` produziert den
  deutschen Architekt-Prompt mit Fingerprint+Recommendation als JSON,
  `SplitResponse` parst `===SPEC===`/`===PLAN===`-Marker, `Persist`
  schreibt nach `docs/superpowers/specs/` und `…/plans/`.
- **Subkommando `lodestone plan <rec-id>`** mit `--dry-run` (zeigt
  nur den Prompt) und `--model`-Override.
- **Subkommando `lodestone init`**: legt `.lodestone.yaml` an, hängt
  `.gitignore`-Snippet (`.lodestone/` ignoriert, `decisions.log`
  ausgenommen) an, installiert vier Skills nach `.claude/skills/`.
- **Audit-Log** (`internal/lodestone/audit/`): jeder
  `fingerprint`/`ingest`/`score`/`plan`-Aufruf appended einen JSONL-
  Eintrag in `.lodestone/decisions.log`.
- **Vier Claude-Skills** unter `flavors/lodestone/skills/`:
  `lodestone-scout`, `-recommend`, `-plan`, `-review-trends` als
  Markdown-Frontmatter; im Binary via `go:embed` eingebettet.

### Added — Phase 1 MVP

Funktional vollständige, LLM-freie und deterministische Pipeline:
`fingerprint → ingest → score → signals`.

- Initiales Repo-Skelett: Go-Modul `github.com/jmt-labs/lodestone`,
  Cobra-CLI `cmd/lodestone`, CI-/Release-Workflows, GoReleaser-
  Konfiguration, Makefile, Linter-Konfiguration.
- Doku-Grundgerüst: `README.md`, `CONTRIBUTING.md`, `CLAUDE.md`,
  `AGENTS.md`, `base/models.yaml`, ausführlicher User-Guide unter
  `docs/lodestone.md`.
- Spec + Phase-1-Plan unter `docs/superpowers/`.
- **T1** — Schemas für `Signal`, `Fingerprint`, `Recommendation`,
  `WorkPackage` (`internal/lodestone/schema/`) inkl. JSON-Roundtrip-
  Tests.
- **T2** — Datei-basierter Store (`internal/lodestone/store/`) mit
  JSONL-Signals, JSON-Fingerprint, JSONL-Recommendations und atomarem
  Replace; In-Memory-Index für O(1)-`Has`.
- **T3** — `Source`-Interface und GitHub-Trending-Adapter
  (`/search/repositories`) mit konservativer Default-Query
  (`stars:>=50 pushed:>{90daysago}`), tagesgenauem Cache, optionalem
  `$GITHUB_TOKEN`, exponentiellem Backoff (max 3 Versuche).
- **T4** — HackerNews-Adapter (`/v0/topstories.json` + `/v0/item/<id>`)
  mit Story-Type- und Keyword-Filter (Default: `ai, llm, mcp, claude,
  agent`), Limit 50 pro Lauf, gleiche Cache-Semantik.
- **T5** — Fingerprint für Go und Node: Walker mit Skip-Verzeichnissen
  (`.git`, `vendor`, `node_modules`, `dist`, `build`), Regex-basiertes
  `go.mod`-Parsing (inline + Block-Require), `package.json`-Parsing
  inkl. devDependencies, Framework-Heuristik, Test-Ratio, CI- und
  MCP-Detection.
- **T6** — Deterministisches Scoring:
  - `compatibility`: gewichtete Jaccard, Language-Match 1.5×,
    Framework-Match 1.0×, normiert auf `[0,1]`
  - `effort`: Default `M`, `XL` bei 0 Match, `S` bei Match + Stars<100
  - `risk`: `low` bei Stars≥500 ∧ LastCommit<90d ∧ License; `high` bei
    fehlender License oder LastCommit>180d; sonst `med`
  - Sortierung `compatibility DESC, stars DESC, id ASC`; Determinismus
    über drei JSON-Byte-Vergleichs-Läufe verifiziert.
- **T7** — Config-Loader für `.lodestone.yaml` mit Defaults
  (`min_stars: 50`, `min_age_days: 30`, `max_last_commit_age_days:
  180`, `require_license: true`); Goals und TechInterests
  flowen in den Fingerprint.
- **T8** — Cobra-Subkommandos `fingerprint`, `ingest`, `score`,
  `signals` mit `--root`, `--source`, `--since`, `--top`, `--json`,
  `--mock` plus Source-Stubs für Phase 2-4.
- **T9** — End-to-End-Test `e2e/lodestone_test.sh` mit
  `LODESTONE_MOCK_FIXTURES` (offline-fähig), CI-Job `e2e`, Mock-
  Fixtures unter `e2e/fixtures/signals/`, Determinismus-Diff über
  zwei aufeinanderfolgende Score-Läufe.

### Dependencies

- `github.com/spf13/cobra` v1.10.2 (CLI)
- `gopkg.in/yaml.v3` v3.0.1 (Konfig-Parsing)
- Sonst nur Go-Standardbibliothek.
