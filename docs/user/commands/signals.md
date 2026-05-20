# `lodestone signals`

Filtert und zeigt gespeicherte Signals aus `.lodestone/signals.jsonl`.
Read-only — schreibt nichts.

## Synopsis

```sh
lodestone signals [--top <N>] [--since <YYYY-MM-DD>] [--source <name>] [--json] [--root <pfad>]
```

## Flags

| Flag | Default | Zweck |
|---|---|---|
| `--top` | alle | Nur die ersten N Signals zeigen |
| `--since` | alle | Nur Signals ab Datum |
| `--source` | alle | Nur eine bestimmte Source |
| `--json` | false | JSON-Output statt Tabelle |
| `--root` | `$PWD` | Projekt-Wurzel |

## Verhalten

Liest `.lodestone/signals.jsonl`, filtert nach den Flags und gibt das
Ergebnis aus:

- **Default-Sortierung:** nach Stars (DESC).
- **Tabelle:** Stars · Source · Sprache · Title · URL.
- **JSON:** Array von Signal-Objekten, eine Recommendation pro Eintrag.

## Beispiele

```sh
# Alle Signals, sortiert nach Stars
lodestone signals

# Top 20
lodestone signals --top 20

# Nur ab einem Datum
lodestone signals --since 2026-05-01

# Nur HackerNews
lodestone signals --source hackernews

# JSON für weiterverarbeitung
lodestone signals --top 50 --json | jq '.[] | select(.language == "Go")'
```

## Exit-Codes & Fehler

| Code | Bedeutung |
|---|---|
| 0 | Erfolg |
| ≠0 | `signals.jsonl` fehlt oder beschädigt |

## Verwandt

- [`ingest`](ingest.md) — Signals erst holen.
- [`score`](score.md) — Recommendations daraus erzeugen.
- [Datenmodell § Signal](../../internals/data-model.md#signal) — Felder.
