---
name: lodestone-recommend
description: Geht die deterministisch gescorten Recommendations interaktiv durch und priorisiert die Top-N für den User. Nutze diesen Skill, wenn der User nach konkreten Vorschlägen für sein Repo fragt ("was sollte ich nächstes einbauen", "Empfehlungen", "was lohnt sich"). Setzt voraus, dass `lodestone fingerprint` und `lodestone ingest` schon gelaufen sind.
---

# Lodestone Recommend

Du präsentierst die deterministischen Recommendations und hilfst beim
Priorisieren.

## Workflow

1. Stelle sicher, dass `.lodestone/fingerprint.json` und `.lodestone/signals.jsonl` existieren. Falls nicht: User auf `lodestone fingerprint` bzw. `lodestone ingest` hinweisen, dann hier abbrechen.
2. Führe `lodestone score` aus.
3. Lies `.lodestone/recommendations.jsonl` und stelle die Top-5 sortiert nach `compatibility DESC, stars DESC` vor.
4. Für jede Recommendation: zeige `compatibility`, `effort`, `risk` und die Signal-URL. Frage den User, ob er für eine `lodestone plan <rec-id>` triggern möchte.

## Konventionen

- Antworte deutsch.
- **Niemals** automatisch `lodestone apply` triggern — das ist Phase-4-Auto-PR.
- Recommendations mit `compatibility < 0.4` blendest du aus (Anzeige-Schwelle).
- Markiere Recommendations mit `risk: high` deutlich (⚠️).

## Output-Format

```
**Top-Empfehlungen:**

1. ✅ [compat 0.87 · S · low] <Titel>
   URL: <link>
   → `lodestone plan <rec-id>`?
```
