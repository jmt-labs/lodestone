# `lodestone ingest`

Holt Signale von externen Quellen und schreibt sie deduppliziert nach
`.lodestone/signals.jsonl`. Zweite Pipeline-Stufe — kann parallel zu
`fingerprint` laufen.

## Synopsis

```sh
lodestone ingest [--source <name>...] [--mock] [--root <pfad>]
```

## Flags

| Flag | Default | Zweck |
|---|---|---|
| `--source` | alle | Spezifische Quelle(n); mehrfach setzbar |
| `--mock` | false | Offline-Modus aus `$LODESTONE_MOCK_FIXTURES` |
| `--root` | `$PWD` | Projekt-Wurzel |

## Verhalten

Sechs Sources, alle teilen Cache- und Retry-Helper:

| Source | Endpoint | Default-Filter |
|---|---|---|
| `github_trending` | `api.github.com/search/repositories` | `stars:>=50 pushed:>{30daysago}`; liest `$GITHUB_TOKEN` |
| `hackernews` | `hacker-news.firebaseio.com/v0/topstories.json` | Story-Type, Keywords `ai, llm, mcp, claude, agent`, Limit 50 |
| `arxiv` | `export.arxiv.org/api/query` | `cat:cs.AI`, sortBy=submittedDate, max=30 |
| `anthropic_changelog` | HTML-Scrape `docs.anthropic.com/en/release-notes/api` | `<h2>`/`<h3>` mit `YYYY-MM-DD`, max 30 |
| `openai_changelog` | HTML-Scrape `platform.openai.com/docs/changelog` | wie oben |
| `npm_trending` | `registry.npmjs.org/-/v1/search` | `keywords:ai`, `popularity=1.0`, size=20 |

**Cache:** Pro Source und Tag genau ein Cache-File unter
`.lodestone/cache/<source>-YYYY-MM-DD.json`. Zweiter Aufruf am selben
Tag liest aus dem Cache — ideal für `score`-Iterationen ohne neue
HTTP-Calls.

**Retry:** Generischer `retryFetch[T]` mit drei Versuchen,
exponentiellem Backoff, retryable bei 5xx und 429.

**Deduplikation:** Über Signal-ID (sha256 aus Source + URL).

## Beispiele

```sh
# Alle Sources
lodestone ingest

# Nur eine
lodestone ingest --source github_trending

# Mehrere selektiv
lodestone ingest --source github_trending --source hackernews

# Offline aus Fixtures
LODESTONE_MOCK_FIXTURES=./e2e/fixtures/signals lodestone ingest --mock
```

## Exit-Codes & Fehler

| Code | Bedeutung |
|---|---|
| 0 | Erfolg |
| ≠0 | Source nicht bekannt, Cache-Schreib-Fehler, `--mock` ohne `$LODESTONE_MOCK_FIXTURES` |

Häufige Fehler:

- `Error: unknown source "X"` — Source-Name prüfen, siehe Tabelle oben.
- `Error: github_trending: max retries exceeded` — Rate-Limit oder
  Netzwerk; `$GITHUB_TOKEN` setzen reduziert die Rate-Limit-Wahrscheinlichkeit
  erheblich.
- `Error: --mock requires $LODESTONE_MOCK_FIXTURES` — Env-Var muss auf
  ein Verzeichnis mit `<source>.json`-Fixtures zeigen.

## Verwandt

- [`fingerprint`](fingerprint.md) — Repo-Profil parallel erzeugen.
- [`score`](score.md) — Signals × Fingerprint scoren.
- [`signals`](signals.md) — gespeicherte Signals filtern und anzeigen.
- [Troubleshooting](../troubleshooting.md) — typische Ingest-Fehler.
