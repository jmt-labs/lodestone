# PR-Checkliste

Vor jedem Merge in `main` müssen alle Punkte dieser Liste abgehakt sein.
Eine PR ohne Checkliste wird zurückgewiesen.

## Vor dem Push

- [ ] `make test` grün (mit `-race`, `-count=1`).
- [ ] `make lint` grün (golangci-lint v2).
- [ ] `make vuln` grün (govulncheck).
- [ ] `make e2e` grün (Wegwerf-Repo via mktemp, Determinismus-Diff).

Diese vier sind das CI-Gate — siehe
[Testing](testing.md). Bei Fail: nicht pushen, erst fixen.

## Im PR-Body

- [ ] **Zusammenfassung** in 1–3 Sätzen, fokussiert auf das *Warum*.
- [ ] **Spec-Link** auf `docs/superpowers/specs/…`.
- [ ] **Plan-Task** referenziert (z. B. „T3 aus Phase-2-Plan").
- [ ] **Epic-Issue** via `Updates #N` (oder „kein Epic-Bezug").
- [ ] **Test-Plan** als Markdown-Checkliste.

## Inhaltliche Gates

- [ ] Sprache: **deutsch** für Commit-Messages, PR-Body, Doku, Specs,
      Pläne; englisch für Code-Identifier und API-Felder.
- [ ] Keine neuen Go-Dependencies ohne Spec-Diskussion (siehe
      [Architektur § Dependency-Budget](../internals/architecture.md)).
- [ ] Keine LLM-Calls im Score-Pfad
      ([ADR-0006](../internals/adr/0006-deterministisches-scoring.md)).
- [ ] Bei Bug-Fix: Regressionstest in einem separaten Commit vor dem
      Fix (siehe [Skills-Policy § Regression-First](skills-policy.md#regression-first-regel)).
- [ ] Keine direkt auf `main` gepushten Commits.

## Doku-Synchronität

- [ ] Wenn User-facing Verhalten geändert wurde: betroffene
      `docs/user/commands/*.md` aktualisiert.
- [ ] Wenn neue Phase erreicht oder Phasen-Status geändert wurde:
      [Roadmap](../internals/roadmap.md), [README.md § Status](../../README.md),
      `CHANGELOG.md` synchronisiert.
- [ ] Wenn neuer Begriff eingeführt wurde:
      [Glossar](../user/glossary.md) erweitert.
- [ ] Wenn neue Design-Entscheidung getroffen wurde: ADR unter
      [`docs/internals/adr/`](../internals/adr/README.md) angelegt.

## Coverage

- [ ] Neue `internal/`-Pakete erreichen ≥ 70 % Coverage. Tests
      beschreiben Verhalten, nicht Implementierung.

## Merge

- PRs werden **nur auf explizite User-Aufforderung erstellt**. Push auf
  Feature-Branch ist immer erlaubt; `gh pr create` braucht einen
  expliziten Auftrag.
- Reviewer-Rolle entscheidet über Merge — bei Architektur-Bedenken
  Re-Loop zum Planner.

## Verwandt

- [Workflow](workflow.md), [Testing](testing.md),
  [Skills-Policy](skills-policy.md), [Docs-Maintenance](docs-maintenance.md).
