# Konfiguration

Zwei Stellen — eine Datei pro Repo und ein paar Environment-Variablen.

## `.lodestone.yaml`

Optional im Repo-Wurzelverzeichnis. Wird durch `lodestone init` mit
Defaults angelegt. Alle Felder sind optional; nicht-gesetzte Felder
werden aus dem Default zugemischt.

```yaml
goals:
  - reliability
  - speed
  - shipping

tech_interests:
  - mcp
  - llm-tools
  - agent-frameworks

lodestone:
  min_stars: 50                # github-trending Stars-Filter
  min_age_days: 30             # github-trending: nur Repos pushed innerhalb von N Tagen
  max_last_commit_age_days: 180
  require_license: true
```

### Felder

| Pfad | Typ | Default | Wirkung |
|---|---|---|---|
| `goals` | `[]string` | `[]` | Wird in den Fingerprint übernommen; LLM-Planning ab Phase 2 priorisiert Goal-Matches. |
| `tech_interests` | `[]string` | `[]` | Wie `goals`, aber für technische Themen (z. B. `mcp`, `vector-db`). |
| `lodestone.min_stars` | `int` | `50` | `github_trending`-Filter. Anti-Hype-Default. |
| `lodestone.min_age_days` | `int` | `30` | `github_trending`-Filter: nur Repos, die innerhalb der letzten N Tage gepushed wurden. |
| `lodestone.max_last_commit_age_days` | `int` | `180` | Schwelle für Risk=`high` (stale). |
| `lodestone.require_license` | `bool` | `true` | Risk=`high` wenn Lizenz fehlt; Recommendations werden nicht herausgefiltert, aber als hochrisikant gewertet. |

Hintergrund zu den Default-Werten:
[ADR-0005 — Anti-Hype-Defaults](../internals/adr/0005-anti-hype-defaults.md).

## Environment-Variablen

| Variable | Zweck | Default |
|---|---|---|
| `GITHUB_TOKEN` | Optionaler Token für `github_trending` — entlastet das Rate-Limit erheblich. | leer |
| `LODESTONE_MOCK_FIXTURES` | Verzeichnis mit `<source>.json`-Fixtures für `lodestone ingest --mock`. Nur für Tests und Offline-Demos. | leer |

## Modell-Routing

Datei: `base/models.yaml` (im Lodestone-Repo, nicht im User-Projekt).
Definiert, welches Claude-Modell für welche Rolle eingesetzt wird:

| Rolle | Modell |
|---|---|
| `planning` | `claude-opus-4-7` (Specs, Pläne, Architektur) |
| `default` | `claude-sonnet-4-6` (Implementierung, Reviews) |
| `mechanical` | `claude-haiku-4-5-20251001` (Format-Konvertierung, Rationale) |
| `review` | `claude-sonnet-4-6` (PR-Review, Spec-Critique) |

Über `--model <id>` kann pro Aufruf überschrieben werden — siehe
[`plan`](commands/plan.md) und [`apply`](commands/apply.md).
Rollen-Details: [`AGENTS.md`](../../AGENTS.md).
