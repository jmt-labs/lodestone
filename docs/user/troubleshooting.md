# Troubleshooting

Typische Fehler und ihre Ursachen. Wenn dein Symptom hier nicht steht,
hilft oft `--mock` und ein Blick in `.lodestone/decisions.log`.

## Pipeline-Reihenfolge

| Symptom | Ursache / Lösung |
|---|---|
| `Error: no signals in store (run \`lodestone ingest\` first)` | `score` braucht erst Signale. `lodestone ingest` oder `--mock` mit `$LODESTONE_MOCK_FIXTURES`. |
| `Error: read fingerprint (run \`lodestone fingerprint\` first)` | `score` und `apply` brauchen den Fingerprint. Erst `lodestone fingerprint`. |
| `Error: recommendation "<id>" not found` | Recommendation-ID aus `lodestone signals --top 20 --json` holen oder erneut `score` laufen lassen. |

## Ingest

| Symptom | Ursache / Lösung |
|---|---|
| `Error: unknown source "X"` | Source-Name prüfen. Erlaubt: `github_trending`, `hackernews`, `arxiv`, `anthropic_changelog`, `openai_changelog`, `npm_trending`. |
| `Error: github_trending: max retries exceeded` | Rate-Limit oder Netzwerk. `$GITHUB_TOKEN` setzen reduziert die Rate-Limit-Wahrscheinlichkeit erheblich. |
| `--mock requires $LODESTONE_MOCK_FIXTURES` | Env-Var muss auf ein Verzeichnis mit `<source>.json`-Fixtures zeigen — z. B. `e2e/fixtures/signals/`. |

## Planning (Phase 2)

| Symptom | Ursache / Lösung |
|---|---|
| `claude: command not found` | `claude`-CLI installieren oder PATH anpassen. Siehe [FAQ](faq.md). |
| `plan: marker ===SPEC=== not found in response` | Claude-Modell hat das Prompt-Format nicht beachtet. Erneut versuchen, ggf. mit anderem Modell (`--model claude-sonnet-4-6`). |
| Specs/Plans landen am falschen Ort | `--root` zeigt auf das falsche Verzeichnis. Lodestone schreibt nach `<root>/docs/superpowers/{specs,plans}`. |

## Apply (Phase 4)

| Symptom | Ursache / Lösung |
|---|---|
| `apply rejected: risk != low` | Recommendation ist nicht risikoarm genug. Andere Recommendation wählen oder manuell mit `plan` weiterarbeiten. |
| `apply rejected: effort != XS` | Auto-Apply ist auf XS beschränkt. Größere Empfehlungen über `plan` ausarbeiten. |
| `apply rejected: compatibility < 0.85` | Schwelle ist hartcodiert — siehe [ADR-0008](../internals/adr/0008-apply-safety-gates.md). |
| `apply rejected: cooldown active (last apply Xh ago)` | 24h-Sperre. Warten oder `lodestone undo <branch>` aufrufen, um das letzte Apply zu invalidieren (zählt aber als „undone", nicht als Reset des Cooldowns). |
| `apply rejected: git status not clean` | Erst `git stash` oder commit; `apply` will einen sauberen Workdir. |
| `gh: command not found` | GitHub-CLI `gh` installieren und `gh auth login` laufen lassen. |

## MCP-Server

Siehe [MCP-Server § Troubleshooting](mcp-server.md#troubleshooting).

## Allgemein

- **Audit-Trail einsehen:** `tail -f .lodestone/decisions.log` zeigt
  jedes Verb mit Argumenten und Outcome.
- **Determinismus-Verdacht:** Wenn `score` zwei unterschiedliche
  Outputs liefert, ist das ein Bug — bitte gegen
  [Determinismus](../internals/determinism.md) und den Unit-Test
  `TestScoreDeterminism` abgleichen und melden.
- **Cache leeren:** `rm -rf .lodestone/cache/` — `ingest` baut beim
  nächsten Lauf alles neu auf.
