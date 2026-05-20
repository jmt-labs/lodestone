# Lodestone Cross-Repo-Sharing — Privacy-Spec

**Status:** Vor-Implementierungs-Diskussion. Phase 4 implementiert das
**nicht**. Diese Spec hält die Anforderungen fest, bevor Code entsteht.

## Problem

Lodestone-Decisions aus mehreren Repos zusammengeführt könnten
nützliche Erkenntnisse liefern (welche Recommendations setzen sich
durch? Welche Quellen liefern Signal-Müll?). Die Daten sind aber
sensibel: Repo-Pfade, Goal-Felder, Dependency-Listen.

## Felder pro Datenklasse

### `Signal`

Veröffentlichbar (gefahrlos):
- `Source`, `URL`, `Title`, `License`, `Language`, `Stars`, `TopicTags`
- `CapturedAt`, `LastCommit`

Niemals teilen:
- nichts — `Signal` enthält nur öffentlich abrufbare Daten

### `Fingerprint`

Veröffentlichbar nach Anonymisierung:
- `Languages`, `Frameworks` (k-Anonymität ab k=5: nur wenn dieselbe
  Kombination mindestens in 5 Repos auftritt)
- `HasCI`, `CIProvider`

Niemals teilen ohne explizite Opt-In:
- `Goals`, `TechInterests` (frei-text, kann Geschäftsstrategie verraten)
- `Deps` (komplette Version-Map)
- `LOCPerLanguage` (Größenordnung der Codebase)
- `MCPServers` (interne Toolchain)

### `Recommendation`

Veröffentlichbar:
- `Compatibility`, `Effort`, `ROI`, `Risk` (aggregierte Statistik)
- `SignalID` (referenziert öffentliches Signal)

Niemals teilen:
- `Rationale`, `CounterEvidence` (können Repo-Kontext leaken)

### `Decision`-Log

Veröffentlichbar nach Opt-In:
- Anzahl `verb`-Aufrufe je Tag

Niemals teilen:
- `Args` (frei-text)
- `Detail` (frei-text)
- `Timestamp` exakt — auf Tages-Granularität truncaten

## Mechanik

1. **Opt-In explizit pro Repo**: `.lodestone.yaml` enthält
   `sharing.enabled: true` mit Datum des Consent.
2. **Geteilte Daten**: nur die "veröffentlichbar"-Felder, in einer
   `share-YYYY-MM-DD.json` exportiert und in `.lodestone/share/`
   abgelegt. Der User entscheidet, **wie** er diese Datei weitergibt.
3. **Re-Identifikations-Schutz**:
   - Repo-Identifier (Pfad, GitHub-Slug) niemals teilen.
   - Goals/TechInterests nur mit k=5-Anonymität teilen, sonst weg.
   - Stars/LOC auf Bucket-Granularität rounden (z. B. 10, 100, 1000,
     10000).
4. **Revoke**: `sharing.enabled: false` plus expliziter Löschauftrag
   in einem (noch zu schaffenden) Aggregator-Repo.

## Aggregator (Phase 4+)

Nicht Teil dieser Spec. Wenn ein Aggregator-Service entsteht, braucht
er eine eigene Spec mit:

- Storage-Policy (Verschlüsselung at-rest, kein PII)
- Retention (z. B. 6 Monate)
- Zugriffs-Modell (öffentliche read-only Dashboards?)
- Lösch-Garantien (Right-to-Erasure innerhalb von 30 Tagen)

## Implementierungs-Hinweise

- `lodestone share` (Phase 4+, derzeit Stub) wird die obigen Regeln
  enforced.
- `internal/lodestone/share/` Package hält die Anonymisierungs-Logik
  zentral und ist 100 % testbar (reine Funktionen über Schema-Typen).
- Vor jedem Release der Sharing-Funktion: Security-Audit der
  Anonymisierungs-Funktion mit echten Repo-Beispielen.

## Offene Fragen

1. Wie wird Consent UI-mäßig abgefragt? (CLI-Prompt? Skill?)
2. Wer betreibt den Aggregator? Eigenes Org-Repo oder freier
   Marketplace?
3. Welche Bucket-Granularität ist „grob genug" ohne Daten unbrauchbar
   zu machen?

Diese Fragen müssen vor Implementierungs-Start beantwortet sein.
