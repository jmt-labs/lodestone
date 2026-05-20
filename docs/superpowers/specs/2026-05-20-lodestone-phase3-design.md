# Lodestone Phase 3 — Design

**Datum:** 2026-05-20
**Voraussetzung:** Phase 2 abgeschlossen.

## Ziele

Phase 3 macht lodestone **fernsteuerbar** und **periodisch ausführbar**:

1. **`lodestone-mcp`** — separates Go-Binary mit Model-Context-Protocol-
   Server über stdio. Tools: `list_signals`, `query_trends`, `score_repo`,
   `generate_plan`, `record_decision`.
2. **GitHub-Action-Template** — wiederverwendbarer Workflow, der wöchentlich
   `lodestone fingerprint && ingest && score` ausführt, das `.lodestone/`-
   Verzeichnis als PR-Artefakt commitet und einen Summary-PR gegen das
   Ziel-Repo öffnet.
3. **Memory-Konsolidierung** — `lodestone memory` extrahiert „Decisions"
   aus `.lodestone/decisions.log` in `.claude/memory.json` (kurze,
   beschriftete Einträge für LLM-Kontext).

## MCP-Protokoll

Minimal-Stack: JSON-RPC 2.0 über stdio, `initialize` und `tools/call`. Kein
SSE, kein WebSocket — das deckt 95 % der Claude-Code-Integrations-Use-Cases
ab.

**Tools:**

| Tool | Args | Output |
|---|---|---|
| `list_signals` | `{source?, since?, top?}` | `[]Signal` |
| `query_trends` | `{since?}` | `{count_by_source, avg_stars}` |
| `score_repo` | `{root?}` | `{fingerprint_summary, top_recommendations}` |
| `generate_plan` | `{rec_id, model?}` | `{spec_md, plan_md, spec_path, plan_path}` |
| `record_decision` | `{verb, args, outcome, detail?}` | `{ok: true}` |

Jeder Tool-Call ruft direkt in die Phase-1/2-Pakete (`ingest`, `store`,
`scoring`, `planning`, `audit`). Die MCP-Schicht ist dünn.

## Code-Layout

```
cmd/lodestone-mcp/main.go      ← Stdio-MCP-Server
internal/lodestone/mcp/        ← Protocol + Tool-Implementierungen
  server.go
  tools.go
  protocol.go
  server_test.go
internal/lodestone/memory/     ← Memory-Konsolidierung
  memory.go
  memory_test.go
cmd/lodestone/memory.go        ← Subkommando `lodestone memory`
.github/workflows/templates/lodestone-weekly.yml
                                ← reusable workflow template
```

## GitHub-Action-Template

Datei: `.github/workflows/templates/lodestone-weekly.yml`. Wird per
`workflow_call` oder direkt per `cron` (Default: Sonntagnacht UTC)
ausgeführt. Schritte:

1. Checkout
2. Setup Go + install lodestone binary
3. `lodestone fingerprint && lodestone ingest && lodestone score`
4. Diff von `.lodestone/recommendations.jsonl` mit `git diff`
5. Wenn Diff nicht leer: neuer Branch `lodestone/weekly-YYYY-MM-DD`, commit,
   `gh pr create` mit Summary-Body (Top-5 Recommendations).

Das Template ist opt-in: Nutzer copy-paste's in ihren `.github/workflows/`.

## Memory-Konsolidierung

`lodestone memory` liest `.lodestone/decisions.log` der letzten N Tage
(Default: 90) und schreibt strukturierte Memory-Einträge nach
`.claude/memory.json` im Schema `{decisions: [{date, verb, summary}]}`.
Idempotent; bestehende Einträge werden nicht überschrieben.

Wichtig: Keine LLM-Aufrufe — reine Aggregation. Spätere Phasen können
einen LLM-Sub-Pass für Zusammenfassungen ergänzen.

## Test-Strategie

- `internal/lodestone/mcp/server_test.go`: round-trip-Tests pro Tool über
  In-Memory-stdio (Pipes).
- `internal/lodestone/memory/memory_test.go`: tmpdir-Tests.
- `cmd/lodestone-mcp` ist nur ein dünner `main()`-Wrapper, kein eigener Test.
- GitHub-Action-Template wird mit `actionlint` (falls verfügbar) gecheckt;
  ansonsten manuelle Syntax-Prüfung im CI.

## Phase-3-Tasks (kompakt)

- [ ] **P3-T1** — MCP-Protocol-Typen + JSON-RPC-Loop (`mcp/protocol.go`).
- [ ] **P3-T2** — Tool-Registry + 5 Tools (`mcp/tools.go`).
- [ ] **P3-T3** — `cmd/lodestone-mcp/main.go` + Build-Target.
- [ ] **P3-T4** — `lodestone memory`-Subkommando + Package.
- [ ] **P3-T5** — GitHub-Action-Template.
- [ ] **P3-T6** — README / CHANGELOG / docs/lodestone.md.
