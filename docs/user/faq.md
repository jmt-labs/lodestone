# FAQ

## Welcher Befehl für meinen Use-Case?

| Du willst… | Du brauchst… |
|---|---|
| neue AI-Trends sehen, sortiert nach Stars | `lodestone ingest` + `lodestone signals --top 20` |
| sehen, welche Trends zu **deinem** Repo passen | `lodestone fingerprint` + `lodestone ingest` + `lodestone score` + `lodestone signals --top 10` |
| eine Recommendation in einen Spec + Plan überführen (mit Claude) | `lodestone plan <rec-id>` |
| eine Recommendation **direkt** als Draft-PR einspielen | `lodestone apply <rec-id>` (vier Safety-Gates müssen passen) |
| einen Apply zurückrollen | `lodestone undo <branch-or-rec-id>` |
| Audit-Trail-Statistik sehen | `lodestone memory` (konsolidiert) und `lodestone stats` (Apply-States) |
| Lodestone in einem neuen Projekt aufsetzen | `lodestone init` |
| Lodestone in Claude Desktop / Claude Code nutzen | [MCP-Server-Setup](mcp-server.md) |

## Brauche ich die Claude-CLI?

Nur für `lodestone plan` und `lodestone apply` (intern ruft `apply`
ebenfalls die Planning-Engine auf). Phase 1 (`fingerprint`, `ingest`,
`score`, `signals`) läuft komplett ohne LLM und ohne API-Keys. Siehe
[ADR-0003](../internals/adr/0003-claude-cli-shellout.md).

## Was kostet das?

Phase 1 ist LLM-frei und damit kostenlos. Phase 2 ruft die `claude`-CLI
auf — die Kosten richten sich nach deinem Anthropic-Plan, nicht nach
Lodestone. Die [Anti-Hype-Defaults](../internals/adr/0005-anti-hype-defaults.md)
sorgen dafür, dass nur wenige, gut gefilterte Recommendations
entstehen — du bestimmst selbst, welche du in `plan` oder `apply`
weiterleitest.

## Privacy?

Alles bleibt lokal:

- `.lodestone/` lebt in deinem Repo.
- Telemetrie wird nicht gesendet.
- Source-Adapter holen ausschließlich öffentliche Daten von GitHub,
  HackerNews, arXiv und den genannten Changelogs.
- Cross-Repo-Sharing (Phase 5+) ist opt-in und muss explizit aktiviert
  werden — Privacy-Spec liegt unter
  [`docs/superpowers/specs/2026-05-20-lodestone-sharing-privacy.md`](../superpowers/specs/2026-05-20-lodestone-sharing-privacy.md).

## Lässt sich Lodestone in CI automatisieren?

Ja, über die GitHub-Action unter
`.github/workflows/templates/lodestone-weekly.yml`. Das Template läuft
Sonntag 03:00 UTC und auf `workflow_dispatch`, ist aber opt-in
(es wird nicht aktiviert, bevor du es kopierst).

## Warum kein Daemon, kein Auto-Push auf `main`?

Konservative Design-Entscheidung. Lodestone ist explizit kein
Hintergrund-Prozess und kein Auto-Editor — siehe
[`CLAUDE.md`](../../CLAUDE.md) § „Was lodestone NICHT ist". Auto-PRs
(Phase 4) sind harten Schranken unterworfen
([ADR-0008](../internals/adr/0008-apply-safety-gates.md)).

## Funktioniert Lodestone für Nicht-Go-/Nicht-Node-Projekte?

Heute werden Go und Node am besten erkannt. Andere Sprachen liefern
einen schwächeren Fingerprint (Languages-Erkennung über
Dateierweiterungen funktioniert, Framework-Heuristik nicht). Eine
Erweiterung pro Sprache geht über
`internal/lodestone/fingerprint/<language>.go` analog zu `golang.go` und
`node.go` — siehe [Architektur § Erweiterungspunkte](../internals/architecture.md).

## Wo finde ich die Roadmap?

[`docs/internals/roadmap.md`](../internals/roadmap.md) ist die Single
Source of Truth für den Phasen-Status.
