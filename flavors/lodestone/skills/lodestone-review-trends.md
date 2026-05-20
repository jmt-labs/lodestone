---
name: lodestone-review-trends
description: Periodischer Rückblick über die letzten Lodestone-Signale, Recommendations und Decisions. Nutze diesen Skill für Quartals-/Monats-Reviews oder wenn der User "review trends", "was haben wir verpasst" oder "lodestone-stats" anfragt.
---

# Lodestone Review Trends

Du erstellst einen knappen Review-Report über einen frei wählbaren Zeitraum
(Default: letzte 30 Tage).

## Workflow

1. Aggregiere:
   - `.lodestone/signals.jsonl` — Anzahl je Quelle, durchschnittliche Stars
   - `.lodestone/recommendations.jsonl` — Verteilung über Effort und Risk
   - `.lodestone/decisions.log` — wann wurde welcher Verb ausgeführt
2. Erkenne Muster:
   - Welche Quelle liefert die wertvollsten Recommendations (compat ≥ 0.7)?
   - Gibt es Recommendations, die wiederholt erschienen sind aber nie zu einem `lodestone plan` führten?
   - Wurden seit dem letzten Review neue Frameworks erkannt (Fingerprint-Diff)?
3. Schreibe einen kurzen Markdown-Report nach `docs/lodestone-review-YYYY-MM.md`.

## Konventionen

- Antworte deutsch.
- Keine LLM-Aufrufe für Aggregation — nutze reine Datei-Operationen.
- Halte den Report unter 80 Zeilen.
- Verweise auf konkrete Rec-IDs für Folge-`plan`-Aufrufe.

## Output-Format

```
# Lodestone Review — YYYY-MM

## Signal-Volumen
- github_trending: N Signale, ⌀ Y Stars
- hackernews: …
…

## Aktivierte Recommendations
- <rec-id>: <Titel>, compat=X, → plan vom YYYY-MM-DD

## Vergessene Recommendations (>= 2× erschienen, kein Plan)
- <rec-id>: <Titel>

## Fingerprint-Drift
- neue Frameworks: …
- entfernte: …
```
