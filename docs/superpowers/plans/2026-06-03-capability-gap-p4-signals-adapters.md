# Capability-Gap Phase 4 — Reicher Signal-Content + Adapter

**Spec:** [`../specs/2026-06-03-capability-gap-design.md`](../specs/2026-06-03-capability-gap-design.md) §C
**Abhängigkeit:** Phase 2 (Enrichment-Qualität profitiert vom reicheren
Content). Eigenständig mergebar.
**Charakter:** **null neue Dependencies** (Stdlib `net/http` + bestehende
Helfer), konsistent mit ADR-0005.

## Ziel

`analyze` bekommt mehr Material als nur Signal-Titel, und die `analyze`-
Engine wird über MCP-Server und einen Skill ansprechbar (dünne Adapter,
ADR-0001).

## Tasks

- [ ] **P4-T1** — `schema/signal.go`: optionales `ReadmeExcerpt` /
      `Capabilities` (`omitempty`), `SignalSchemaVersion` 1 → 2. Tests.
- [ ] **P4-T2** — `ingest/github_trending.go`: README via
      `GET /repos/{owner}/{repo}/readme` (gleiche Auth wie heute),
      gekürzt + via `stripTags` bereinigt; Caching greift über
      `cachePath`/`saveCache`/`loadCache`. Tests mit HTTP-Mock.
- [ ] **P4-T3** — `ingest/changelog.go`: `parseChangelogHTML` erweitern,
      um die ersten Absätze pro Entry in `Summary` zu erfassen. Tests.
- [ ] **P4-T4** — `mcp/tools.go`: neues Tool `analyze` als dünner Adapter
      auf die `analyze`-Engine (ADR-0001). `server_test.go` ergänzen.
- [ ] **P4-T5** — `flavors/lodestone/skills/`: Skill für gap-getriebene
      Tool-Adoption + Embed-Mirror in `internal/lodestone/skills/data/`
      (erfüllt `make skills-coverage`, ADR-0007).
- [ ] **P4-T6** — `make test lint vuln e2e` grün; E2E-Fixture unter
      `e2e/fixtures/signals` mit AI-Repo + GitHub-Tool-Signal, das ein
      reproduzierbares Enrichment erzeugt.

## Definition of Done

- Ein GitHub-Signal trägt einen `readme_excerpt`; das daraus erzeugte
  Enrichment referenziert konkrete Tool-Capabilities.
- MCP-Tool `analyze` liefert dasselbe Ergebnis wie die CLI.

## Modus

Feature-Branch `feat/p4-t<task>-<slug>`, TDD, atomare Commits, PR mit
`Updates #<epic>`.
