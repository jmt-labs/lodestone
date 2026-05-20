# `lodestone stats`

Aggregiert Apply-State-Counts aus `.lodestone/applies.jsonl`.
Read-only.

## Synopsis

```sh
lodestone stats [--root <pfad>]
```

## Flags

| Flag | Default | Zweck |
|---|---|---|
| `--root` | `$PWD` | Projekt-Wurzel |

## Verhalten

Liest `.lodestone/applies.jsonl` und gruppiert nach Status. Typische
Status-Werte:

| Status | Bedeutung |
|---|---|
| `draft_open` | Draft-PR offen (Standard nach erfolgreichem `apply`) |
| `branch_pushed_no_pr` | Branch gepusht, PR-Erstellung scheiterte |
| `undone` | Per `lodestone undo` zurückgerollt |

Output ist ein einfaches Textformat, eine Zeile pro Status.

## Beispiele

```sh
lodestone stats
# Apply-Stats über 7 Einträge:
#   draft_open             4
#   undone                 2
#   branch_pushed_no_pr    1
```

## Exit-Codes & Fehler

| Code | Bedeutung |
|---|---|
| 0 | Erfolg |
| ≠0 | `applies.jsonl` fehlt oder beschädigt |

## Verwandt

- [`apply`](apply.md) — Wie Einträge in `applies.jsonl` entstehen.
- [`undo`](undo.md) — Wie der Status `undone` entsteht.
