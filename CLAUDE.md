# CLAUDE.md — Vorgaben für Claude in diesem Repo

> Diese Datei wird von Claude Code automatisch gelesen. Sie definiert
> die nicht-verhandelbaren Konventionen für jegliche KI-getriebene
> Arbeit an `lodestone`.

## Sprache

**Antworte auf Deutsch.** Commit-Messages, PR-Bodies, Doku, Specs und
Pläne auf Deutsch. Code-Identifier, API-Felder, Log-Messages und
Datei-Inhalte technischer Natur (z. B. `go.mod`) bleiben englisch.

## Pflicht-Skills

| Situation | Skill | Anforderung |
|---|---|---|
| Feature / Bug-Fix | `superpowers:brainstorming` | **MUSS vor Code aufgerufen werden** |
| Implementierung | `superpowers:test-driven-development` | **MUSS vor Code aufgerufen werden** |
| Vor Commit / PR | `superpowers:verification-before-completion` | **MUSS ausgeführt werden** |
| Debugging | `superpowers:systematic-debugging` | **MUSS vor Fix aufgerufen werden** |

**Bug-Regel:** Bei jedem Bug schreibst du den **Regressionstest ZUERST**,
committest ihn separat, und committest dann den Fix.

## Branch- & PR-Workflow

1. Brainstorming → Design abgestimmt.
2. **Spec** in `docs/superpowers/specs/YYYY-MM-DD-<thema>-design.md`.
3. **Plan** in `docs/superpowers/plans/YYYY-MM-DD-<thema>.md` mit
   Checkbox-Tasks (`- [ ] Tx: …`).
4. **Branch** anlegen — Schema `feat/p<phase>-t<task>-<slug>` für
   Plan-Tasks, `fix/<slug>` für Bugs, `chore/<slug>` für Wartung.
5. **TDD-Implementierung** mit kurzen, atomaren Commits.
6. **PR gegen `main`** mit Body, der die Spec/Plan-Datei verlinkt und
   das Epic-Issue mit `Updates #N` referenziert.

### Niemals direkt auf `main`

**Ausnahme:** Der erste Bootstrap-Commit dieses Repos wurde explizit auf
`main` autorisiert. Danach gilt: **alles über Branch + PR**, nie direkt
auf `main`. Kein force-push auf `main`.

### PR-Erstellung

PRs werden **nur auf explizite Aufforderung des Users erstellt**.
Push auf Feature-Branch ist erlaubt, PR-Eröffnung erfordert
expliziten Auftrag.

## YAGNI & Code-Qualität

- Keine spekulativen Features.
- Keine Abstraktionen ohne konkreten zweiten Aufrufer.
- Keine Error-Handling-Pfade für Szenarien, die nicht eintreten können.
- Keine Backward-Compat-Shims für Code, der noch nicht released ist.
- **Keine Kommentare im Code**, außer das WARUM ist non-obvious
  (versteckte Constraint, subtile Invariante, Workaround für Bug).

## Historische Invarianten (Phase 1) und aktuelle Geltung

Diese vier Invarianten wurden für Phase 1 aufgestellt. Punkte 1 und 3
sind ab Phase 2 explizit gelockert; Punkte 2 und 4 gelten unverändert
weiter und sind als ADRs festgehalten.

1. **~~Keine LLM-Aufrufe~~** — ab Phase 2 erlaubt in `lodestone plan`
   und `lodestone apply` (über die `claude`-CLI). **Niemals** im
   Score-Pfad — siehe
   [ADR-0006](docs/internals/adr/0006-deterministisches-scoring.md).
2. **Deterministische Pipeline** — zwei Score-Läufe mit identischem
   Input müssen byte-identische sortierte Outputs liefern. Gilt
   weiter, siehe
   [ADR-0006](docs/internals/adr/0006-deterministisches-scoring.md)
   und [Determinismus](docs/internals/determinism.md).
3. **~~Minimaler Dependency-Footprint~~** — ab Phase 2 gelockert auf
   „Standardbibliothek bevorzugen; neue Deps brauchen
   Spec-Diskussion". Aktueller Stand: weiter nur `cobra` und
   `yaml.v3`.
4. **Anti-Hype-Defaults konservativ** — `min_stars: 50`,
   `min_age_days: 30`, `max_last_commit_age_days: 180`,
   `require_license: true`. Gilt weiter, siehe
   [ADR-0005](docs/internals/adr/0005-anti-hype-defaults.md).

Detaillierte Contributor-Vorgaben:
[`docs/contributor/workflow.md`](docs/contributor/workflow.md),
[`docs/contributor/skills-policy.md`](docs/contributor/skills-policy.md).
Phasen-Status: [`docs/internals/roadmap.md`](docs/internals/roadmap.md).

## Modell-Routing

Siehe `base/models.yaml`. Kurzfassung:

- `planning` → `claude-opus-4-7` (Specs, Pläne, Architektur)
- `default` → `claude-sonnet-4-6` (Implementierung, Reviews)
- `mechanical` → `claude-haiku-4-5-20251001` (Roundtrips,
  Format-Konvertierung, Rationale-Generierung)
- `review` → `claude-sonnet-4-6` (PR-Review, Spec-Critique)

## Testing

- **Coverage-Ziel** für `internal/`-Pakete: ≥ 70 %.
- Tests beschreiben Verhalten, nicht Implementierung.
- E2E in `e2e/`, ausführbar via `make e2e`.
- Vor PR: `make test lint vuln` muss grün sein.

## Was lodestone NICHT ist

- Kein Daemon. Keine Hintergrund-Prozesse.
- Kein Telemetrie-Sender. Lokale Artefakte bleiben lokal.
- Kein Auto-Editor für `main`-Branches. Auto-PRs (Phase 4) immer auf
  Feature-Branch, immer als Draft, immer mit harten Schranken.
