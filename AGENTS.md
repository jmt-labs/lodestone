# AGENTS.md — Team-Rollen & Modell-Zuordnung

`lodestone` arbeitet mit einem schlanken Multi-Modell-Setup. Jede Rolle
hat einen klaren Zuständigkeitsbereich und ein zugewiesenes Modell aus
`base/models.yaml`.

## Rollen

### Planner — `planning`-Modell

- Designt Specs in `docs/superpowers/specs/`.
- Schreibt Pläne mit Checkbox-Tasks in `docs/superpowers/plans/`.
- Tradeoff-Analyse, Architektur-Entscheidungen, Phasen-Planung.
- **Niemals** direkt Code committen — übergibt an Implementer.

### Implementer — `default`-Modell

- Setzt Plan-Tasks um, eine pro Branch.
- Folgt TDD (Test → fehlschlagen → Implementierung → grün).
- Bei Bugs: Regressionstest **zuerst**, dann Fix.
- Hält sich an YAGNI; bricht nicht in Architektur-Diskussionen aus.

### Reviewer — `review`-Modell

- Reviewt PRs vor Merge.
- Prüft: Spec-Treue, Test-Coverage, Determinismus, YAGNI-Verstöße,
  Dependency-Hygiene.
- Eskaliert architektonische Bedenken zurück an den Planner.

### Mechanic — `mechanical`-Modell

- Generiert Rationale-Strings und Counter-Evidence für Recommendations
  (ab Phase 2 mit echten LLM-Aufrufen).
- Format-Konvertierungen, Schema-Roundtrips, einfache Refactorings.
- Schnelles, günstiges Modell für hohes Volumen.

## Workflow zwischen Rollen

```
Brainstorming (Planner)
    ↓
Spec + Plan (Planner)
    ↓
Branch + TDD-Impl (Implementer)
    ↓
PR (Implementer)
    ↓
Review (Reviewer) — bei Bedarf Re-Loop zum Implementer oder Planner
    ↓
Merge → main
```

## Wann welches Modell?

- **Komplexe Design-Entscheidung?** Planner (`opus`).
- **Mechanische Umsetzung eines klaren Plans?** Implementer (`sonnet`).
- **Eine einzelne JSON-Roundtrip-Test-Generierung, Rationale-Satz, oder
  Format-Konvertierung?** Mechanic (`haiku`).
- **PR-Review mit Architektur-Implikationen?** Reviewer (`sonnet`).

## Pflicht-Skills pro Rolle

| Rolle | Skill |
|---|---|
| Planner | `superpowers:writing-plans`, `superpowers:brainstorming` |
| Implementer | `superpowers:test-driven-development`, `superpowers:verification-before-completion` |
| Reviewer | `superpowers:verification-before-completion` |
| Mechanic | — (atomare Aufgaben) |

Bei Bug-Fixes (jede Rolle): `superpowers:systematic-debugging` mit
**Regressionstest VOR Fix**.
