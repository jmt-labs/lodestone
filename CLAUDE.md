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

## Phase-1-Invarianten

Diese Regeln gelten bis zum Phase-1-Abschluss:

1. **Keine LLM-Aufrufe** in produktivem Code-Pfad. Erst ab Phase 2.
2. **Deterministische Pipeline** — zwei Score-Läufe mit identischem
   Input müssen byte-identische sortierte Outputs liefern.
3. **Minimaler Dependency-Footprint** — nur `github.com/spf13/cobra`
   und `gopkg.in/yaml.v3`. Keine neuen externen Go-Deps ohne
   explizite Diskussion.
4. **Anti-Hype-Defaults konservativ** — `min_stars: 50`,
   `min_age_days: 30`, `max_last_commit_age_days: 180`,
   `require_license: true`.

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
