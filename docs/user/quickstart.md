# Quickstart

In 60 Sekunden vom frischen Klon zur ersten Recommendation.

## Voraussetzungen

- Lodestone ist [installiert](installation.md).
- Du arbeitest in einem Git-Repo (Go- oder Node-Projekt erkennt der
  Fingerprint besonders gut — andere Stacks funktionieren auch, mit
  schwächerem Profil).
- Optional: `$GITHUB_TOKEN` gesetzt (reduziert Rate-Limits beim
  Ingest).

## In 60 Sekunden

```sh
cd dein-projekt/

# 1. Bootstrap: Konfig, .gitignore, vier Skills
lodestone init

# 2. Repo analysieren
lodestone fingerprint

# 3. Externe AI-Signale holen (sechs Quellen)
lodestone ingest

# 4. Signale gegen den Fingerprint scoren
lodestone score

# 5. Top-Empfehlungen ansehen
lodestone signals --top 10
```

Alle Outputs liegen unter `.lodestone/` — siehe
[Artefakte](../internals/artifacts.md).

## Was ist passiert?

| Schritt | Ergebnis | Erklärung |
|---|---|---|
| `init` | `.lodestone.yaml`, `.gitignore`-Snippet, `.claude/skills/*` | Bootstrap eines Lodestone-Projekts. |
| `fingerprint` | `.lodestone/fingerprint.json` | Profil deines Repos (Sprachen, Frameworks, Deps, Goals). Siehe [Glossar § Fingerprint](glossary.md#fingerprint). |
| `ingest` | `.lodestone/signals.jsonl` | Trends aus sechs Quellen, deduppliziert. Siehe [Glossar § Signal](glossary.md#signal). |
| `score` | `.lodestone/recommendations.jsonl` | Deterministisch sortierte Empfehlungen. Siehe [Scoring](../internals/scoring.md). |
| `signals --top 10` | Tabelle auf stdout | Read-only Übersicht. |

## Nächste Schritte

- **Konfiguration anpassen:** [`docs/user/configuration.md`](configuration.md)
  erklärt `.lodestone.yaml`-Felder.
- **Empfehlung umsetzen:** [`plan`](commands/plan.md) erzeugt aus einer
  Recommendation einen Spec/Plan über die `claude`-CLI.
- **Auto-PR:** [`apply`](commands/apply.md) öffnet (bei passenden
  Safety-Gates) einen Draft-PR.
- **MCP-Server:** [`docs/user/mcp-server.md`](mcp-server.md) zeigt das
  Setup für Claude Desktop und Claude Code.
- **Skills aktivieren:** [`docs/user/skills.md`](skills.md) — vier
  Claude-Skills, installiert durch `lodestone init`.

## Fehler beim ersten Lauf?

Siehe [Troubleshooting](troubleshooting.md) — die typischen
„Pipeline-Reihenfolge", „Rate-Limit" und „`claude`-CLI fehlt"-Fälle
sind dort beschrieben.
