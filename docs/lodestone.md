# lodestone — User-Guide (Phase 1)

> Stand: 2026-05-20, Phase-1-MVP. LLM-frei, deterministisch.

## Installation

> **Noch nicht released.** Sobald `v0.1.0` getagged ist:

```sh
# Go install (immer verfügbar)
go install github.com/jmt-labs/lodestone/cmd/lodestone@latest

# Homebrew (geplant nach v0.1.0)
brew install jmt-labs/tap/lodestone

# apt (geplant nach v0.1.0) — über GoReleaser-erzeugtes .deb
```

Bis dahin lokal bauen:

```sh
git clone https://github.com/jmt-labs/lodestone.git
cd lodestone
make build      # → bin/lodestone
```

## Quickstart

```sh
cd dein-go-oder-node-projekt/
lodestone fingerprint
lodestone ingest
lodestone score
lodestone signals --top 10
```

Alle Artefakte landen unter `.lodestone/`. Voreingestellt ist nur
`decisions.log` (Audit-Trail) committable; alles andere gehört
typischerweise in `.gitignore`.

## Subkommandos

### `lodestone fingerprint`

Scannt das aktuelle Repo (oder `--root <pfad>`) und erzeugt
`.lodestone/fingerprint.json`. Erkennt:

- **Sprachen** anhand Dateierweiterungen (`.go`, `.js`/`.jsx`/`.ts`/`.tsx`)
- **Frameworks** über `go.mod` und `package.json` (cobra, anthropic-sdk,
  gin, echo, chi, react, vue, next, svelte, mcp-sdk, @anthropic-ai/sdk)
- **Dependencies** aus `go.mod` (Regex-Parse) und `package.json`
  (`dependencies` + `devDependencies`)
- **LOC pro Sprache** (Skip-Verzeichnisse: `.git`, `vendor`,
  `node_modules`, `dist`, `build`)
- **Test-Ratio** (Test-LOC / Non-Test-LOC)
- **CI-Provider** (`github_actions` / `gitlab_ci` / `circleci`)
- **MCP-Konfiguration** (Existenz von `.mcp.json`)
- **Goals / TechInterests** aus `.lodestone.yaml`

### `lodestone ingest`

Holt Signale von externen Quellen und schreibt sie deduppliziert nach
`.lodestone/signals.jsonl`.

```sh
lodestone ingest                                # alle Quellen
lodestone ingest --source github_trending       # nur eine
lodestone ingest --source github_trending --source hackernews
lodestone ingest --mock                         # offline aus $LODESTONE_MOCK_FIXTURES
```

**Quellen (Phase 1):**

- `github_trending` — `api.github.com/search/repositories`, gefiltert
  nach `stars:>=N pushed:>YYYY-MM-DD`. Liest `$GITHUB_TOKEN` (optional,
  empfohlen wegen Rate-Limit).
- `hackernews` — `hacker-news.firebaseio.com/v0/topstories.json`,
  Story-Type-Filter, Keyword-Match (Default: `ai, llm, mcp, claude,
  agent`).

**Cache:** Pro Quelle und Tag genau ein Cache-File unter
`.lodestone/cache/<source>-YYYY-MM-DD.json`. Zweiter Aufruf am selben
Tag liest aus dem Cache.

### `lodestone score`

Lädt `fingerprint.json` und alle Signale, berechnet pro Signal die
Dimensionen und schreibt nach `.lodestone/recommendations.jsonl`
(deterministisch sortiert: `compatibility DESC, stars DESC, id ASC`).

| Dimension | Werte | Berechnung |
|---|---|---|
| `compatibility` | 0.0 – 1.0 | Gewichtete Jaccard auf Signal-Tags ∩ Repo-Frameworks/-Languages. Language-Match 1.5×, Framework-Match 1.0×. |
| `effort` | XS – XL | Default `M`; `XL` wenn 0 Match; `S` bei Match und Stars < 100. |
| `risk` | low / med / high | `low` bei Stars≥500 ∧ LastCommit<90d ∧ License; `high` bei fehlender License oder LastCommit>180d (stale); sonst `med`. |

### `lodestone signals`

```sh
lodestone signals                          # alle Signale, sortiert nach Stars
lodestone signals --top 20                 # Top-20
lodestone signals --since 2026-05-01       # nur ab Datum
lodestone signals --source hackernews      # nur eine Quelle
lodestone signals --json                   # JSON statt Tabelle
```

## Konfiguration (`.lodestone.yaml`)

Optional im Repo-Wurzelverzeichnis. Komplettes Beispiel:

```yaml
goals:
  - reliability
  - speed
  - shipping

tech_interests:
  - mcp
  - llm-tools
  - agent-frameworks

lodestone:
  min_stars: 50                # github-trending Stars-Filter
  min_age_days: 30             # github-trending: nur Repos pushed innerhalb von N Tagen
  max_last_commit_age_days: 180
  require_license: true
```

Alle Felder sind optional; Defaults werden zugemischt. `goals` und
`tech_interests` landen im `fingerprint.json` und beeinflussen
zukünftiges LLM-basiertes Planning ab Phase 2.

## Lokale Artefakte

```
.lodestone/
├── cache/
│   ├── github_trending-2026-05-20.json   # tagesgenauer Quell-Cache
│   └── hackernews-2026-05-20.json
├── signals.jsonl                          # JSON-Lines, append-only, dedupliziert
├── fingerprint.json                       # einzelne JSON-Datei, atomic write
├── recommendations.jsonl                  # JSON-Lines, atomic replace
└── decisions.log                          # Audit-Trail (ab Phase 2)
```

**Empfehlung für `.gitignore`** (das Skill-Install ab Phase 2 fügt das
automatisch ein):

```
.lodestone/
!.lodestone/decisions.log
```

## Determinismus

Zwei aufeinanderfolgende `lodestone score`-Läufe mit identischem
Fingerprint und identischer Signal-Liste produzieren byte-identische
`recommendations.jsonl`. Verifiziert durch:

- Unit-Test in `internal/lodestone/scoring/scoring_test.go`
  (`TestScoreDeterminism`, drei Läufe mit JSON-Byte-Vergleich).
- End-to-End-Test in `e2e/lodestone_test.sh` (Determinismus-Schritt
  diffed `recommendations.jsonl` zwischen zwei Score-Läufen).

## Troubleshooting

| Symptom | Ursache / Lösung |
|---|---|
| `Error: no signals in store (run \`lodestone ingest\` first)` | `score` braucht erst Signale. `lodestone ingest` ausführen oder Mock-Modus mit `--mock` + `$LODESTONE_MOCK_FIXTURES`. |
| `Error: unknown source "X"` | Phase-1-Quellen heißen `github_trending` und `hackernews`. |
| `--mock requires $LODESTONE_MOCK_FIXTURES` | Env-Var muss auf ein Verzeichnis mit `<source>.json`-Fixtures zeigen. |
| `Error: read fingerprint (run \`lodestone fingerprint\` first)` | `score` braucht den Fingerprint. Erst `lodestone fingerprint`. |
| `Error: github_trending: max retries exceeded` | Rate-Limit oder Netzwerk. `$GITHUB_TOKEN` setzen reduziert die Rate-Limit-Wahrscheinlichkeit erheblich. |

## Was lodestone nicht tut (Phase 1)

- **Keine LLM-Aufrufe.** Compatibility/Effort/Risk sind rein
  deterministisch. LLM-Planning kommt in Phase 2.
- **Keine Auto-PRs.** Phase 4 wird das mit harten Schranken einführen.
- **Keine Telemetrie.** Alle Daten bleiben lokal.
- **Kein Daemon.** Manuelle Triggerung; Scheduling via Cron oder
  GitHub-Action-Template (Phase 3) ist opt-in.

Vollständige Roadmap: siehe
[`docs/superpowers/specs/2026-05-20-lodestone-design.md`](superpowers/specs/2026-05-20-lodestone-design.md).
