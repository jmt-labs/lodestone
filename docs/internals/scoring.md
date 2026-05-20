# Scoring-Algorithmus

`lodestone score` lädt den [Fingerprint](data-model.md#fingerprint) und
alle [Signals](data-model.md#signal) und berechnet pro Signal drei
Dimensionen plus eine deterministische Sortierung. Der Pfad ist
LLM-frei — siehe [ADR-0006](adr/0006-deterministisches-scoring.md).

## Compatibility (0.0 – 1.0)

Gewichtete Jaccard-Ähnlichkeit zwischen den `TopicTags` eines Signals
und den `Frameworks`/`Languages` des Fingerprints.

| Faktor | Gewicht |
|---|---|
| Match auf Language (`Go`, `JavaScript`, …) | **1.5×** |
| Match auf Framework (`cobra`, `react`, …) | **1.0×** |

Formel (vereinfacht):

```
compat = sum(gewichteter_match) / max(|TopicTags|, |Frameworks|+|Languages|)
clamp(compat, 0.0, 1.0)
```

Hintergrund: Language-Matches sind stärker, weil sie die Implementierung
direkt beeinflussen; Framework-Matches sind schwächer, weil ein Signal
auch für einen anderen Framework-Stack relevant sein kann.

## Effort (XS / S / M / L / XL)

Kategorischer Score für den geschätzten Umsetzungsaufwand. Heuristik:

| Bedingung | Effort |
|---|---|
| Match-Count == 0 | `XL` |
| Match-Count ≥ 1 und Stars < 100 | `S` |
| Sonst | `M` |

`XS` und `L` sind reserviert für zukünftige Verfeinerung (LLM-gestützt
ab Phase 5+). `XS` ist heute eine Voraussetzung für Auto-Apply — siehe
[Apply-Befehl](../user/commands/apply.md).

## Risk (low / med / high)

Kategorischer Score basierend auf Maintenance- und License-Signalen:

| Bedingung | Risk |
|---|---|
| Stars ≥ 500 **und** LastCommit < 90 d **und** License vorhanden | `low` |
| License fehlt **oder** LastCommit > 180 d (stale) | `high` |
| Sonst | `med` |

`low` ist Voraussetzung für Auto-Apply.

## Sortierschlüssel

Output `recommendations.jsonl` ist deterministisch sortiert nach:

```
compatibility DESC, stars DESC, id ASC
```

Stars werden über `SignalID` aus `signals.jsonl` aufgelöst. Tie-Breaker
`id ASC` garantiert eindeutige Reihenfolge auch bei identischen
Compatibility- und Stars-Werten.

## Recommendation-ID

Reproduzierbar über identische Inputs:

```
ID = "sha256:" + hex(sha256(signal_id + "|" + json(fingerprint)))
```

Damit erzeugt ein zweiter Score-Lauf mit identischem Fingerprint und
identischer Signal-Liste identische IDs — die Voraussetzung für
[Determinismus](determinism.md).

## Erweiterungspunkte

| Wunsch | Pfad |
|---|---|
| Neuer Compatibility-Faktor | `internal/lodestone/scoring/compatibility.go` |
| Andere Effort-Heuristik | `internal/lodestone/scoring/effort.go` |
| Andere Risk-Schwellen | `internal/lodestone/scoring/risk.go` |
| Andere Sortierung | bricht Determinismus-Test — Spec-Diskussion zuerst |

Jede Änderung muss `TestScoreDeterminism` grün halten und sollte
Golden-Fixtures mit aktualisieren.
