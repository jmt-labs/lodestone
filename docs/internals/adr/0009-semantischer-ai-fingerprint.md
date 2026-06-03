# ADR-0009 — Semantischer AI-Fingerprint (deterministisch, kein LLM)

## Status

Proposed, 2026-06-03. Teil des Capability-Gap-Epics.

## Kontext

Der Fingerprint erfasst heute Struktur (Sprachen, Frameworks, Deps, LOC,
CI, MCP-Server), aber nicht, was ein Repo KI-mäßig tut. Damit bleibt das
Matching von Trend-Signalen gegen das Repo syntaktisch
(Sprach-/Framework-Overlap in `compatibility`). Für eine
Capability-Gap-Analyse braucht lodestone die IST-KI-Fähigkeiten des Repos.

## Entscheidung

Der Fingerprint wird um ein optionales Sub-Struct `AICapabilities`
erweitert (`Providers`, `SDKs`, `ModelIDs`, `Patterns`, `MCPRole`,
`EvalHarness`). `FingerprintSchemaVersion` steigt von 1 auf 2.

Die Erkennung ist **rein deterministisch per Code-Heuristik** — der
Fingerprint erhebt Fakten, fällt keine Urteile:

- `Providers`/`SDKs`/`EvalHarness` über Dep-Lookup-Maps (`goAISDKs`,
  `nodeAISDKs`) analog zu `goFrameworks`/`nodeFrameworks`.
- `MCPRole` über `.mcp.json`, MCP-SDK-Deps und `cmd/*-mcp/`.
- `ModelIDs`/`Patterns` über Regex/Symbol-Heuristik, eingehängt in den
  bestehenden `walk()` ohne Doppel-I/O.

**Keine LLM-Calls und kein Netzwerk** im Fingerprint — er bleibt eine
reproduzierbare, offline ausführbare Faktenerhebung. Das Interpretieren
(Intent, Reifegrad, Qualität) übernimmt die separate `analyze`-Stufe
([ADR-0010](0010-analyze-stufe.md)).

## Konsequenzen

- **Plus:** Capability-Gap-Analyse bekommt eine semantische Grundlage,
  ohne den deterministischen Charakter von `fingerprint` aufzugeben.
- **Plus:** Alle neuen Felder sind `omitempty` in einem `nil`-baren
  Sub-Struct → Nicht-AI-Repos erzeugen kein Schema-Rauschen, bestehende
  Fingerprints bleiben unverändert.
- **Plus:** Keine neuen Dependencies; nur Stdlib-Regex und bestehende
  Detektor-Maps.
- **Minus:** Heuristik liefert False Positives (z. B. Modell-IDs in
  Kommentaren). Akzeptiert, weil `analyze` den Kontext interpretiert.
- **Achtung Determinismus:** `fp.AI` fließt via `canonicalFingerprint` in
  die `recommendationID`. Bei AI-Repos ändern sich IDs gegenüber heute
  (beabsichtigt). `score` selbst liest `fp.AI` nie; `TestScoreDeterminism`
  bleibt byte-identisch, weil die Altfixtures kein `fp.AI` haben.

## Alternativen

- **LLM-Fingerprint.** Verworfen — bricht den offline/deterministischen
  Charakter von `fingerprint` und vermischt Faktenerhebung mit Reasoning.
- **Viele Top-Level-Felder statt Sub-Struct.** Verworfen — bläht das
  Schema bei Nicht-AI-Repos auf; `nil`-bares Sub-Struct ist YAGNI-konform.

## Quelle

[Capability-Gap-Design](../../superpowers/specs/2026-06-03-capability-gap-design.md),
[ADR-0006](0006-deterministisches-scoring.md),
`internal/lodestone/fingerprint/`, `internal/lodestone/schema/fingerprint.go`.
