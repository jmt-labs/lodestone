# ADR-0002 — JSONL statt SQLite als Persistierung

## Status

Accepted, 2026-05-20.

## Kontext

Lodestone muss Signals, den Fingerprint und Recommendations lokal
unter `.lodestone/` persistieren. Naive Optionen sind eine
eingebettete Datenbank (SQLite via cgo, BoltDB, Badger) oder
flache JSON-Lines-Dateien.

## Entscheidung

JSON-Lines plus atomic `tmp + rename`:

- `signals.jsonl` — append-only, dedupliziert über `ID`.
- `fingerprint.json` — single-file, komplett-replace.
- `recommendations.jsonl` — komplett-replace via `tmp + rename`.
- `applies.jsonl` — append-only.
- `decisions.log` — append-only Audit-Trail.

## Konsequenzen

- **Plus:** Null externe Dependencies; reine stdlib (`encoding/json`,
  `os`).
- **Plus:** Mit `cat`, `jq`, `grep` direkt inspizierbar — wichtig für
  Audit und Debugging.
- **Plus:** Git-freundlich; Diffs sind menschenlesbar.
- **Plus:** Atomarität über `tmp + rename` ist OS-garantiert (POSIX).
- **Minus:** Keine Indizes — In-Memory-Lookup bei jedem Lauf. Für die
  erwartete Größenordnung (≤ 10.000 Signals) vernachlässigbar.
- **Minus:** Keine Transaktionen über mehrere Files — wird kompensiert
  durch klar isolierte Schreib-Pfade pro Datei.

## Alternativen

- **SQLite (cgo).** Cross-Compile schwierig, bricht das
  „pure-Go-Toolchain"-Versprechen. Verworfen.
- **modernc.org/sqlite (pure Go).** 7 MB Binary-Wachstum für eine
  Pipeline, die nie >10k Records sieht. Verworfen.
- **BoltDB / Badger.** Custom-Format, nicht inspizierbar mit Standard-
  Tools. Verworfen.

## Quelle

[Phase-1-Design](../../superpowers/specs/2026-05-20-lodestone-design.md)
§ Storage, `internal/lodestone/store/`.
