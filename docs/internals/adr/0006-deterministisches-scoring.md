# ADR-0006 — Deterministisches Scoring, keine LLM im Score-Pfad

## Status

Accepted, 2026-05-20.

## Kontext

Recommendations sollen reproduzierbar sein — derselbe Fingerprint plus
dieselbe Signal-Liste muss dieselbe Ausgabe produzieren, sonst sind
Audits, Diffs und Vertrauensaufbau unmöglich.

## Entscheidung

`lodestone score` ist rein deterministisch:

- **Keine LLM-Calls** in Compatibility, Effort, Risk.
- **Stabile Sortierung:** `compatibility DESC, stars DESC, id ASC`.
- **Reproduzierbare IDs:** `sha256:hex(signal_id + "|" + json(fingerprint))`.
- **Verifiziert:** Drei Score-Läufe in
  `internal/lodestone/scoring/TestScoreDeterminism` produzieren
  byte-identische `recommendations.jsonl`. E2E-Diff zwischen zwei
  Läufen in `e2e/lodestone_test.sh`.

LLMs sind ab Phase 2 erlaubt, aber **nur** in `lodestone plan` und
Hilfs-Aktionen — niemals im Score-Pfad.

## Konsequenzen

- **Plus:** Code-Review, Bug-Reports und CI-Diffs sind möglich, weil
  identische Inputs identische Outputs ergeben.
- **Plus:** Audit-Trail (`decisions.log`) wird reproduzierbar
  interpretierbar.
- **Plus:** Tests sind ohne API-Keys, ohne Netzwerk und ohne LLM-Kosten
  ausführbar.
- **Minus:** Compatibility ist heute eine grobe Heuristik (gewichtete
  Jaccard). LLM-gestützte Verfeinerung muss explizit als
  Pre-Computation außerhalb des Score-Pfads erfolgen (z. B. als
  Rationale, die in `Recommendation.Rationale` einfließt).

## Alternativen

- **LLM-Compatibility.** Verworfen — bricht Determinismus, macht jeden
  Lauf API-Kosten-pflichtig.
- **Probabilistisches Scoring.** Verworfen — selbst mit seed wäre die
  Output-Konsistenz schwerer zu garantieren als mit reiner Arithmetik.

## Quelle

[Phase-1-Design](../../superpowers/specs/2026-05-20-lodestone-design.md),
[Determinismus](../determinism.md), [Scoring](../scoring.md),
`internal/lodestone/scoring/`.
