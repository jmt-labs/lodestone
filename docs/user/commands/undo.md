# `lodestone undo <branch-or-rec-id>`

Rollt einen Apply zurück: schließt den zugehörigen PR und löscht den
Branch lokal sowie remote.

## Synopsis

```sh
lodestone undo <branch-or-rec-id> [--root <pfad>]
```

## Flags

| Flag | Default | Zweck |
|---|---|---|
| `<branch-or-rec-id>` | — | Branch-Name (`lodestone/apply-…`) **oder** Recommendation-ID |
| `--root` | `$PWD` | Projekt-Wurzel |

## Verhalten

1. Apply-State aus `.lodestone/applies.jsonl` über den Identifier
   auflösen.
2. `gh pr close --delete-branch <pr-number>` ausführen.
3. Lokalen Branch löschen.
4. Apply-State auf `undone` setzen.
5. Audit-Eintrag in `.lodestone/decisions.log`.

Funktioniert nur für Applies, die `lodestone apply` selbst angelegt
hat — manuell erstellte PRs werden nicht angerührt.

## Beispiele

```sh
# Per Branch
lodestone undo lodestone/apply-abc123-2026-05-20

# Per Recommendation-ID
lodestone undo sha256:abc123…
```

## Exit-Codes & Fehler

| Code | Bedeutung |
|---|---|
| 0 | Erfolg (PR geschlossen, Branch gelöscht) |
| ≠0 | Identifier nicht gefunden, `gh` nicht im PATH, PR bereits geschlossen |

## Verwandt

- [`apply`](apply.md) — Recommendation als Draft-PR.
- [`stats`](stats.md) — Status-Überblick (`undone` taucht dort auf).
