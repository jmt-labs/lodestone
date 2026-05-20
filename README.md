# lodestone

**Liest das AI-Ökosystem für dein Repo.**

`lodestone` ist ein Trend-Intelligence-System für Code-Repositories: Es
sammelt Signale aus der AI-Welt (GitHub-Trending, HackerNews, MCP-Registry,
Framework-Releases, ArXiv, npm/PyPI), bildet einen Repo-Fingerprint und
scort jedes Signal gegen dein Projekt entlang **Compatibility**, **Effort**,
**ROI** und **Risk**. Daraus entstehen reproduzierbare Empfehlungen — und
optional Spec/Plan/Tasks-Triples im superpowers-Format.

Der Name ist Programm: Ein Lodestone ist ein natürlich magnetischer
Stein — ein Kompass, der zeigt, ohne den Kurs vorzuschreiben.

## Status

**Pre-Alpha / Phase 1 (MVP)** — siehe
[`docs/superpowers/plans/2026-05-20-lodestone-mvp.md`](docs/superpowers/plans/2026-05-20-lodestone-mvp.md).

Phase 1 ist bewusst **LLM-frei**: deterministisches Scoring, lokale
JSONL-Dateien, keine externen Calls außer den konfigurierten Quellen.

## Installation

> **Noch nicht released.** Sobald das erste Tag steht:

```sh
# Homebrew (geplant)
brew install jmt-labs/tap/lodestone

# apt (geplant) — via GoReleaser-erzeugtes .deb

# Go install
go install github.com/jmt-labs/lodestone/cmd/lodestone@latest
```

## Subkommandos (Phase 1)

| Verb | Status | Zweck |
|---|---|---|
| `lodestone ingest` | T3 / T4 | Externe Signale abrufen → `.lodestone/signals.jsonl` |
| `lodestone fingerprint` | T5 | Aktuelles Repo analysieren → `.lodestone/fingerprint.json` |
| `lodestone score` | T6 | Signale × Fingerprint scoren → `.lodestone/recommendations.jsonl` |
| `lodestone signals` | T8 | Gespeicherte Signale filtern / anzeigen |
| `lodestone plan` | Phase 2 | Spec / Plan / Tasks aus Recommendation generieren |
| `lodestone recommend` | Phase 2 | Empfehlungen interaktiv durchgehen |
| `lodestone apply` / `undo` | Phase 4 | Auto-PR-Engine |
| `lodestone stats` / `calibrate` / `share` | Phase 3+ | Telemetrie- / Tuning-Werkzeuge |

## Lokale Artefakte

`lodestone` schreibt alle Zustände in `.lodestone/` im Zielprojekt:

```
.lodestone/
├── cache/                    # Roh-Fetches mit TTL-Datum
├── signals.jsonl             # Normalisierte, deduplizierte Signale
├── fingerprint.json          # Letzter Repo-Scan
├── recommendations.jsonl     # Gescorte Vorschläge
└── decisions.log             # Audit-Trail (einzige Datei nicht in .gitignore)
```

## Architektur

Drei Ebenen:

1. **CLI `lodestone`** (Cobra, Go) — deterministische Pipeline.
2. **Claude-Skills** (`flavors/`-Stil, Phase 2) — direkt installierbar
   in beliebige Repos.
3. **MCP-Server `lodestone-mcp`** (Phase 3) — separates Go-Binary mit
   `list_signals`, `query_trends`, `score_repo`, `generate_plan`,
   `record_decision`.

Details: [`docs/superpowers/specs/2026-05-20-lodestone-design.md`](docs/superpowers/specs/2026-05-20-lodestone-design.md).

## Entwicklung

Siehe [`CONTRIBUTING.md`](CONTRIBUTING.md). Pflicht-Workflow: Spec →
Plan → Branch → TDD → PR.

## Lizenz

MIT — siehe [`LICENSE`](LICENSE).
