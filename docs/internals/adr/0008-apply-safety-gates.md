# ADR-0008 — Vier Safety-Gates für Auto-Apply

## Status

Accepted, 2026-05-20.

## Kontext

Phase 4 erlaubt `lodestone apply <rec-id>`, einen Draft-PR aus einer
Recommendation. Ohne harte Schranken droht ein dauerlaufender
Auto-Editor — genau das, was lodestone laut `CLAUDE.md` **nicht** sein
soll.

## Entscheidung

Vier unabhängige Safety-Gates plus eine git-Hygiene-Prüfung. Alle
müssen passen, sonst lehnt `apply` ab:

1. `recommendation.risk == "low"`.
2. `recommendation.effort == "XS"`.
3. `recommendation.compatibility >= 0.85`.
4. Kein Apply in den letzten 24 h (Cooldown via
   `.lodestone/applies.jsonl`).
5. `git status` sauber — kein staged/unstaged Diff.

Zusätzlich:

- Branch immer `lodestone/apply-<rec-suffix>-<date>`.
- PR ist **immer Draft**, **immer gegen `main`**, **kein Auto-Merge**.
- `lodestone undo <branch-or-rec-id>` rollt zurück (PR schließen +
  Branch löschen).

## Konsequenzen

- **Plus:** Auto-Apply ist auf hochpassende, sehr kleine, sehr
  risikoarme Recommendations beschränkt — der Erwartungswert eines
  Fehlers ist gering, die Recovery-Story (`undo`) trivial.
- **Plus:** Cooldown verhindert Tag-für-Tag-Schwemme.
- **Plus:** `FakeGit` / `FakePR` machen den Pfad ohne echte
  Git-/GitHub-Aufrufe testbar.
- **Plus:** Niemals direkter `main`-Commit; passt zum
  `main`-Push-Verbot aus `CLAUDE.md`.
- **Minus:** Ein „Auto-Apply" ist heute selten, weil die Gates streng
  sind. Genau gewollt für Phase 4 — die Schwelle darf später lockerer
  werden, aber nur über expliziten ADR-Folger.

## Alternativen

- **Nur risk-low als Gate.** Verworfen — ein großes Refactoring kann
  „low-risk" sein und trotzdem nicht für Auto-Apply geeignet.
- **Konfigurierbare Schwellen ohne Hard-Limit.** Verworfen — würde den
  Schutz-Charakter aushöhlen. Konfig wird erst eingeführt, wenn echte
  Apply-Daten zeigen, welche Schwellen praxistauglich sind.

## Quelle

[Phase-4-Design](../../superpowers/specs/2026-05-20-lodestone-phase4-design.md),
[Apply-Befehl](../../user/commands/apply.md), `internal/lodestone/apply/`.
