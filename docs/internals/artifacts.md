# Lokale Artefakte

Alle Outputs landen unter `.lodestone/` im Zielprojekt. Default-mäßig
gehört dieses Verzeichnis in `.gitignore`; einzige Ausnahme ist
`decisions.log` als committable Audit-Spur. `lodestone init` legt den
entsprechenden `.gitignore`-Snippet automatisch an.

## Layout

```
.lodestone/
├── cache/                          # Roh-Fetches mit TTL-Datum
│   ├── github_trending-2026-05-20.json
│   ├── hackernews-2026-05-20.json
│   ├── arxiv-2026-05-20.json
│   ├── anthropic_changelog-2026-05-20.json
│   ├── openai_changelog-2026-05-20.json
│   └── npm_trending-2026-05-20.json
├── signals.jsonl                   # append-only, dedupliziert
├── fingerprint.json                # einzeln, atomar via tmp+rename
├── recommendations.jsonl           # atomar via tmp+rename
├── applies.jsonl                   # Auto-PR-Tracking (Phase 4)
└── decisions.log                   # Audit-Trail, JSONL
```

## Datei-Eigenschaften

| Datei | Format | Schreib-Strategie | Lebenszyklus |
|---|---|---|---|
| `cache/<source>-<date>.json` | JSON | Atomic write | Per Source und Tag genau eine Datei |
| `signals.jsonl` | JSON Lines | Append, dedupliziert über ID | Wächst mit jedem `ingest` |
| `fingerprint.json` | JSON | Atomic via `tmp + rename` | Komplett-Replace bei jedem `fingerprint` |
| `recommendations.jsonl` | JSON Lines | Atomic via `tmp + rename` | Komplett-Replace bei jedem `score` |
| `applies.jsonl` | JSON Lines | Append | Wächst mit jedem `apply` |
| `decisions.log` | JSON Lines | Append-only | Wächst mit jeder Aktion |

## `.gitignore`-Empfehlung

```
.lodestone/
!.lodestone/decisions.log
```

Der Audit-Trail (`decisions.log`) ist standardmäßig ausgenommen, damit
Entscheidungs-Geschichte committable bleibt. Wer auch `decisions.log`
lokal halten will, lässt die zweite Zeile weg.

## Konsolidiertes Memory

`lodestone memory` aggregiert die letzten N Tage aus
`.lodestone/decisions.log` nach `.claude/memory.json` — idempotent über
`(Datum, Verb, Summary)`-Dedup. Default-Fenster: 90 Tage.
