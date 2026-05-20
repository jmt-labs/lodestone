# Workflow

Lodestone-Entwicklung folgt einem schlanken Spec-Plan-PR-Workflow. Kein
Sprint-Boards, keine Tickets — der Plan ist Single Source of Truth pro
Phase, und Branches korrespondieren 1:1 mit Plan-Tasks.

## Pipeline

```
Brainstorming
    ↓
Spec (docs/superpowers/specs/YYYY-MM-DD-<thema>-design.md)
    ↓
Plan (docs/superpowers/plans/YYYY-MM-DD-<thema>.md, Checkbox-Tasks)
    ↓
Branch (feat/p<phase>-t<task>-<slug>)
    ↓
TDD-Implementierung (eine Task pro Branch, atomare Commits)
    ↓
PR gegen main (nur auf explizite Aufforderung)
    ↓
Review (Reviewer-Rolle, siehe AGENTS.md)
    ↓
Merge → main
```

## Branch-Namen-Schema

| Pattern | Zweck |
|---|---|
| `feat/p<phase>-t<task>-<slug>` | Plan-Task umsetzen (z. B. `feat/p1-t3-github-trending`) |
| `fix/<slug>` | Bug-Fix mit Regressionstest |
| `chore/<slug>` | Wartung (CI, Doku, Refactoring ohne Verhaltensänderung) |
| `docs/<slug>` | Reine Doku-Änderung |

## Direkt-`main`-Verbot

> Ausnahme: Der erste Bootstrap-Commit dieses Repos wurde explizit
> autorisiert. Danach gilt: **alles über Branch + PR**, niemals direkt
> auf `main`. Kein force-push auf `main`. Siehe [`CLAUDE.md`](../../CLAUDE.md).

## Spec → Plan

- **Spec** = Was und Warum. Architektur-Skizze, Akzeptanzkriterien,
  Alternativen. Format: [Spec-Format](spec-format.md).
- **Plan** = Wie und in welcher Reihenfolge. Checkbox-Tasks
  (`- [ ] T1: …`). Format: [Plan-Format](plan-format.md).

## TDD-Regel

| Situation | Reihenfolge |
|---|---|
| Neues Feature | Test → fehlschlagen → Implementierung → grün. |
| Bug-Fix | **Regressionstest ZUERST** committen (eigener Commit), dann den Fix als separater Commit. |
| Refactoring | Tests bleiben unverändert, müssen vor und nach dem Refactor grün sein. |

Hintergrund:
[Skills-Policy](skills-policy.md) (`superpowers:test-driven-development`,
`superpowers:systematic-debugging`).

## PR-Erstellung

PRs werden **nur auf explizite Aufforderung** geöffnet. Push auf den
Feature-Branch ist immer erlaubt; das `gh pr create` braucht eine
explizite User-Anweisung.

## PR-Body-Template

```markdown
## Zusammenfassung
- Was hat sich geändert?
- Warum?

## Spec / Plan
- Spec: docs/superpowers/specs/…
- Plan-Task: T<n>

## Test-Plan
- [ ] make test grün
- [ ] make lint grün
- [ ] make vuln grün
- [ ] make e2e grün
- [ ] Doku updated (falls User-facing)

## Referenzen
Updates #<epic-issue-nummer>
```

Vor dem Merge: alle Punkte in [PR-Checkliste](pr-checklist.md) abhaken.

## Sprache

- **Doku, Specs, Pläne, Commit-Messages: deutsch.**
- Code-Identifier und API-Felder: englisch.
- Keine unnötigen Kommentare im Code — siehe
  [`CLAUDE.md`](../../CLAUDE.md) § YAGNI.
