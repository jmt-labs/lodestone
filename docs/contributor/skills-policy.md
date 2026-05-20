# Skills-Policy

Pflicht-Skills für KI-getriebene Arbeit am Lodestone-Repo. Quelle:
[`CLAUDE.md`](../../CLAUDE.md).

## Pflicht-Skills

| Situation | Skill |
|---|---|
| Vor Feature / Implementierung | `superpowers:brainstorming` |
| Während Implementierung | `superpowers:test-driven-development` |
| Vor Commit / PR | `superpowers:verification-before-completion` |
| Bei Bugs | `superpowers:systematic-debugging` (Regressionstest ZUERST) |

## Regression-First-Regel

Bei jedem Bug:

1. **Regressionstest schreiben**, der den Bug reproduziert.
   Lokal verifizieren, dass er ohne Fix rot ist.
2. **Den Test alleine committen.** Der Commit zeigt den Bug an.
3. **Den Fix als separaten Commit** anhängen. Der Test wird grün.

Das macht die Bug-Geschichte im Git-Log nachvollziehbar und schließt
strukturell aus, dass jemand „den Fix einbaut, ohne ihn zu beweisen".

## Rollen und Skill-Zuordnung

Aus [`AGENTS.md`](../../AGENTS.md):

| Rolle | Modell | Pflicht-Skills |
|---|---|---|
| Planner | `claude-opus-4-7` | `superpowers:writing-plans`, `superpowers:brainstorming` |
| Implementer | `claude-sonnet-4-6` | `superpowers:test-driven-development`, `superpowers:verification-before-completion` |
| Reviewer | `claude-sonnet-4-6` | `superpowers:verification-before-completion` |
| Mechanic | `claude-haiku-4-5-20251001` | — |

## YAGNI als Begleit-Regel

- Keine spekulativen Features.
- Keine Abstraktionen ohne konkreten zweiten Aufrufer.
- Keine Error-Handling-Pfade für Szenarien, die nicht eintreten können.
- Keine Backward-Compat-Shims für Code, der noch nicht released ist.
- Keine Kommentare im Code, außer das **Warum** ist non-obvious
  (versteckte Constraint, subtile Invariante, Workaround für Bug).

Quelle: [`CLAUDE.md`](../../CLAUDE.md) § YAGNI & Code-Qualität.

## Verwandt

- [Workflow](workflow.md) — wo die Skills in den Spec-Plan-Branch-PR-Loop
  eingebaut sind.
- [Testing](testing.md) — was die Skills durchsetzen.
- [PR-Checkliste](pr-checklist.md) — die finale Gate vor Merge.
