# Determinismus

Zwei aufeinanderfolgende `lodestone score`-Läufe mit identischem
Fingerprint und identischer Signal-Liste produzieren **byte-identische**
`recommendations.jsonl`. Diese Eigenschaft ist nicht zufällig — sie
ist Kern-Vertrag des Phase-1-MVP und wird strukturell verifiziert.

## Garantie

- **Score-Pfad LLM-frei.** Compatibility, Effort und Risk werden rein
  deterministisch aus dem Fingerprint und den Signal-Tags berechnet.
  Phase 2+ darf LLMs nur in `lodestone plan` und Hilfs-Aktionen
  einsetzen — niemals im Score-Pfad.
- **Stabile Sortierung.** Schlüssel: `compatibility DESC, stars DESC,
  id ASC`. Vollständig deterministisch bei Tie-Breakern.
- **Reproduzierbare IDs.** Recommendation-ID =
  `sha256:hex(signal_id + "|" + json(fingerprint))`. Identische Inputs
  produzieren identische IDs.
- **Atomare Datei-Operationen.** Score-Output wird über `tmp + rename`
  geschrieben — kein partieller Zustand bei Crash.

## Verifikation

| Ebene | Test |
|---|---|
| Unit | `internal/lodestone/scoring/scoring_test.go::TestScoreDeterminism` — drei aufeinanderfolgende Score-Läufe mit `json.Marshal`-Byte-Vergleich. |
| End-to-End | `e2e/lodestone_test.sh` — Snapshot vor zweitem `score`, `diff -q` gegen `recommendations.jsonl` nach zweitem Lauf. |
| CI | Job `e2e` in `.github/workflows/ci.yml` führt den E2E-Test bei jeder PR aus. |

## Konsequenzen

- Neue Score-Heuristiken müssen den Determinismus-Test grün halten.
- Neue Sources müssen ihre Outputs stabil sortieren; flatternde
  Reihenfolge bricht die JSONL-Hashes.
- Caching ist tagesgenau pro Source, nicht stundengenau — sonst würde
  derselbe Befehl an unterschiedlichen Uhrzeiten unterschiedliche
  Caches sehen.

Hintergrund:
[ADR-0006 — Deterministisches Scoring](adr/0006-deterministisches-scoring.md).
