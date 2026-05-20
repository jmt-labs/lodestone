# Contributing zu lodestone

Danke für dein Interesse! Dieses Projekt folgt einem schlanken
Spec-Plan-PR-Workflow.

## Workflow

1. **Brainstorming** → kurzes Design-Gespräch (Issue oder Chat).
2. **Spec** in `docs/superpowers/specs/YYYY-MM-DD-<thema>-design.md`.
3. **Plan** in `docs/superpowers/plans/YYYY-MM-DD-<thema>.md` mit
   Checkbox-Tasks.
4. **Branch** nach Schema `feat/p<phase>-t<task>-<slug>` (z. B.
   `feat/p1-t3-github-trending`).
5. **TDD**: Test zuerst, dann Implementierung. Bei Bug-Fixes:
   **Regressionstest VOR dem Fix** committen.
6. **PR gegen `main`**, Body referenziert Spec/Plan und das Epic-Issue
   via `Updates #N`. PRs werden **nur auf explizite Aufforderung**
   erstellt.

## Sprachkonvention

- **Doku, Specs, Pläne, Commit-Messages: deutsch.**
- Code-Identifier und API-Felder: englisch.
- Keine unnötigen Kommentare im Code.

## Lokale Entwicklung

```
make build         # baut cmd/lodestone in ./bin
make test          # go test ./...
make lint          # golangci-lint
make vuln          # govulncheck
make e2e           # E2E-Smoke-Test (ab T9)
```

## Code-Stil

- Go 1.24+, Standardbibliothek bevorzugen.
- Phase 1: nur Cobra (CLI) und `gopkg.in/yaml.v3` (Konfig, ab T7).
- **YAGNI** — keine spekulativen Features, keine Abstraktionen
  ohne konkreten zweiten Aufrufer.

## Pflicht-Skills (für Claude-getriebene Arbeit)

| Situation | Skill |
|---|---|
| Vor Feature/Implementierung | `superpowers:brainstorming` |
| Während Implementierung | `superpowers:test-driven-development` |
| Vor Commit/PR | `superpowers:verification-before-completion` |
| Bei Bugs | `superpowers:systematic-debugging` (Regression-Test ZUERST) |
