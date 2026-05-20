# MCP-Server `lodestone-mcp`

Zweites Binary neben `lodestone`. Spricht
[Model Context Protocol](https://modelcontextprotocol.io) über
**stdio JSON-RPC 2.0** (Protocol-Version `2024-11-05`) und stellt
Lodestones Kern-Funktionen als Tools für MCP-Clients bereit.

## Installation

```sh
go install github.com/jmt-labs/lodestone/cmd/lodestone-mcp@latest
```

Verifikation:

```sh
lodestone-mcp --version
# lodestone-mcp v0.1.0-alpha
```

Beide Binaries (`lodestone` und `lodestone-mcp`) sind ab `v0.1.0` in
jedem GoReleaser-Archive enthalten.

## Setup für Claude Desktop

Datei: `~/Library/Application Support/Claude/claude_desktop_config.json`
(macOS) bzw. `%APPDATA%\Claude\claude_desktop_config.json` (Windows).

```json
{
  "mcpServers": {
    "lodestone": {
      "command": "lodestone-mcp",
      "cwd": "/absoluter/pfad/zu/deinem/projekt"
    }
  }
}
```

`cwd` ist wichtig — der Server liest und schreibt nach
`<cwd>/.lodestone/`. Nach dem Speichern Claude Desktop neu starten.

## Setup für Claude Code

Datei: `.mcp.json` im Projekt-Wurzelverzeichnis.

```json
{
  "mcpServers": {
    "lodestone": {
      "command": "lodestone-mcp"
    }
  }
}
```

Claude Code übergibt seinen `cwd`. Wenn das Binary nicht im `$PATH`
liegt, vollständigen Pfad angeben (`"command": "/Users/.../bin/lodestone-mcp"`).

## Tools

Fünf Built-in-Tools rufen direkt in die Phase-1/2-Pakete — kein
Duplikat-Code zur CLI (siehe
[ADR-0001](../internals/adr/0001-three-layer-model.md)).

| Tool | Argumente | Output |
|---|---|---|
| `list_signals` | `{source?, since?, top?}` | JSON-Array von Signals |
| `query_trends` | `{since?}` | `{count_by_source, avg_stars, total}` |
| `score_repo` | `{}` | `{fingerprint_summary, top_recommendations}` |
| `generate_plan` | `{rec_id, model?}` | `{spec_md, plan_md, spec_path, plan_path, model}` |
| `record_decision` | `{verb, outcome, detail?, args?}` | `{ok: true}` |

Schema-Details aller Tools sind in `internal/lodestone/mcp/tools.go`
hinterlegt und werden zur Laufzeit über `tools/list` ausgespielt.

## Verifikation

In Claude Desktop oder Claude Code, in einem Repo mit `.lodestone/`-Daten:

> „List die Top-5 Signals aus lodestone."

Der Client sollte das `list_signals`-Tool mit `{top: 5}` aufrufen und
das Ergebnis wiedergeben.

## Troubleshooting

| Symptom | Ursache / Lösung |
|---|---|
| Server startet nicht | `cwd` zeigt auf nicht-existierendes Verzeichnis. Pfad korrigieren. |
| `tools/list` zeigt 0 Tools | Falsches Binary aufgerufen. `which lodestone-mcp` prüfen. |
| `score_repo` schlägt fehl | `.lodestone/fingerprint.json` und `.lodestone/signals.jsonl` müssen existieren — `lodestone fingerprint` und `lodestone ingest` zuerst laufen lassen. |
| `generate_plan` schlägt fehl | `claude`-CLI fehlt im Pfad oder kein `--model`-Argument; siehe [`plan`-Befehl](commands/plan.md). |
| Logs anschauen | Claude Desktop: `Help → View Logs`. Claude Code: stderr des MCP-Subprozesses. |

## Verwandt

- [Architektur](../internals/architecture.md) — Drei-Ebenen-Modell.
- [ADR-0001](../internals/adr/0001-three-layer-model.md) — warum drei
  Frontends auf demselben Kern.
- [Befehle](commands/README.md) — die CLI-Verben, die identische
  Funktionen anbieten.
