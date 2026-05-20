# `lodestone memory`

Konsolidiert das Audit-Log `.lodestone/decisions.log` zu einer
übersichtlichen `.claude/memory.json`, die Claude-Skills und MCP-Tools
als Kontext nutzen können. Idempotent.

## Synopsis

```sh
lodestone memory [--days <N>] [--out <pfad>] [--root <pfad>]
```

## Flags

| Flag | Default | Zweck |
|---|---|---|
| `--days` | 90 | Wie viele Tage rückwärts berücksichtigen |
| `--out` | `.claude/memory.json` (relativ zu `--root`) | Ziel-Pfad |
| `--root` | `$PWD` | Projekt-Wurzel |

## Verhalten

Liest `.lodestone/decisions.log` (JSONL, append-only), filtert
Einträge der letzten N Tage, aggregiert sie per `(Datum, Verb,
Summary)`-Dedup und schreibt sie nach `.claude/memory.json`. Atomar
via `tmp + rename`.

Idempotenz: zwei aufeinanderfolgende `memory`-Läufe ohne neue
Decisions produzieren identische Outputs.

## Beispiele

```sh
# Standard-Konsolidierung
lodestone memory

# Letzte 30 Tage
lodestone memory --days 30

# In ein anderes Verzeichnis
lodestone memory --out artifacts/memory.json
```

## Exit-Codes & Fehler

| Code | Bedeutung |
|---|---|
| 0 | Erfolg |
| ≠0 | `decisions.log` fehlt oder beschädigt, Schreibfehler |

## Verwandt

- [Artefakte § Konsolidiertes Memory](../../internals/artifacts.md).
- [`stats`](stats.md) — Apply-Statistik aus dem Audit-Trail.
