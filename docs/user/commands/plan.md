# `lodestone plan <rec-id>`

Generiert aus einer Recommendation ein Spec/Plan-Paar im
superpowers-Format — ruft die lokale `claude`-CLI auf
([ADR-0003](../../internals/adr/0003-claude-cli-shellout.md)).
Voraussetzung: die `claude`-CLI ist installiert und konfiguriert.

## Synopsis

```sh
lodestone plan <rec-id> [--model <id>] [--dry-run]
```

## Flags

| Flag | Default | Zweck |
|---|---|---|
| `<rec-id>` | — | Recommendation-ID (siehe `lodestone signals` oder `recommendations.jsonl`) |
| `--model` | aus `base/models.yaml` (`planning` → `claude-opus-4-7`) | Modell-Override |
| `--dry-run` | false | Nur den Prompt anzeigen, keinen Claude-Aufruf machen |

## Verhalten

Schritte:

1. Recommendation per `<rec-id>` aus
   `.lodestone/recommendations.jsonl` laden.
2. Fingerprint aus `.lodestone/fingerprint.json` laden.
3. Deutschen Architekt-Prompt bauen
   (`internal/lodestone/planning/BuildPrompt`).
4. `claude --print --model <id>` aufrufen.
5. Response an den Markern `===SPEC===` und `===PLAN===` trennen.
6. Spec nach `docs/superpowers/specs/YYYY-MM-DD-<thema>-design.md`
   und Plan nach `docs/superpowers/plans/YYYY-MM-DD-<thema>.md`
   schreiben.
7. Eintrag in `.lodestone/decisions.log` anhängen.

Im `--dry-run`-Modus stoppt der Lauf nach Schritt 3 und gibt nur den
generierten Prompt aus — nützlich zum Debuggen oder zum Kosten-Check.

## Beispiele

```sh
# Top-Recommendation finden
lodestone signals --top 5 --json | jq '.[0].id'

# Spec + Plan generieren
lodestone plan sha256:abc123…

# Nur Prompt-Vorschau, ohne API-Call
lodestone plan sha256:abc123… --dry-run

# Mit Sonnet statt Opus
lodestone plan sha256:abc123… --model claude-sonnet-4-6
```

## Exit-Codes & Fehler

| Code | Bedeutung |
|---|---|
| 0 | Erfolg (Spec + Plan geschrieben) |
| ≠0 | `claude` nicht im `$PATH`, ungültige Recommendation-ID, Marker-Parsing fehlgeschlagen |

Häufige Fehler:

- `claude: command not found` — `claude`-CLI installieren, siehe
  [FAQ](../faq.md).
- `recommendation "<id>" not found` — `lodestone signals --top 20`
  laufen lassen und ID prüfen.

## Verwandt

- [`score`](score.md) — Recommendations erzeugen.
- [`apply`](apply.md) — eine Recommendation direkt als Draft-PR.
- [Modell-Routing](../configuration.md#modell-routing) —
  `base/models.yaml`.
- [ADR-0003](../../internals/adr/0003-claude-cli-shellout.md) — warum
  Shell-out statt SDK.
