# Contributing zu lodestone

Danke für dein Interesse! Lodestone folgt einem schlanken
Spec-Plan-PR-Workflow.

## Kurz-Einstieg

1. **Brainstorming** → Design-Idee.
2. **Spec** in `docs/superpowers/specs/YYYY-MM-DD-<thema>-design.md`.
3. **Plan** in `docs/superpowers/plans/YYYY-MM-DD-<thema>.md` mit
   Checkbox-Tasks.
4. **Branch** nach Schema `feat/p<phase>-t<task>-<slug>`.
5. **TDD-Implementierung** — bei Bug-Fix Regressionstest ZUERST.
6. **PR gegen `main`** — nur auf explizite Aufforderung.

## Sprache

- Doku, Specs, Pläne, Commit-Messages: **deutsch**.
- Code-Identifier und API-Felder: englisch.

## Lokale Entwicklung

```sh
make build         # baut bin/lodestone und bin/lodestone-mcp
make test          # go test ./... (mit -race)
make lint          # golangci-lint v2
make vuln          # govulncheck
make e2e           # End-to-End-Smoke-Test
```

Alle vier müssen vor jedem PR grün sein.

## Detailliert

Vollständige Doku im `docs/contributor/`-Bereich:

- [Workflow](docs/contributor/workflow.md) — Pipeline, Branch-Schema, PR-Body-Template.
- [Spec-Format](docs/contributor/spec-format.md) — Pflicht-Sektionen.
- [Plan-Format](docs/contributor/plan-format.md) — Checkbox-Tasks.
- [Skills-Policy](docs/contributor/skills-policy.md) — Pflicht-Skills, Regression-First-Regel.
- [Testing](docs/contributor/testing.md) — Test-Pyramide, Coverage-Ziel.
- [PR-Checkliste](docs/contributor/pr-checklist.md) — Pre-Merge-Gates.
- [Release-Prozess](docs/contributor/release-process.md) — GoReleaser-Workflow.
- [Doku-Wartung](docs/contributor/docs-maintenance.md) — wie die Doku synchron bleibt.

Vorgaben für KI-Tools: [`CLAUDE.md`](CLAUDE.md). Rollen-Übersicht:
[`AGENTS.md`](AGENTS.md).
