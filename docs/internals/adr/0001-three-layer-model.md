# ADR-0001 — Drei-Ebenen-Modell (CLI + Skills + MCP)

## Status

Accepted, 2026-05-20.

## Kontext

Lodestone soll von Menschen am Terminal, von LLMs über Claude-Skills
und von MCP-Clients (Claude Desktop, Claude Code, IDE-Extensions)
gleichzeitig nutzbar sein. Ohne Architektur-Entscheidung droht
Duplikat-Code zwischen drei Codebasen.

## Entscheidung

Drei dünne Frontend-Ebenen rufen denselben Kern auf:

1. **CLI** (`cmd/lodestone/`, Cobra) — primäre Mensch-Schnittstelle.
2. **Claude-Skills** (`flavors/lodestone/skills/`, via `go:embed`
   ins Binary eingebettet) — Anweisungen für LLMs, die die CLI
   benutzen.
3. **MCP-Server** (`cmd/lodestone-mcp/`, stdio JSON-RPC 2.0) — direkte
   Programmschnittstelle für MCP-Clients.

Aller Domänen-Code liegt unter `internal/lodestone/*`. Die drei
Frontends sind reine Adapter-Schichten.

## Konsequenzen

- **Plus:** Eine Code-Quelle, drei Konsumenten — keine Drift.
- **Plus:** Neue Features sind sofort in allen drei Ebenen verfügbar.
- **Plus:** Tests gegen den Kern decken alle drei Ebenen ab.
- **Minus:** Drei Frontend-Repräsentationen müssen synchron gehalten
  werden (CLI-Flag, Skill-Prompt, MCP-Tool-Schema). Mitigiert durch
  Doku-Coverage-Checks (siehe
  [docs-maintenance](../../contributor/docs-maintenance.md)).
- **Minus:** Der MCP-Server bringt zwei Binaries mit sich; GoReleaser
  packt beide in jedes Archive.

## Alternativen

- **Nur CLI:** Lehnt MCP-Clients aus; LLMs müssten Shell ausführen.
  Verworfen, weil MCP der explizite Phase-3-Auftrag war.
- **MCP als Library im CLI-Prozess:** Spart ein Binary, aber bricht
  die stdio-Konvention und macht den CLI-Prozess zustandsbehaftet.
  Verworfen.

## Quelle

[Phase-1-Design](../../superpowers/specs/2026-05-20-lodestone-design.md),
[Phase-3-Design](../../superpowers/specs/2026-05-20-lodestone-phase3-design.md),
[Architektur](../architecture.md).
