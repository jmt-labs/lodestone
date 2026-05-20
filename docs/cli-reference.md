# lodestone — CLI-Reference

Vollständige Subkommando-Übersicht je Phase. Für eine erzählende
Einführung siehe den [User-Guide](lodestone.md).

## Phase 1 — Deterministische Pipeline

| Verb | Status | Zweck |
|---|---|---|
| `lodestone fingerprint` | ✅ | Aktuelles Repo analysieren → `.lodestone/fingerprint.json` |
| `lodestone ingest` | ✅ | Externe Signale abrufen → `.lodestone/signals.jsonl` |
| `lodestone score` | ✅ | Signale × Fingerprint scoren → `.lodestone/recommendations.jsonl` |
| `lodestone signals` | ✅ | Gespeicherte Signale filtern / anzeigen |

## Phase 2 — Planning & Onboarding

| Verb | Status | Zweck |
|---|---|---|
| `lodestone init` | ✅ | `.lodestone.yaml` + `.gitignore`-Snippet + Skills nach `.claude/skills/` |
| `lodestone plan <rec-id>` | ✅ | Spec / Plan aus Recommendation generieren (ruft `claude`-CLI) |

## Phase 3 — Remote-Schnittstellen

| Verb | Status | Zweck |
|---|---|---|
| `lodestone memory` | ✅ | `decisions.log` → `.claude/memory.json` konsolidieren |
| `lodestone-mcp` (2. Binary) | ✅ | MCP-Server über stdio (`list_signals`, `query_trends`, `score_repo`, `generate_plan`, `record_decision`) |

## Phase 4 — Auto-PR-Engine

| Verb | Status | Zweck |
|---|---|---|
| `lodestone apply <rec-id>` | ✅ | Recommendation als Draft-PR (vier Safety-Gates) |
| `lodestone undo <branch>` | ✅ | Apply rückgängig: PR schließen + Branch löschen |
| `lodestone stats` | ✅ | Apply-Erfolgs-Statistik aus `.lodestone/applies.jsonl` |

## Phase 5+ (geplant)

| Verb | Status | Zweck |
|---|---|---|
| `lodestone recommend` | Skill (Phase 2) | Empfehlungen interaktiv durchgehen — `flavors/lodestone/skills/lodestone-recommend.md` |
| `lodestone calibrate` | Stub | Scoring-Gewichte gegen Decision-Log nachjustieren |
| `lodestone share` | Stub | Cross-Repo-Sharing (siehe [Privacy-Spec](superpowers/specs/2026-05-20-lodestone-sharing-privacy.md)) |

## Globale Flags

| Flag | Default | Zweck |
|---|---|---|
| `--root` | `$PWD` | Projekt-Wurzel, in der `.lodestone/` lebt |

## Ingest-Quellen

| Name | API | Default-Filter |
|---|---|---|
| `github_trending` | `api.github.com/search/repositories` | `stars:>=50 pushed:>{30daysago}` |
| `hackernews` | `hacker-news.firebaseio.com/v0/topstories.json` + `item/<id>` | Type=story, Keywords `ai, llm, mcp, claude, agent`, Limit 50 |
| `arxiv` | `export.arxiv.org/api/query` | `cat:cs.AI`, sortBy=submittedDate, max=30 |
| `anthropic_changelog` | HTML-Scrape `docs.anthropic.com/en/release-notes/api` | `<h2>`/`<h3>` mit `YYYY-MM-DD`, max 30 |
| `openai_changelog` | HTML-Scrape `platform.openai.com/docs/changelog` | wie oben |
| `npm_trending` | `registry.npmjs.org/-/v1/search` | `keywords:ai`, `popularity=1.0`, size=20 |

Alle Quellen nutzen den gleichen Cache (`.lodestone/cache/<source>-<date>.json`)
und Retry-Helper (3 Versuche, exponentieller Backoff, 5xx/429
retryable).

## MCP-Tools (`lodestone-mcp`)

| Tool | Args | Output |
|---|---|---|
| `list_signals` | `{source?, since?, top?}` | JSON-Array von Signals |
| `query_trends` | `{since?}` | `{count_by_source, avg_stars, total}` |
| `score_repo` | `{}` | `{fingerprint_summary, top_recommendations}` |
| `generate_plan` | `{rec_id, model?}` | `{spec_md, plan_md, spec_path, plan_path, model}` |
| `record_decision` | `{verb, outcome, detail?, args?}` | `{ok: true}` |

Protocol: JSON-RPC 2.0, Protocol-Version `2024-11-05`, stdio-Transport.

## Auto-PR-Engine — Safety-Gates

Alle vier Gates müssen passen, sonst lehnt `lodestone apply <rec-id>` ab:

1. `recommendation.risk == low`
2. `recommendation.effort == XS`
3. `recommendation.compatibility >= 0.85`
4. Kein Apply in den letzten 24 h (`.lodestone/applies.jsonl`)
5. `git status` sauber (kein staged/unstaged-Diff)

Branch-Name immer `lodestone/apply-<rec-suffix>-<date>`. PR ist immer
**Draft**, immer gegen `main`, kein Auto-Merge.

## Konfiguration (`.lodestone.yaml`)

```yaml
goals: ["reliability", "shipping"]
tech_interests: ["mcp", "llm-tools"]
lodestone:
  min_stars: 50
  min_age_days: 30
  max_last_commit_age_days: 180
  require_license: true
```

Alle Felder optional — Defaults werden zugemischt. Details:
[User-Guide → Konfiguration](lodestone.md#konfiguration-lodestoneyaml).
