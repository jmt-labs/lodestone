# lodestone — Dokumentation

> Topic-Index für alle Doku-Artefakte. Drei Zielgruppen, drei Bereiche.

`lodestone` sammelt AI-Ökosystem-Signale, scort sie deterministisch
gegen einen Repo-Fingerprint und liefert reproduzierbare Empfehlungen
— als CLI, als MCP-Server und als Claude-Skill-Pack. Phasen 1–4 sind
auf `main` gemerged.

Aktueller Phasen-Status: siehe [`internals/roadmap.md`](internals/roadmap.md).

## Für User

Du willst lodestone in deinem Projekt einsetzen.

| Dokument | Inhalt |
|---|---|
| [Installation](user/installation.md) | `go install`, Build aus Source, geplante Pakete |
| [Quickstart](user/quickstart.md) | In 60 Sekunden vom Klon zur ersten Empfehlung |
| [Befehle](user/commands/README.md) | Alle 11 Verben + `lodestone-mcp`, ein Detail-Doc pro Verb |
| [Konfiguration](user/configuration.md) | `.lodestone.yaml` und Environment-Variablen |
| [MCP-Server](user/mcp-server.md) | `lodestone-mcp` mit Claude Desktop / Claude Code einrichten |
| [Skills](user/skills.md) | Vier Claude-Skills, installierbar via `lodestone init` |
| [FAQ](user/faq.md) | Welcher Befehl für welchen Use-Case? Decision-Tree |
| [Troubleshooting](user/troubleshooting.md) | Häufige Fehler und ihre Ursachen |
| [Glossar](user/glossary.md) | Fingerprint, Signal, Recommendation, Compatibility, Effort, Risk, Phasen |

## Für Contributors

Du willst etwas an lodestone ändern oder erweitern.

| Dokument | Inhalt |
|---|---|
| [Workflow](contributor/workflow.md) | Spec → Plan → Branch → TDD → PR |
| [Spec-Format](contributor/spec-format.md) | Konvention für `docs/superpowers/specs/` |
| [Plan-Format](contributor/plan-format.md) | Checkbox-Tasks-Konvention |
| [Skills-Policy](contributor/skills-policy.md) | Pflicht-Skills, Regression-First-Regel |
| [Testing](contributor/testing.md) | Test-Pyramide, Coverage-Ziel, Determinismus-Tests |
| [PR-Checkliste](contributor/pr-checklist.md) | Pre-PR-Gates, Definition of Done |
| [Release-Prozess](contributor/release-process.md) | Tag-getriebener GoReleaser-Workflow |
| [Doku-Wartung](contributor/docs-maintenance.md) | Wie Doku synchron zum Code bleibt |

## Für Architekten / KI

Du willst verstehen, wie lodestone intern funktioniert.

| Dokument | Inhalt |
|---|---|
| [Architektur](internals/architecture.md) | Drei-Ebenen-Modell, Code-Layout, Datenfluss, Erweiterungspunkte |
| [Datenmodell](internals/data-model.md) | Signal, Fingerprint, Recommendation, WorkPackage |
| [Scoring-Algorithmus](internals/scoring.md) | Compatibility / Effort / Risk: Formeln und Schwellen |
| [Determinismus](internals/determinism.md) | Byte-Identität-Garantie und Verifikation |
| [Artefakte](internals/artifacts.md) | Layout unter `.lodestone/` |
| [Roadmap](internals/roadmap.md) | Phasen 1–5+ mit Status |
| [ADRs](internals/adr/README.md) | Architecture Decision Records, rückwirkend extrahiert |

## Superpowers-Specs und -Pläne

Historische Designs und Implementierungspläne pro Phase. Index mit
logischen Aliasen: [`superpowers/README.md`](superpowers/README.md).
