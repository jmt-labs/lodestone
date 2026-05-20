# Plan-Format

Pläne sind das ausführbare Begleitstück zur Spec — sie zerteilen das
Vorhaben in atomare, jeweils mergebare Tasks. Pfad:

```
docs/superpowers/plans/YYYY-MM-DD-<thema>.md
```

## Aufbau

```markdown
# <Titel> — Plan

> Spec: docs/superpowers/specs/YYYY-MM-DD-<thema>-design.md
> Phase: X

## Ziele
Wiederholung der Akzeptanzkriterien aus der Spec, eins zu eins.

## Tasks

- [ ] **T1: <kurzer Task-Name>**
  - Branch: `feat/p<phase>-t1-<slug>`
  - Files: `internal/lodestone/...`, `cmd/lodestone/...`
  - Test: Unit-Test in `..._test.go`; E2E-Schritt in `e2e/`
  - DoD: <konkrete Done-Bedingung>

- [ ] **T2: …**
  - Branch: `feat/p<phase>-t2-<slug>`
  - …

## Reihenfolge & Abhängigkeiten
T1 → T2 → T3 (sequentiell), T4 unabhängig.

## Out-of-Scope
Was bewusst nicht in diesem Plan steckt (verweis auf späteren Plan oder
auf YAGNI).
```

## Regeln

- **Eine Task ≈ ein Branch ≈ ein PR.** Wenn eine Task in mehr als 2–3
  Commits zerfällt, ist sie zu groß — vor dem Start aufteilen.
- **Branch-Namen folgen** [Workflow § Branch-Namen-Schema](workflow.md#branch-namen-schema).
- **DoD ist überprüfbar.** „Funktioniert" ist kein DoD; „Unit-Test `X`
  grün, E2E-Schritt `Y` ergänzt" schon.
- **Tasks haken sich beim Merge ab.** Der Plan ist live — wer eine Task
  abschließt, updated die Checkbox im selben PR oder im Folge-PR.

## Beispiel: existierende Pläne

- [Phase-1-MVP-Plan](../superpowers/plans/2026-05-20-lodestone-mvp.md)
  — 12 Tasks, alle abgehakt.
- [Phase-2-Plan](../superpowers/plans/2026-05-20-lodestone-phase2.md)
  — kompakter, weil Phase 2 weniger getreuen Steps brauchte.

## Verwandt

- [Spec-Format](spec-format.md) — das was zum wie.
- [Workflow](workflow.md) — Spec → Plan → Branch → TDD → PR.
- [PR-Checkliste](pr-checklist.md) — was vor dem Merge verifiziert wird.
