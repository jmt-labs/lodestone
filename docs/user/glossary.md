# Glossar

Zentrales Begriffsverzeichnis. Alphabetisch sortiert.

### Apply

Die Aktion, eine Recommendation als Draft-PR ins Repository
einzuspielen. Ausgeführt durch `lodestone apply <rec-id>`. Vier
Safety-Gates müssen passen: `risk == low`, `effort == XS`,
`compatibility >= 0.85`, kein Apply in den letzten 24 h. Zusätzlich
muss `git status` sauber sein. Siehe [Apply-Befehl](commands/apply.md)
und [Safety-Gates-ADR](../internals/adr/0008-apply-safety-gates.md).

### Compatibility

Skalarer Score zwischen `0.0` und `1.0`, der angibt, wie gut ein
Signal zum Repository-Fingerprint passt. Berechnet als gewichtete
Jaccard-Ähnlichkeit von Signal-Tags mit Repo-Frameworks/-Languages
(Language-Match 1.5×, Framework-Match 1.0×). Formel:
[Scoring-Algorithmus](../internals/scoring.md).

### Decision-Log

Append-only JSONL-Datei unter `.lodestone/decisions.log`, in die jedes
Verb (`fingerprint`, `ingest`, `score`, `plan`, `apply`, …) einen
Eintrag mit Timestamp, Argumenten und Outcome schreibt. Wird via
`lodestone memory` periodisch nach `.claude/memory.json` konsolidiert.

### Effort

Kategorischer Score `XS` / `S` / `M` / `L` / `XL`, der den geschätzten
Umsetzungsaufwand einer Recommendation beschreibt. Default `M`; `XL`
bei 0 Match; `S` bei Match und Stars < 100. `XS` ist Voraussetzung für
einen Auto-Apply.

### Fingerprint

Strukturiertes Profil des analysierten Repositories — Sprachen,
Frameworks, Dependencies, LOC pro Sprache, Test-Ratio, CI-Provider,
MCP-Konfiguration, plus Goals und TechInterests aus `.lodestone.yaml`.
Erzeugt durch `lodestone fingerprint`, gespeichert in
`.lodestone/fingerprint.json`. Schema:
[Datenmodell](../internals/data-model.md).

### MCP-Tool

Eine Funktion, die der MCP-Server `lodestone-mcp` über JSON-RPC 2.0
anbietet. Aktuell fünf Built-ins: `list_signals`, `query_trends`,
`score_repo`, `generate_plan`, `record_decision`. Aufrufbar aus
Claude Desktop, Claude Code und anderen MCP-Clients. Siehe
[MCP-Server-Setup](mcp-server.md).

### Phase 1

Die initiale Pipeline ohne LLM-Calls: `fingerprint → ingest → score
→ signals`. Vollständig deterministisch, nur GitHub-Trending und
HackerNews als Quellen. Abgeschlossen, Tag-Ziel `v0.1.0-alpha`.

### Phase 2

LLM-Integration über die Claude-CLI: `lodestone init` für Onboarding,
`lodestone plan` für Spec/Plan-Generierung. Vier zusätzliche
Ingest-Quellen (arXiv, Anthropic-Changelog, OpenAI-Changelog,
npm-Trending), Audit-Log. Abgeschlossen.

### Phase 3

Remote-Schnittstellen: zweites Binary `lodestone-mcp` (MCP-Server
über stdio), `lodestone memory` zur Decision-Log-Konsolidierung,
GitHub-Action-Template für wöchentliche Scans. Abgeschlossen.

### Phase 4

Auto-PR-Engine: `lodestone apply` mit vier Safety-Gates, `lodestone
undo` zum Zurückrollen, `lodestone stats` für Erfolgsstatistiken.
Niemals direkt auf `main`, immer Draft, maximal ein Apply pro 24 h.
Abgeschlossen.

### Phase 5+

Geplant, noch nicht implementiert: `lodestone recommend`
(interaktive Empfehlungs-Loop als Skill), `lodestone calibrate`
(Scoring-Gewichte aus Decision-Log nachjustieren), `lodestone share`
(anonymes Cross-Repo-Sharing — Privacy-Spec liegt bereits unter
[`docs/superpowers/specs/2026-05-20-lodestone-sharing-privacy.md`](../superpowers/specs/2026-05-20-lodestone-sharing-privacy.md)).

### Recommendation

Eine deterministisch sortierte Empfehlung — Resultat von `lodestone
score`. Verbindet ein Signal mit dem Fingerprint und liefert
`compatibility`, `effort`, `risk` und eine Rationale. ID-Schema:
`sha256:hex(signal_id + "|" + json(fingerprint))`. Gespeichert in
`.lodestone/recommendations.jsonl`, sortiert nach
`compatibility DESC, stars DESC, id ASC`.

### Risk

Kategorischer Score `low` / `med` / `high`. `low` bei Stars ≥ 500
∧ LastCommit < 90 d ∧ License vorhanden. `high` bei fehlender
License oder LastCommit > 180 d (stale). Sonst `med`. `low` ist
Voraussetzung für Auto-Apply.

### Safety-Gate

Eine der vier Bedingungen, die alle gleichzeitig erfüllt sein müssen,
damit `lodestone apply` einen Draft-PR öffnet. Siehe
[Apply-Befehl](commands/apply.md). Zusätzlich gilt: `git status`
muss sauber sein.

### Signal

Ein atomares Trend-Ereignis aus einer externen Quelle — etwa ein
GitHub-Trending-Repo, eine HackerNews-Story oder ein arXiv-Paper.
Gespeichert als eine Zeile in `.lodestone/signals.jsonl`. Felder:
ID, Name, URL, Sprache, Tags, Stars, Release-Datum, Source.
Schema: [Datenmodell](../internals/data-model.md).

### Source

Ein Ingest-Adapter, der Signale aus einer externen API holt. Phase 1+2
liefern sechs Sources: `github_trending`, `hackernews`, `arxiv`,
`anthropic_changelog`, `openai_changelog`, `npm_trending`. Alle teilen
Cache- und Retry-Helper.

### WorkPackage

Strukturierter Output der Planning-Engine — eine umgesetzte
Recommendation als Spec-Plus-Plan-Paar im superpowers-Format
(`docs/superpowers/specs/…` und `…/plans/…`). Generiert durch
`lodestone plan <rec-id>` über die Claude-CLI. Schema:
[Datenmodell](../internals/data-model.md).
