# ADR-0005 — Konservative Anti-Hype-Defaults

## Status

Accepted, 2026-05-20.

## Kontext

Lodestone soll Aufmerksamkeit auf reife AI-Ökosystem-Trends lenken,
nicht auf jeden frischen GitHub-Trend. Ohne klare Schwellen würden
Stars-getriebene Hype-Repos ohne Maintenance-Spur die Top-N
dominieren.

## Entscheidung

Konservative Defaults in `.lodestone.yaml`, zugemischt wenn nicht
explizit gesetzt:

| Feld | Default | Zweck |
|---|---|---|
| `min_stars` | 50 | Filtert reine Hobby-Projekte |
| `min_age_days` | 30 | Filtert Repos, die noch keinen Push-Zyklus überlebt haben |
| `max_last_commit_age_days` | 180 | Filtert verwaiste Projekte |
| `require_license` | true | Excludes Repos ohne Lizenz aus den Empfehlungen |

Diese Defaults sind in `.lodestone.yaml` explizit überschreibbar; das
Tool weist beim ersten `init` darauf hin, aber drängt nicht zu lockerer
Einstellung.

## Konsequenzen

- **Plus:** Die Top-10-Liste enthält per Default keine frischen
  Hype-Pakete ohne Maintenance.
- **Plus:** Risk-Score `low` setzt License voraus — Anti-Hype-Default
  und Risk-Logik bleiben konsistent.
- **Plus:** Kein „Day-1-Trend-Bias" — ein Repo muss mindestens einen
  Monat Bestand haben, um aufzutauchen.
- **Minus:** Tagesfrische Releases bekannter Maintainer werden ggf.
  unterdrückt. User können `min_age_days: 0` setzen, wenn sie das
  wollen.

## Alternativen

- **Keine Defaults, alles per User-Konfig.** Verworfen — der erste
  Lauf wäre dann eine Stars-Liste, also genau das Anti-Pattern.
- **LLM-basiertes Maturity-Scoring.** Verworfen — würde Phase-1-
  Determinismus brechen, siehe
  [ADR-0006](0006-deterministisches-scoring.md).

## Quelle

[CLAUDE.md § Phase-1-Invarianten](../../../CLAUDE.md),
[Phase-1-Design](../../superpowers/specs/2026-05-20-lodestone-design.md),
`internal/config/config.go`.
