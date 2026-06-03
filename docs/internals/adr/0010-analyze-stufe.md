# ADR-0010 — `analyze`-Stufe: LLM-Anreicherung nach dem Score

## Status

Proposed, 2026-06-03. Teil des Capability-Gap-Epics. Erweitert
[ADR-0006](0006-deterministisches-scoring.md).

## Kontext

Für gezielte Verbesserungspläne muss lodestone semantisch verstehen,
welche Capability-Lücke ein Trend-Signal im aufrufenden Repo füllt und wie
der Migrationspfad aussieht. Das erfordert LLM-Reasoning. Gleichzeitig darf
der Score-Pfad nach ADR-0006 byte-deterministisch und LLM-frei bleiben.

## Entscheidung

Eine neue, **optionale** Pipeline-Stufe `analyze` läuft **nach** `score`:

```
fingerprint → ingest → score → analyze (LLM) → plan → apply
```

- `analyze` reichert die Top-N-Recommendations an und schreibt das Ergebnis
  in einen **separaten** Output `.lodestone/enrichments.jsonl`.
  `recommendations.jsonl` bleibt unangetastet — die Determinismus-Grenze
  ist auch auf Datei-Ebene sichtbar.
- Das Ergebnis ist ein **eigener Typ** `schema.Enrichment`
  (`CapabilityGap`, `CurrentApproach`, `AffectedFiles`, `MigrationPath`,
  `CounterEvidence`, `Confidence`) — kein Vermischen von deterministischem
  und LLM-Output im selben Record.
- LLM-Zugriff über denselben Shellout zu `claude --print` wie `plan`
  ([ADR-0003](0003-claude-cli-shellout.md)). Das Runner-Pattern wird nach
  `internal/lodestone/llm/` extrahiert und von `analyze` und `planning`
  geteilt.
- Die Antwort kommt als JSON-Block zwischen `===ENRICHMENT===`/`===END===`
  und wird per `json.Unmarshal` deserialisiert — robust und testbar.
- `analyze` ist optional. Ohne diese Stufe verhält sich lodestone exakt wie
  heute; `plan` nutzt ein Enrichment nur, wenn es existiert.

ADR-0006 wird damit explizit erweitert: LLMs sind jetzt in `plan` **und**
`analyze` erlaubt — weiterhin **nie** im Score-Pfad.

## Konsequenzen

- **Plus:** Semantische Anreicherung ohne Bruch der Score-Determinismus-
  Garantie; die in ADR-0006 §Konsequenzen genannte „Pre-Computation
  außerhalb des Score-Pfads" wird hier konkret eingelöst.
- **Plus:** `analyze` über `FakeRunner` vollständig ohne API-Keys, Netzwerk
  oder Kosten testbar (wie `planning_test.go`).
- **Plus:** Getrennte Persistenz hält Audits und Diffs sauber.
- **Minus:** Eine zusätzliche Stufe und ein zusätzlicher Output-Artefakt.
  Gerechtfertigt durch den Kernnutzen (gezielte Pläne).
- **Minus:** LLM kann ungültiges JSON liefern — expliziter Fehlerpfad statt
  stillem Fallback.

## Alternativen

- **Die leeren Felder `Rationale`/`CounterEvidence`/`SuggestedNext` in
  `Recommendation` füllen.** Verworfen — zu dünn für strukturierte
  Gap-Analyse und vermischt LLM-Output mit dem deterministischen Record.
- **LLM-Re-Ranking der Recommendations.** Verworfen — bräche die
  byte-stabile Sortierung und damit ADR-0006.
- **LLM direkt im Score-Pfad.** Verworfen — siehe ADR-0006.

## Quelle

[Capability-Gap-Design](../../superpowers/specs/2026-06-03-capability-gap-design.md),
[ADR-0003](0003-claude-cli-shellout.md),
[ADR-0006](0006-deterministisches-scoring.md),
`internal/lodestone/planning/`.
