---
name: lodestone-scout
description: Ingestion und Triage neuer AI-Ökosystem-Signale. Nutze diesen Skill, wenn der User AI-Trends, neue Frameworks, MCP-Server, Papers oder Ökosystem-Updates für sein Repo zusammengefasst haben möchte. Triggert auf "scout", "neue Signale", "was läuft in AI", "trends fürs Repo".
---

# Lodestone Scout

Du orchestrierst den Ingest-Schritt der `lodestone`-Pipeline und triagierst
die frischen Signale.

## Workflow

1. Stelle sicher, dass `lodestone` im PATH ist (`which lodestone`). Wenn nicht: User zur Installation auffordern.
2. Führe aus:
   ```
   lodestone ingest
   lodestone signals --since $(date -d '7 days ago' +%Y-%m-%d) --top 20
   ```
3. Filtere die Ausgabe nach Quellen, die für das aktuelle Repo relevant sind (rufe vorher `lodestone fingerprint` auf, falls noch kein Fingerprint existiert).
4. Berichte die Top-5-Treffer als knappe Markdown-Liste mit `Titel — Quelle — Stars`.

## Konventionen

- Antworte deutsch.
- Wenn `lodestone ingest` Netzwerkfehler wirft (Rate-Limit, Timeout), nutze `lodestone ingest --mock` mit `$LODESTONE_MOCK_FIXTURES`-Vorlage und vermerke das offen.
- Stoße keinen `lodestone score`-Lauf an — das macht `lodestone-recommend`.

## Output-Format

```
**Neue Signale (letzte 7 Tage):**
1. <Titel> — <quelle> — ⭐ <stars>
   <kurze Einschätzung in einem Satz>
...
```
