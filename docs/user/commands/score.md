# `lodestone score`

Lädt Fingerprint und alle Signals, berechnet pro Signal eine
Recommendation und schreibt sie deterministisch sortiert nach
`.lodestone/recommendations.jsonl`. Dritte Pipeline-Stufe.

## Synopsis

```sh
lodestone score [--root <pfad>]
```

## Flags

| Flag | Default | Zweck |
|---|---|---|
| `--root` | `$PWD` | Projekt-Wurzel |

## Verhalten

Vorbedingungen:

- `.lodestone/fingerprint.json` muss existieren (vorher `fingerprint`
  laufen lassen).
- `.lodestone/signals.jsonl` muss mindestens einen Eintrag haben
  (vorher `ingest`).

Pro Signal werden drei Dimensionen berechnet — vollständig
deterministisch, ohne LLM:

| Dimension | Werte | Kurz-Regel |
|---|---|---|
| `compatibility` | 0.0 – 1.0 | Gewichtete Jaccard, Language-Match 1.5×, Framework-Match 1.0× |
| `effort` | XS – XL | Default `M`; `XL` bei 0 Match; `S` bei Match + Stars<100 |
| `risk` | low / med / high | `low` bei Stars≥500 ∧ LastCommit<90 d ∧ License; `high` bei fehlender License oder LastCommit>180 d; sonst `med` |

Vollständige Formeln: [Scoring-Algorithmus](../../internals/scoring.md).

**Output:** `.lodestone/recommendations.jsonl`, sortiert nach
`compatibility DESC, stars DESC, id ASC`. Atomar via `tmp + rename`
geschrieben. Recommendation-ID ist
`sha256:hex(signal_id + "|" + json(fingerprint))` — identische Inputs
ergeben identische IDs.

**Determinismus-Garantie:** Zwei aufeinanderfolgende `score`-Läufe
mit identischem Fingerprint und identischer Signal-Liste produzieren
byte-identische Outputs. Verifiziert in Unit-Test
`TestScoreDeterminism` und E2E. Siehe
[Determinismus](../../internals/determinism.md).

## Beispiele

```sh
# Standard-Lauf (Pipeline-Mitte)
lodestone fingerprint && lodestone ingest && lodestone score

# Anschauen
cat .lodestone/recommendations.jsonl | jq .

# Top 5 schnell sehen
lodestone signals --top 5
```

## Exit-Codes & Fehler

| Code | Bedeutung |
|---|---|
| 0 | Erfolg |
| ≠0 | Fingerprint fehlt, keine Signals, Schreibfehler |

Häufige Fehler:

- `Error: no signals in store (run \`lodestone ingest\` first)` —
  Pipeline-Reihenfolge umkehren oder Mock-Modus.
- `Error: read fingerprint (run \`lodestone fingerprint\` first)` —
  vor `score` muss `fingerprint` laufen.

## Verwandt

- [`signals`](signals.md) — Ergebnis als Top-N-Liste anzeigen.
- [`plan`](plan.md) — eine Recommendation zu Spec + Plan ausbauen.
- [`apply`](apply.md) — eine Recommendation zum Draft-PR machen.
- [Scoring-Algorithmus](../../internals/scoring.md) — Formeln.
- [Determinismus](../../internals/determinism.md) — Garantie und
  Verifikation.
