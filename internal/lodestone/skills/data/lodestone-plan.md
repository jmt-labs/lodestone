---
name: lodestone-plan
description: Generiert Spec und Plan im superpowers-Format aus einer Lodestone-Recommendation. Nutze diesen Skill, wenn der User eine konkrete Recommendation in einen ausführbaren Plan überführen will ("plan für diesen Vorschlag", "spec für rec-id X"). Setzt eine bestehende Recommendation in `.lodestone/recommendations.jsonl` voraus.
---

# Lodestone Plan

Du nutzt die Planning-Engine via `lodestone plan` und überprüfst den Output.

## Workflow

1. Falls der User keine Rec-ID nennt, erst `lodestone-recommend` aufrufen bzw. auf vorhandene Recommendations hinweisen.
2. Führe aus:
   ```
   lodestone plan <rec-id>
   ```
   Das Binary ruft Claude über die CLI mit der `planning`-Modell-Wahl (`claude-opus-4-7`), parst die Antwort an `===SPEC===`/`===PLAN===`-Markern und schreibt zwei Dateien unter `docs/superpowers/specs/` und `docs/superpowers/plans/`.
3. Öffne beide erzeugten Dateien und lies sie zur Sanity-Check. Achte auf:
   - YAGNI-Verletzungen (spekulative Features)
   - fehlende Tradeoff-Analyse
   - Plan-Tasks ohne klare Akzeptanzkriterien
4. Wenn Auffälligkeiten: dem User kurz melden und Anpassungsvorschläge machen, **bevor** der Plan gebrancht wird.

## Konventionen

- Antworte deutsch.
- Verwende `--dry-run`, wenn der User nur den Prompt sehen möchte (z. B. zum Eigen-Audit).
- Nie selbst Spec/Plan freischalten — der User entscheidet, was nach `docs/superpowers/` final eincheckt.
