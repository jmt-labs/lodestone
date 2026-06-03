# Capability-Gap Phase 1 — Semantischer AI-Fingerprint

**Spec:** [`../specs/2026-06-03-capability-gap-design.md`](../specs/2026-06-03-capability-gap-design.md) §A
**Abhängigkeit:** keine. Eigenständig mergebar.
**Charakter:** rein deterministisch, **keine LLM-Calls** (ADR-0009).

## Ziel

Der Fingerprint erfasst die KI-Fähigkeiten eines Repos per Code-Heuristik:
Provider, SDKs, Modell-IDs, Patterns, MCP-Rolle, Eval-Harnesses. Reine
Faktenerhebung ohne Netzwerk und ohne LLM.

## Tasks

- [ ] **P1-T1** — `schema/fingerprint.go`: `AICapabilities`-Sub-Struct
      (`Providers`, `SDKs`, `ModelIDs`, `Patterns`, `MCPRole`,
      `EvalHarness`), Feld `AI *AICapabilities` mit `omitempty`,
      `FingerprintSchemaVersion` 1 → 2. JSON-Roundtrip-/`omitempty`-Tests.
- [ ] **P1-T2** — `fingerprint/golang.go`: Map `goAISDKs` (z. B.
      `github.com/anthropics/anthropic-sdk-go`,
      `github.com/sashabaranov/go-openai`, `github.com/mark3labs/mcp-go`)
      → Provider/SDK-Extraktion. Tests.
- [ ] **P1-T3** — `fingerprint/node.go`: Map `nodeAISDKs`
      (`@anthropic-ai/sdk`, `openai`, `ai`, `@modelcontextprotocol/sdk`),
      analog. Tests.
- [ ] **P1-T4** — neue Datei `fingerprint/aidetect.go`: `detectModelIDs`
      (Regex `claude-…`/`gpt-…`/`gemini-…`) und `detectAIPatterns`
      (Symbol-Heuristik tool_use/rag/streaming/prompt_caching/eval/
      embeddings/agent_loop). In `walk()` eingehängt — der bereits für
      `countLines` gelesene Dateiinhalt wird durchgereicht, **kein
      Doppel-I/O**. Tests.
- [ ] **P1-T5** — `fingerprint/fingerprint.go`: `detectMCPServers` zu
      `MCPRole` (server/client/both/"") erweitern; `Analyze()` setzt
      `fp.AI` nach `detectMCPServers()`, bleibt `nil` ohne AI-Signale.
- [ ] **P1-T6** — Fixtures `fingerprint/testdata/ai_go` (anthropic-sdk-go
      + tool_use + Modell-ID) und `ai_node` (`@anthropic-ai/sdk` +
      streaming). Detektions-Tests gegen erwartete `AICapabilities`.
- [ ] **P1-T7** — Verifizieren: `scoring/TestScoreDeterminism` bleibt
      byte-identisch (Altfixtures haben `fp.AI == nil`). `make test lint
      vuln` grün.
- [ ] **P1-T8** — Doku: `docs/user/commands/fingerprint.md` und
      `docs/internals/*` um die AI-Felder ergänzen; Roadmap-Eintrag
      aktualisieren.

## Definition of Done

- `lodestone fingerprint` auf einem AI-Repo zeigt befülltes `ai`-Feld,
  auf einem Nicht-AI-Repo kein `ai`-Feld.
- Determinismus-Test grün, Coverage ≥ 70 % im `fingerprint`-Paket.

## Modus

Feature-Branch `feat/p1-t<task>-<slug>`, kurze atomare Commits, TDD
(Test vor Implementierung). PR gegen `main` mit `Updates #<epic>`.
