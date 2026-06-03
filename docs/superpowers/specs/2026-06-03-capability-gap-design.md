# Lodestone Capability-Gap — Epic-Design

**Datum:** 2026-06-03
**Voraussetzung:** Phasen 1–4 auf `main`.
**Status:** Forward-looking Epic-Spec. Dieser Durchgang liefert nur die
Planungsartefakte (Spec, vier Phasen-Pläne, ADRs); der Produktionscode
folgt phasenweise über eigene Feature-Branches und PRs.

## Motivation

Lodestone liest heute das AI-Ökosystem und gleicht Signale gegen das
aufrufende Repo ab — aber das Matching ist **rein syntaktisch**.
`compatibility` (`internal/lodestone/scoring/compatibility.go`) ist ein
gewichteter Jaccard-Overlap von Sprachen und Frameworks. Der Fingerprint
(`internal/lodestone/schema/fingerprint.go`) erfasst Struktur (Sprachen,
Deps, LOC, CI, MCP-Server), aber **nicht, was ein Repo KI-mäßig tut**:
welche Provider und Modelle es nutzt, ob es Tool-Use, RAG, Prompt-Caching
oder Eval-Harnesses einsetzt.

Dadurch bleibt das Kernversprechen — *brandaktuelle KI-Trends mit dem Repo
abgleichen und gezielte, detaillierte Verbesserungspläne erzeugen* — nur
oberflächlich eingelöst. Eine Recommendation trägt heute nur
`compatibility`, `effort`, `risk`; die Felder `Rationale`,
`CounterEvidence` und `SuggestedNext` sind in der Praxis leer. Die
Plan-Generierung (`internal/lodestone/planning/`) hat damit kaum Material
für „detailliert".

## Ziel

Eine **Capability-Gap-Analyse** einführen: Lodestone erfasst die IST-KI-
Nutzung des Repos, gleicht sie gegen reife Trend-Signale ab und liefert
konkrete Lücken plus Migrationspfade, die in die Plan-Generierung fließen.

**Anker-Use-Case dieses Epics:** „Neue Tools adoptieren", gap-getrieben.
Also nicht „dieses Repo trendet", sondern: *„Dein Repo löst X heute mit
Ansatz A; dieses reife Tool/Library füllt Lücke Y besser — hier der
konkrete Integrationspfad."* Generisches Modell-/API-Upgrade
(veraltete Modell-IDs, fehlendes Caching) ist bewusst **außerhalb** dieses
Epics und kann ein späteres aufsetzen.

## Leitprinzip: Determinismus bleibt unangetastet

Der Score-Pfad bleibt LLM-frei und byte-deterministisch (ADR-0006). Alle
Semantik landet in einer neuen, klar getrennten Stufe `analyze`, die
*nach* `score` läuft. Die Pipeline wird:

```
fingerprint → ingest → score → analyze (NEU, LLM, optional) → plan → apply
```

`analyze` ist optional. Wer es überspringt, bekommt exakt das heutige
Verhalten. Die Persistenz ist getrennt: `analyze` schreibt nach
`.lodestone/enrichments.jsonl`, `recommendations.jsonl` bleibt
unangetastet. Diese Trennung macht die Determinismus-Grenze auch auf
Datei-Ebene sichtbar.

Die Entscheidung ist in [ADR-0010](../../internals/adr/0010-analyze-stufe.md)
festgehalten; der semantische Fingerprint in
[ADR-0009](../../internals/adr/0009-semantischer-ai-fingerprint.md).

## Architektur

### A) Semantischer AI-Fingerprint (deterministisch, kein LLM)

Neues optionales Sub-Struct in `schema/fingerprint.go`,
`FingerprintSchemaVersion` 1 → 2. Alle Felder `omitempty`, das Sub-Struct
selbst `nil`-bar — bei Nicht-AI-Repos bleibt es leer und erzeugt kein
Schema-Rauschen:

```go
type AICapabilities struct {
    Providers   []string // anthropic, openai, google, ollama, local
    SDKs        []string // anthropic-sdk-go, @anthropic-ai/sdk, openai, ai (Vercel)
    ModelIDs    []string // Regex aus Quelltext: claude-…, gpt-…, gemini-…
    Patterns    []string // tool_use, rag, streaming, prompt_caching, eval, embeddings, agent_loop
    MCPRole     string   // server | client | both | ""
    EvalHarness []string // promptfoo, langsmith, deepeval, custom
}
```

Detektion **rein per Code-Heuristik** — der Fingerprint erhebt Fakten,
fällt keine Urteile (das Interpretieren übernimmt `analyze`):

| Feld | Heuristik | Wiederverwendung |
|---|---|---|
| `Providers`/`SDKs` | Dep-Lookup-Maps | `goAISDKs` in `golang.go`, `nodeAISDKs` in `node.go` (analog `goFrameworks`/`nodeFrameworks`) |
| `EvalHarness` | Dep-Lookup + Datei-Heuristik (`evals/`, `*.eval.*`) | Dep-Map + bestehender `walk()` |
| `MCPRole` | `.mcp.json` → server; `@modelcontextprotocol/sdk`/`mcp-go` in Deps → client; `cmd/*-mcp/` → server | erweitert `detectMCPServers` |
| `ModelIDs` | Regex `claude-[a-z0-9-]+`, `gpt-[0-9o.-]+`, `gemini-[a-z0-9.-]+` über Quelltext | neue `aidetect.go`, in `walk()` eingehängt |
| `Patterns` | Symbol-Heuristik (`tool_use`/`tools=` → tool_use; `embeddings`/`pgvector`/`chroma` → rag; `.stream(` → streaming; `cache_control` → prompt_caching) | `aidetect.go` |

`ModelIDs`/`Patterns` hängen in den bestehenden `walk()` ein **ohne
Doppel-I/O**: der Dateiinhalt wird bereits für `countLines` gelesen und
durchgereicht. `Analyze()` setzt `fp.AI` nach `detectMCPServers()`; ohne
AI-Signale bleibt es `nil`.

**Was LLM bräuchte und deshalb NICHT hier liegt:** Intent, Reifegrad und
Qualität des KI-Ansatzes. Das ist Aufgabe von `analyze`.

### B) Die `analyze`-Stufe (LLM, separat)

Neues Paket `internal/lodestone/analyze/`, strukturell ein Klon des
`planning`-Patterns (`analyze.go`, `prompt.go`, `analyze_test.go`).

Das Runner-Pattern aus `planning/runner.go` (`Runner`, `ClaudeRunner`,
`FakeRunner`) wird in ein gemeinsames Mini-Paket
`internal/lodestone/llm/` extrahiert und von `analyze` und `planning`
geteilt — das vermeidet eine `analyze` → `planning`-Abhängigkeit.
`planning` re-exportiert die Typen für Rückwärtskompatibilität. Der
Shellout zu `claude --print` bleibt unverändert (ADR-0003).

**Output: neuer Typ** `schema/enrichment.go`. Die schmalen optionalen
Felder in `recommendation.go` sind zu dünn für strukturierte Gap-Analyse
und gehören dem deterministischen Pfad — daher ein separater Record:

```go
type Enrichment struct {
    SchemaVersion    int
    RecommendationID string
    SignalID         string
    GeneratedAt      time.Time
    Model            string
    CapabilityGap    string   // Lücke Y, die das Tool füllt
    CurrentApproach  string   // Ansatz A im Repo heute
    AffectedFiles    []string // konkrete Code-Stellen
    MigrationPath    []string // geordnete atomare Schritte
    CounterEvidence  string   // wann sich die Migration NICHT lohnt
    Confidence       string   // low/med/high
}
```

Persistiert in `.lodestone/enrichments.jsonl` über einen neuen
`EnrichmentStore` (`store/store.go` + `filestore.go`), JSONL +
atomic-rename wie die bestehenden Stores (ADR-0002).

**Prompt:** deutsches Template analog `planning/prompt.go`, JSON-
eingebettete Inputs (`fp` inkl. `AI`, `rec`, `sig`). Die Antwort kommt als
**JSON-Block** zwischen `===ENRICHMENT===` und `===END===` und wird per
`json.Unmarshal` direkt in `schema.Enrichment` deserialisiert — robuster
und testbarer als reine Textmarker. Der Prompt fragt explizit nach:
(1) welche KI-Capability-Lücke das Tool füllt, (2) wie das Repo das heute
löst (Ansatz A mit konkreten Dateien aus dem Fingerprint), (3) Migrations-
pfad in atomaren Schritten, (4) Counter-Evidence.

**`plan`-Anbindung:** `cmd/lodestone/plan.go` lädt zusätzlich das
Enrichment für die rec-id (`FindEnrichment`). `planning.BuildPrompt` wird
um einen **optionalen** Enrichment-Block erweitert (Gap, Migrationspfad,
betroffene Dateien). Fehlt das Enrichment, ist das Verhalten exakt wie
heute — keine Bruchkante.

### C) Reicherer Signal-Content (nur soweit für den Anker-Case nötig)

Changelogs sind heute nur Titel + Datum; für „dieses Tool füllt Lücke Y"
braucht `analyze` mehr Material. Minimal-invasiv, **null neue
Dependencies** (Stdlib `net/http` + bestehende `stripTags`/`cachePath`-
Helfer), konsistent mit ADR-0005:

- `schema/signal.go`: optionales `ReadmeExcerpt` / `Capabilities`,
  `SignalSchemaVersion` 1 → 2.
- `ingest/github_trending.go` (Hauptquelle für Tool-Adoption): README via
  `GET /repos/{owner}/{repo}/readme`, gekürzt + bereinigt; Caching greift
  über die bestehende `cachePath`/`saveCache`/`loadCache`-Mechanik
  automatisch.
- `ingest/changelog.go`: erste Absätze pro Entry in `Summary` (sekundär in
  diesem Epic, relevanter für ein späteres Modell-Upgrade-Szenario).

## Code-Layout

```
internal/lodestone/
  llm/                  runner.go (extrahiert aus planning)
  analyze/              analyze.go  prompt.go  analyze_test.go
  schema/               enrichment.go (neu), fingerprint.go (AICapabilities),
                        signal.go (ReadmeExcerpt/Capabilities)
  fingerprint/          aidetect.go (neu), golang.go/node.go (AI-SDK-Maps)
  store/                store.go/filestore.go (EnrichmentStore)

cmd/lodestone/
  analyze.go            (neu)
  plan.go               (lädt Enrichment)
```

## Phasen-Schnitt

Vier Phasen, jede eigenständig mergebar. **Phase 1 + 2 zusammen sind der
erste vertikale Durchstich** des Anker-Use-Case.

| Phase | Inhalt | Abhängigkeit | Plan |
|---|---|---|---|
| 1 | Semantischer AI-Fingerprint (deterministisch) | — | [P1-Plan](../plans/2026-06-03-capability-gap-p1-fingerprint.md) |
| 2 | `analyze`-Stufe + Anker-Durchstich (LLM) | Phase 1 | [P2-Plan](../plans/2026-06-03-capability-gap-p2-analyze.md) |
| 3 | `plan` konsumiert Enrichment | Phase 2 | [P3-Plan](../plans/2026-06-03-capability-gap-p3-plan-consumes.md) |
| 4 | Reicher Signal-Content + MCP/Skill-Adapter | Phase 2 | [P4-Plan](../plans/2026-06-03-capability-gap-p4-signals-adapters.md) |

## Test-Strategie

- **Determinismus-Regression:** `scoring/TestScoreDeterminism` muss in
  Phase 1 byte-identisch bleiben. `score` liest `fp.AI` nie; Altfixtures
  (`go_minimal`, `node_react`) haben kein `fp.AI`.
- **LLM aus CI heraushalten:** `analyze` und der erweiterte `plan`-Pfad
  werden ausschließlich über `FakeRunner` getestet (wie `planning_test.go`)
  — keine API-Keys, kein Netzwerk, keine Kosten (ADR-0003).
- **JSON-Robustheit:** Fehlerpfade `bad-json` und `runner-error` explizit
  testen, kein stilles Fallback.
- **Detektoren:** neue Fixtures `fingerprint/testdata/ai_go`, `ai_node`.
- **E2E:** `analyze --dry-run` ohne `claude`-Binary reproduzierbar; neue
  Fixture unter `e2e/fixtures/signals` (AI-Repo + GitHub-Tool-Signal →
  reproduzierbares Enrichment).
- **Coverage:** ≥ 70 % in `internal/` (CLAUDE.md), über Fixtures/FakeRunner
  erreichbar.

## Risiken / Trade-offs

- **Determinismus-Grenze:** `fp.AI` fließt via `canonicalFingerprint` in
  die `recommendationID` — bei AI-Repos ändern sich IDs gegenüber heute.
  Beabsichtigt (anderer Input → andere ID); Altfixtures bleiben unberührt.
  In Phase 1 verifizieren.
- **Heuristik-Präzision:** Pattern-/ModelID-Regex liefert False Positives.
  Akzeptabel, weil `analyze` den Kontext interpretiert; der Fingerprint
  liefert nur Hinweise.
- **YAGNI:** Nur der Anker-Use-Case. Kein generisches Modell-Upgrade, keine
  spekulativen Pattern-Detektoren über die belegten hinaus, README-Fetch
  nur für GitHub-Signale (nicht alle sechs Quellen).

## Ausdrücklich nicht in diesem Epic

- LLM im Score-Pfad (ADR-0006 bleibt).
- Generisches Modell-/API-Upgrade (veraltete Modell-IDs, Caching-Migration)
  als Anker — kann ein Folge-Epic werden.
- Automatischer Code-Edit über die Plan-/Spec-Dokumente hinaus.
- Neue Go-Dependencies.

## Anhang: GitHub-Epic-Issue (Vorlage)

Wird **nur nach expliziter Freigabe** angelegt. Body-Vorlage:

> **Epic: Capability-Gap-Analyse — gezielte KI-Verbesserungspläne**
>
> Lodestone von einem syntaktischen Trend-Radar zu einer Capability-Gap-
> Engine weiterentwickeln. Semantischer AI-Fingerprint + neue, optionale
> `analyze`-Stufe (LLM, nach dem deterministischen Score) erzeugen
> konkrete Lücken + Migrationspfade, die in `plan` einfließen.
>
> Design: `docs/superpowers/specs/2026-06-03-capability-gap-design.md`
>
> - [ ] Phase 1 — Semantischer AI-Fingerprint (#sub)
> - [ ] Phase 2 — `analyze`-Stufe + Anker-Durchstich (#sub)
> - [ ] Phase 3 — `plan` konsumiert Enrichment (#sub)
> - [ ] Phase 4 — Reicher Signal-Content + Adapter (#sub)
