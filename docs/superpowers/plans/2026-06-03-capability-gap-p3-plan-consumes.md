# Capability-Gap Phase 3 — `plan` konsumiert Enrichment

**Spec:** [`../specs/2026-06-03-capability-gap-design.md`](../specs/2026-06-03-capability-gap-design.md) §B
**Abhängigkeit:** Phase 2 (braucht `Enrichment` + `EnrichmentStore`).
Eigenständig mergebar (Enrichment-Block ist optional im Prompt).

## Ziel

Die Plan-Generierung wird gezielt und detailliert: liegt für eine rec-id
ein Enrichment vor, fließen Capability-Gap, Migrationspfad und betroffene
Dateien in Spec und Plan ein. Ohne Enrichment bleibt das Verhalten exakt
wie heute.

## Tasks

- [ ] **P3-T1** — `planning/prompt.go`: `BuildPrompt` um einen
      **optionalen** Enrichment-Block erweitern (Gap, Migrationspfad,
      betroffene Dateien). Tests mit und ohne Enrichment.
- [ ] **P3-T2** — `cmd/lodestone/plan.go`: Enrichment für die rec-id via
      `FindEnrichment` laden und an die Engine reichen; fehlt es →
      heutiges Verhalten.
- [ ] **P3-T3** — `planning/planning.go`: Signatur `Plan(ctx, fp, rec,
      *schema.Enrichment)` nil-tolerant erweitern.
- [ ] **P3-T4** — Test: erzeugte Spec enthält Migrationsschritte, wenn ein
      Enrichment vorhanden ist; ohne Enrichment unverändert.
- [ ] **P3-T5** — `make test lint vuln` grün; `docs/user/commands/plan.md`
      um den Enrichment-Hinweis ergänzen.

## Definition of Done

- `lodestone plan <rec-id>` nach vorherigem `analyze` produziert eine Spec
  mit konkretem Migrationspfad statt generischem Text.
- Ohne `analyze` byte-stabil zum heutigen Output (Regressionstest).

## Modus

Feature-Branch `feat/p3-t<task>-<slug>`, TDD, atomare Commits, PR mit
`Updates #<epic>`.
