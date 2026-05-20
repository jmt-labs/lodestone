# Lodestone Phase 2 — Plan

**Spec:** [`../specs/2026-05-20-lodestone-phase2-design.md`](../specs/2026-05-20-lodestone-phase2-design.md)

## Tasks

- [ ] **P2-T1** — Shared Ingest-Helper extrahieren (`cache.go`, `retry.go`).
- [ ] **P2-T2** — ArXiv-Source (Atom-Feed parsing).
- [ ] **P2-T3** — Anthropic-Changelog-Source (HTML-Scrape).
- [ ] **P2-T4** — OpenAI-Changelog-Source (HTML-Scrape).
- [ ] **P2-T5** — npm-Trending-Source (Search-API mit popularity-Sort).
- [ ] **P2-T6** — Subkommando `ingest` registriert die vier neuen Quellen.
- [ ] **P2-T7** — Planning-Engine (`internal/lodestone/planning/`) mit
      `runner`-Interface (Real/Fake), Prompt-Template, Spec/Plan-Output.
- [ ] **P2-T8** — Subkommando `lodestone plan <rec-id>` mit `--dry-run`,
      `--model`-Override.
- [ ] **P2-T9** — Audit-Log (`internal/lodestone/audit/`): JSONL-Append
      pro `lodestone`-Aufruf nach `.lodestone/decisions.log`.
- [ ] **P2-T10** — Subkommando `lodestone init`: `.lodestone.yaml`
      bootstrappen, `.gitignore`-Snippet anhängen, Skill-Install.
- [ ] **P2-T11** — Skills unter `flavors/lodestone/skills/`:
      `-scout`, `-recommend`, `-plan`, `-review-trends`.
- [ ] **P2-T12** — README, CHANGELOG, `docs/lodestone.md` Update.
- [ ] **P2-T13** — E2E-Erweiterung mit Mocks für die vier neuen Quellen
      + `plan`-Pfad mit FakeRunner.

## Modus

Direkt auf `main`. Jeder Task ist ein eigener Commit. CI muss nach jedem
Push grün bleiben.
