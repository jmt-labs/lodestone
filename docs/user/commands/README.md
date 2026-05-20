# Befehle

Vollständige Subkommando-Übersicht. Detail-Doku pro Verb in den
verlinkten Dateien. Für die erzählende Einführung siehe
[Quickstart](../quickstart.md).

## CLI `lodestone`

### Phase 1 — Deterministische Pipeline

| Verb | Status | Zweck |
|---|---|---|
| [`fingerprint`](fingerprint.md) | ✅ | Aktuelles Repo analysieren → `.lodestone/fingerprint.json` |
| [`ingest`](ingest.md) | ✅ | Externe Signale abrufen → `.lodestone/signals.jsonl` |
| [`score`](score.md) | ✅ | Signale × Fingerprint scoren → `.lodestone/recommendations.jsonl` |
| [`signals`](signals.md) | ✅ | Gespeicherte Signale filtern und anzeigen |

### Phase 2 — Planning + Onboarding

| Verb | Status | Zweck |
|---|---|---|
| [`init`](init.md) | ✅ | `.lodestone.yaml`, `.gitignore`-Snippet, Skills nach `.claude/skills/` |
| [`plan`](plan.md) | ✅ | Spec / Plan aus Recommendation generieren (ruft `claude`-CLI) |

### Phase 3 — Remote-Schnittstellen

| Verb | Status | Zweck |
|---|---|---|
| [`memory`](memory.md) | ✅ | `decisions.log` → `.claude/memory.json` konsolidieren |

### Phase 4 — Auto-PR-Engine

| Verb | Status | Zweck |
|---|---|---|
| [`apply`](apply.md) | ✅ | Recommendation als Draft-PR (vier Safety-Gates) |
| [`undo`](undo.md) | ✅ | Apply rückgängig: PR schließen + Branch löschen |
| [`stats`](stats.md) | ✅ | Apply-Erfolgs-Statistik aus `.lodestone/applies.jsonl` |

### Sonstige

| Verb | Status | Zweck |
|---|---|---|
| [`version`](version.md) | ✅ | Versionsinformationen anzeigen |

### Phase 5+ (geplant)

| Verb | Status | Zweck |
|---|---|---|
| `recommend` | 🚧 Stub | Interaktive Empfehlungs-Loop als Claude-Skill |
| `calibrate` | 🚧 Stub | Scoring-Gewichte aus Decision-Log nachjustieren |
| `share` | 🚧 Stub | Cross-Repo-Sharing (siehe [Privacy-Spec](../../superpowers/specs/2026-05-20-lodestone-sharing-privacy.md)) |

## Zweites Binary `lodestone-mcp`

MCP-Server über stdio JSON-RPC 2.0 mit fünf Tools — Setup und Tool-Liste
siehe [MCP-Server](../mcp-server.md).

## Globale Flags

| Flag | Default | Zweck |
|---|---|---|
| `--root` | `$PWD` | Projekt-Wurzel, in der `.lodestone/` lebt |

## Verwandt

- [Konfiguration](../configuration.md) — `.lodestone.yaml` und ENV.
- [Glossar](../glossary.md) — Compatibility, Effort, Risk, …
- [Scoring-Algorithmus](../../internals/scoring.md) — Formeln.
