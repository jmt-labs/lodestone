# Capability-Gap Phase 2 — `analyze`-Stufe + Anker-Durchstich

**Spec:** [`../specs/2026-06-03-capability-gap-design.md`](../specs/2026-06-03-capability-gap-design.md) §B
**Abhängigkeit:** Phase 1 (braucht `fp.AI`). Eigenständig mergebar
(`analyze` ist optional in der Pipeline).
**Charakter:** LLM-Stufe **nach** dem Score — Score bleibt deterministisch
(ADR-0010).

## Ziel

`lodestone analyze` reichert die Top-N-Recommendations semantisch an:
Capability-Gap, IST-Ansatz im Repo, betroffene Dateien, Migrationspfad,
Counter-Evidence. Ergebnis getrennt persistiert in
`.lodestone/enrichments.jsonl`. Erster vertikaler Durchstich des
Anker-Use-Case (gap-getriebene Tool-Adoption) zusammen mit Phase 1.

## Tasks

- [ ] **P2-T1** — Paket `internal/lodestone/llm/`: `Runner`,
      `ClaudeRunner`, `FakeRunner` aus `planning/runner.go` extrahieren;
      `planning` re-exportiert die Typen (keine Verhaltensänderung).
      Bestehende `planning`-Tests bleiben grün.
- [ ] **P2-T2** — `schema/enrichment.go`: `Enrichment`-Typ +
      `EnrichmentSchemaVersion`. Roundtrip-Tests.
- [ ] **P2-T3** — `store/store.go` + `filestore.go`: `EnrichmentStore`
      (`ReplaceEnrichments`, `ListEnrichments`, `FindEnrichment(recID)`),
      `enrichments.jsonl` mit atomic-rename. Tests mit `t.TempDir()`.
- [ ] **P2-T4** — `analyze/prompt.go`: deutsches Template, JSON-eingebettete
      Inputs (`fp`, `rec`, `sig`), Antwort als JSON-Block zwischen
      `===ENRICHMENT===`/`===END===`, `SplitResponse` + `json.Unmarshal`.
      Tests: Prompt encodiert `fp.AI` und Signal-Content.
- [ ] **P2-T5** — `analyze/analyze.go`: Engine `New(opts)` /
      `Analyze(ctx, fp, rec, sig)`. FakeRunner-Tests:
      happy-path / bad-json / runner-error.
- [ ] **P2-T6** — `cmd/lodestone/analyze.go`: Cobra-Verb `lodestone
      analyze [--top N] [--model] [--dry-run]`; lädt Top-N aus dem Store,
      matched Signal über `SignalID`, persistiert Enrichments, schreibt
      Audit-Eintrag (`decisions.log`).
- [ ] **P2-T7** — `docs/user/commands/analyze.md` (erfüllt
      `make docs-cmd-coverage`).
- [ ] **P2-T8** — `make test lint vuln` grün.

## Definition of Done

- `lodestone analyze --top 3` erzeugt für ein AI-Repo + GitHub-Tool-Signal
  ein Enrichment mit befülltem `capability_gap` und `migration_path`.
- `--dry-run` gibt den Prompt ohne `claude`-Aufruf aus.
- `recommendations.jsonl` bleibt unverändert (Determinismus sichtbar
  getrennt).

## Modus

Feature-Branch `feat/p2-t<task>-<slug>`, TDD, atomare Commits, PR mit
`Updates #<epic>`.
