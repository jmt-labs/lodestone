# Lodestone Phase 2 — Design

**Datum:** 2026-05-20
**Status:** In Umsetzung (direkt auf `main` autorisiert)
**Voraussetzung:** Phase 1 abgeschlossen (`f623ad4`).

## Ziele

Phase 2 öffnet die Pipeline für **LLM-getriebene Planung** und erweitert
die Signal-Basis auf vier zusätzliche Quellen. Skills werden installierbar,
Decisions persistiert.

## Scope

1. **Vier neue Sources** (`internal/lodestone/ingest/`):
   - `arxiv` — `http://export.arxiv.org/api/query` (Atom XML), Default-
     Query auf `cat:cs.AI` mit `sortBy=submittedDate`
   - `anthropic_changelog` — Scrape `docs.anthropic.com/en/release-notes/api`
   - `openai_changelog` — Scrape `platform.openai.com/docs/changelog`
   - `npm_trending` — `registry.npmjs.org/-/v1/search` (`popularity` sort)
2. **Planning-Engine** (`internal/lodestone/planning/`):
   - `Plan(ctx, fp, rec) (Spec, Plan, error)` ruft Claude über die
     `claude`-CLI (binary-existence check, sonst Fehler)
   - Modell aus `base/models.yaml` (`planning` → `claude-opus-4-7`)
   - Output: zwei Dateien unter `docs/superpowers/specs/` und
     `docs/superpowers/plans/`
   - Prompt-Template inline im Code
3. **Subkommando `lodestone plan <rec-id>`** mit `--dry-run` (zeigt nur
   den Prompt) und `--model <id>`-Override.
4. **Skills** als Markdown-Files unter `flavors/lodestone/skills/`:
   - `lodestone-scout.md` — Ingest + Vorschlags-Triage
   - `lodestone-recommend.md` — Empfehlungen interaktiv durchgehen
   - `lodestone-plan.md` — ruft die Planning-Engine
   - `lodestone-review-trends.md` — Periodischer Trend-Review
5. **`lodestone init`** — neue Subkommando, das:
   - `.lodestone.yaml` mit Defaults anlegt (falls nicht vorhanden)
   - `.gitignore`-Snippet (`.lodestone/` ignorieren, `decisions.log`
     ausnehmen) anhängt
   - Skill-Frontmatter-Files unter `.claude/skills/` symlinkt oder
     kopiert (opt-in)
6. **`decisions.log`** — append-only Audit-Trail. Jeder Aufruf von
   `lodestone plan`, `lodestone score`, `lodestone ingest` schreibt
   einen JSON-Lines-Eintrag mit Timestamp + Verb + Args + Outcome.

## LLM-Integration: Shell-out, nicht SDK

`internal/lodestone/planning/` ruft das `claude`-Binary via
`exec.Command`. Vorteile:

- Keine neue Go-Dependency (anthropic-sdk-go vermieden)
- Nutzt die User-Auth des Claude-CLI (kein eigenes Key-Management)
- Mock-fähig in Tests (`runner` Interface mit `RealRunner` und
  `FakeRunner`)

Falls `claude` nicht im PATH ist, gibt `lodestone plan` einen klaren
Fehler aus und schlägt eine Installationsanleitung vor.

## Dependency-Politik in Phase 2

Phase-1-Invariante "nur Cobra + yaml.v3" wird gelockert:

- `encoding/xml` für ArXiv (Standardbibliothek, kein neuer Dep)
- `regexp` für HTML-Scrape (Standardbibliothek)
- **Keine neuen externen Deps**.

## Test-Strategie

- Jede neue Source hat httptest-basierte Unit-Tests (Erfolg, Empty,
  Timeout, Cache).
- Planning-Engine-Tests verwenden `FakeRunner`, der scripted Output
  zurückgibt.
- E2E erweitert: zusätzliche Mock-Fixtures für alle Quellen.

## Betroffene Dateien (Phase 2)

```
internal/lodestone/ingest/
  arxiv.go              arxiv_test.go
  anthropic_changelog.go anthropic_changelog_test.go
  openai_changelog.go    openai_changelog_test.go
  npm_trending.go        npm_trending_test.go

internal/lodestone/planning/
  planning.go      planning_test.go
  runner.go        runner_test.go
  prompt.go

internal/lodestone/audit/
  audit.go         audit_test.go

cmd/lodestone/
  plan.go
  init.go

flavors/lodestone/skills/
  lodestone-scout.md
  lodestone-recommend.md
  lodestone-plan.md
  lodestone-review-trends.md
```

## Aus Phase 2 ausgeklammert

- MCP-Server (Phase 3)
- GitHub-Action-Template (Phase 3)
- Auto-PR-Engine (Phase 4)
- Cross-Repo-Sharing (Phase 4, eigene Privacy-Spec)
