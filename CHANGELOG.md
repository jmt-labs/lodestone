# Changelog

Alle nennenswerten Änderungen an diesem Projekt werden in dieser Datei
dokumentiert. Format gemäß [Keep a Changelog](https://keepachangelog.com/de/1.1.0/),
Versionierung nach [Semantic Versioning](https://semver.org/lang/de/).

## [Unreleased]

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
