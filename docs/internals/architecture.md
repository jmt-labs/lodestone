# Architektur

> Stand: 2026-05-20, Phasen 1–4 auf `main`.

Dieses Dokument beschreibt die Code-Struktur und den Datenfluss. Für
die User-Sicht siehe [Quickstart](../user/quickstart.md), für die
Subkommandos die [Befehls-Übersicht](../user/commands/README.md).
Verwandte Detail-Dokus:
[Datenmodell](data-model.md) ·
[Scoring](scoring.md) ·
[Determinismus](determinism.md) ·
[Artefakte](artifacts.md) ·
[Roadmap](roadmap.md) ·
[ADRs](adr/README.md).

## Drei-Ebenen-Modell

```
┌──────────────────────────────────────────────────────────────────┐
│ 1) CLI lodestone (Cobra)                                         │
│    fingerprint · ingest · score · signals · plan · apply · …     │
├──────────────────────────────────────────────────────────────────┤
│ 2) Claude-Skills (flavors/lodestone/skills/, embedded)           │
│    scout · recommend · plan · review-trends                      │
├──────────────────────────────────────────────────────────────────┤
│ 3) MCP-Server lodestone-mcp (stdio JSON-RPC 2.0)                 │
│    list_signals · query_trends · score_repo · generate_plan ·    │
│    record_decision                                               │
└──────────────────────────────────────────────────────────────────┘
        │                                  │
        ▼                                  ▼
   .lodestone/ (lokales Artefakt)    Claude / IDE / Editor
```

Alle drei Ebenen rufen denselben Kern unter `internal/lodestone/*` auf.
Kein Duplikat-Code zwischen CLI, Skill und MCP — der MCP-Server ist
eine dünne Adapter-Schicht.

## Code-Layout

```
cmd/lodestone/                    ← CLI (Cobra)
  main.go            common.go
  fingerprint.go     ingest.go        score.go     signals.go
  init.go            plan.go          memory.go
  apply.go           (apply, undo, stats in einer Datei)

cmd/lodestone-mcp/main.go         ← MCP-Server (stdio JSON-RPC)

internal/lodestone/
  schema/      ← Signal, Fingerprint, Recommendation, WorkPackage
  store/       ← FileStore über .lodestone/*.jsonl
  ingest/      ← Source-Interface + 6 Adapter + cache/retry-Helper
  fingerprint/ ← Walker für Go + Node, Framework-Heuristik
  scoring/    ← compatibility, effort, risk (deterministisch)
  planning/    ← Runner-Interface + ClaudeRunner + Prompt-Template
  audit/       ← decisions.log (JSONL, append-only)
  memory/      ← decisions.log → .claude/memory.json
  mcp/         ← Protocol + ToolRegistry + 5 Built-in-Tools
  apply/       ← Safety-Gates + State + Apply-Engine + Git/PR-Runner
  skills/      ← go:embed der vier Skill-Markdown-Files

internal/config/                   ← .lodestone.yaml Loader (yaml.v3)

flavors/lodestone/skills/          ← Kanonische Skill-Markdowns

base/models.yaml                   ← Modell-Routing
                                     (planning → claude-opus-4-7 etc.)
```

## Datenfluss

```
┌─────────┐   ingest    ┌────────────────┐
│ Sources ├────────────►│ signals.jsonl  │
└─────────┘             └────────┬───────┘
                                 │
┌────────────┐  fingerprint  ┌───┴──────────────┐
│ Repo-Files ├──────────────►│ fingerprint.json │
└────────────┘               └──────┬───────────┘
                                    │
                                    ▼
                     ┌──────────────────────────┐
                     │ score (deterministisch)  │
                     └────────────┬─────────────┘
                                  │
                                  ▼
                     ┌──────────────────────────┐
                     │ recommendations.jsonl    │
                     │ (compat DESC, stars      │
                     │  DESC, id ASC)           │
                     └────┬──────────────┬──────┘
                          │              │
                       plan│              │apply (Phase 4)
                          ▼              ▼
              ┌─────────────────┐   ┌───────────────────┐
              │ Claude-CLI ruft │   │ Safety-Gates →    │
              │ planning.Engine │   │ Branch+Commit+    │
              └────────┬────────┘   │ Push+Draft-PR     │
                       │            └─────────┬─────────┘
                       ▼                      │
       docs/superpowers/{specs,plans}/        ▼
                                    .lodestone/applies.jsonl
```

Jede Aktion (`fingerprint`, `ingest`, `score`, `plan`, `apply`) schreibt
zudem einen Eintrag nach `.lodestone/decisions.log`. Diese Audit-Spur
wird via `lodestone memory` periodisch nach `.claude/memory.json`
konsolidiert.

## Lokale Artefakte und Determinismus

Layout und Schreib-Strategie der `.lodestone/`-Dateien:
[Artefakte](artifacts.md). Die Score-Pipeline garantiert
byte-identische Outputs bei identischen Inputs — siehe
[Determinismus](determinism.md). Test-Pyramide und CI-Gates sind in
[`contributor/testing.md`](../contributor/testing.md) beschrieben.

## Dependency-Budget

- `github.com/spf13/cobra` (CLI)
- `gopkg.in/yaml.v3` (Konfig)
- Sonst nur Go-Standardbibliothek (`encoding/{json,xml}`, `net/http`,
  `regexp`, `crypto/sha256`, `os/exec` für Shell-out).

Neue externe Deps brauchen eine Spec-Diskussion. Phase-1-Invariante
("nur Cobra + yaml.v3") wurde ab Phase 2 gelockert, aber der Geist
bleibt: vor neuen Deps fragen, ob die Standardbibliothek reicht.

## Erweiterungspunkte

| Wunsch | Erweiterungspunkt |
|---|---|
| Neue Source | `internal/lodestone/ingest/<name>.go` implementiert das `Source`-Interface; in `cmd/lodestone/ingest.go::buildSource` registrieren |
| Neuer Fingerprint-Detektor | `internal/lodestone/fingerprint/<language>.go` analog zu `golang.go`/`node.go` |
| Neues MCP-Tool | `internal/lodestone/mcp/tools.go::RegisterBuiltins` erweitern |
| Anderes LLM | `planning.Runner`-Interface implementieren (analog zu `ClaudeRunner` / `FakeRunner`) |
| Andere Konfig-Quelle | `internal/config/config.go` um Loader-Variante ergänzen |

Tests müssen mit. CI-Gate ist `make test lint vuln e2e` — alle vier
grün, sonst kein Merge.
