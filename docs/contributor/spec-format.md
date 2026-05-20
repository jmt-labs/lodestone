# Spec-Format

Specs sind das forward-looking Artefakt — sie definieren, was als
Nächstes gebaut wird. Pfad und Namensschema:

```
docs/superpowers/specs/YYYY-MM-DD-<thema>-design.md
```

Beispiel: `docs/superpowers/specs/2026-05-20-lodestone-phase4-design.md`.

## Pflicht-Sektionen

```markdown
# <Titel> — Design

> Stand: YYYY-MM-DD. Phase X.

## Kontext
Warum überhaupt? Was ist das Problem? Was sind die Constraints
(Phase-Invarianten, Dependency-Budget, …)?

## Ziel
Eine bis drei Bullet-Points, was die Spec liefert. Keine Zeitachsen.

## Lösung
Architektur-Skizze, Datenflüsse, Datenmodelle. Diagramme als
ASCII-Art bevorzugt (Markdown-rendert sie direkt).

## Akzeptanzkriterien
- [ ] AC1: <überprüfbares Kriterium>
- [ ] AC2: …

## Alternativen
Was wurde verworfen und warum? (Wenn die Alternative später relevant
wird, ein eigener ADR.)

## Offene Fragen
Was ist noch unklar? Wer entscheidet?

## Risiken
Was kann schiefgehen? Wie wird das mitigiert?
```

## Stil

- **Deutsch.** Code-Snippets in englischen Identifiern, sonst durchgehend
  deutsch — siehe [Workflow § Sprache](workflow.md#sprache).
- **Knapp.** Eine gute Spec ist 200–500 Zeilen, nicht 2000. Wenn sie
  länger wird, ist sie wahrscheinlich zwei Specs.
- **Konkret.** „Wir verbessern die Performance" ist kein
  Akzeptanzkriterium. „p50-Score-Latenz < 100 ms bei 1000 Signals"
  schon.
- **Verlinkt.** Stütz dich auf vorhandene Artefakte:
  [Architektur](../internals/architecture.md),
  [Datenmodell](../internals/data-model.md),
  [ADRs](../internals/adr/README.md).

## Spec vs. ADR

| Artefakt | Frage |
|---|---|
| Spec | „Was bauen wir als Nächstes?" |
| ADR | „Warum haben wir diese Entscheidung getroffen, die wahrscheinlich noch öfter relevant wird?" |
| `docs/internals/*.md` | „Wie funktioniert das laufende System?" |

Eine Spec darf auf einen geplanten ADR verweisen. Sobald die
Entscheidung mehrfach auftaucht, wird sie als eigener ADR
festgehalten — siehe [ADR-Index](../internals/adr/README.md).

## Workflow

1. Brainstorming (mit `superpowers:brainstorming`) ergibt eine
   Design-Idee.
2. Spec schreiben, im PR gegen `main` reviewen lassen
   (`feat/p<phase>-t<task>-spec`-Branch).
3. Nach Approval: [Plan](plan-format.md) als separater PR mit
   Checkbox-Tasks.
4. Implementierungs-PRs referenzieren beide via `Updates #<issue>` und
   dem Spec-/Plan-Link im PR-Body.
