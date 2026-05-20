# Lodestone Phase 4 — Design

**Datum:** 2026-05-20
**Voraussetzung:** Phasen 1-3 auf `main`.
**Status:** Konservative Auslegung — Auto-PR ist **immer Draft**, **nie
auf `main`**, **nie ohne harte Sicherheitsgates**.

## Ziele

Phase 4 macht aus Recommendations konkrete Pull-Requests gegen das Ziel-
Repo. Das Risiko-Profil ist deutlich höher als Phase 1-3 (Git-Schreib-
operationen, GitHub-API-Calls), daher kommen mehrere unabhängige
Schranken zum Tragen.

## Auto-PR-Engine

`lodestone apply <rec-id>`:

1. **Safety-Gates** (alle müssen erfüllt sein, sonst Abbruch):
   - `risk == low`
   - `effort == XS`
   - `compatibility >= 0.85`
   - kein Auto-PR in den letzten 24h gegen denselben Branch (Rate-Limit
     aus `.lodestone/applies.jsonl`)
   - `git status` sauber (kein staged/unstaged-Diff)
2. **Plan-Generierung** (delegiert an Planning-Engine aus Phase 2):
   schreibt Spec + Plan unter `docs/superpowers/{specs,plans}/`.
3. **Branch** anlegen nach `lodestone/apply-<rec-id-suffix>-<date>`.
   **Nie auf `main`** committen — `git switch -c` erzwungen.
4. **Commit + Push** der zwei Markdown-Dateien.
5. **PR-Erstellung als Draft** über die GitHub-CLI (`gh pr create
   --draft`). Falls `gh` nicht im PATH oder kein Token: Branch ist
   gepusht, PR muss manuell geöffnet werden.
6. **Audit + Apply-Log**: `applies.jsonl`-Eintrag mit Rec-ID, Branch,
   PR-Nummer, Status.

**Mehr als 1 Auto-PR pro Tag pro Repo verhindert** durch Zeitstempel-
Prüfung in `applies.jsonl`.

## Undo

`lodestone undo <branch-or-pr>`:

1. Findet den `applies.jsonl`-Eintrag.
2. Schließt PR (falls vorhanden, via `gh pr close --delete-branch`).
3. Löscht Branch lokal und remote.
4. Schreibt `undo`-Audit-Entry.

Kein Revert auf `main`, weil Auto-PRs nie nach `main` mergen ohne
expliziten User-Approval.

## Success-Tracker

`lodestone stats`: aggregiert aus `applies.jsonl` Erfolgs-Statistiken
(merged / closed / open) je Rec-Quelle und gibt eine Tabelle aus. Reine
Lese-Operation, keine LLM-Aufrufe.

## Cross-Repo-Sharing — Privacy-Spec separat

Cross-Repo-Sharing (z. B. anonymisiertes Teilen von erfolgreichen
Recommendations zwischen Repos) erfordert eine eigene Privacy-Spec, die
explizit klärt:

- Welche Felder eines `Recommendation`/`Signal` dürfen geteilt werden?
- Wie wird der Datenfluss konsentiert und revoziert?
- Wo werden geteilte Daten persistiert (Default: gar nicht)?
- Welche Re-Identifikations-Risiken bleiben nach Anonymisierung?

Diese Spec liegt unter
`docs/superpowers/specs/2026-05-20-lodestone-sharing-privacy.md` und wird
in Phase 4 nur **dokumentiert**, nicht implementiert.

## Code-Layout

```
internal/lodestone/apply/
  apply.go            apply_test.go
  gates.go            gates_test.go
  state.go            state_test.go

cmd/lodestone/
  apply.go
  undo.go
  stats.go
```

## Test-Strategie

- `gates_test.go`: deckt jede Safety-Gate-Verletzung einzeln ab
- `apply_test.go`: nutzt Fake-Git und Fake-`gh`-Runner (Shell-out, mock-
  fähig)
- `state.go` schreibt/liest `applies.jsonl` mit `t.TempDir()`-Tests
- E2E: nicht erweitert (Git-Operationen brauchen echtes Repo + Remote)

## Phase-4-Tasks

- [ ] **P4-T1** — Privacy-Spec für Cross-Repo-Sharing.
- [ ] **P4-T2** — Safety-Gates (`apply/gates.go`).
- [ ] **P4-T3** — Apply-State (`apply/state.go`): `applies.jsonl`.
- [ ] **P4-T4** — Apply-Orchestration (`apply/apply.go`) mit
      injizierbaren Git- und `gh`-Runnern.
- [ ] **P4-T5** — Subkommandos `apply`, `undo`, `stats`.
- [ ] **P4-T6** — Docs Update + Phase-4-Eintrag im CHANGELOG.

## Ausdrücklich nicht in Phase 4

- LLM-getriebene Code-Änderungen (über Spec/Plan-Dokumente hinaus).
- Push direkt auf `main`.
- PR-Auto-Merge.
- Verteiltes Sharing.
