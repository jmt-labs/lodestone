# Roadmap

> **Single Source of Truth** für den Phasen-Status. `README.md §Status`
> und `CHANGELOG.md` zitieren denselben Stand. Konsistenz wird durch
> `make docs-status-check` strukturell abgesichert.

Stand: 2026-05-20. Phasen 1–4 sind auf `main` gemerged und CI-grün.

## Phasen-Übersicht

| Phase | Inhalt | Status | Tag | Spec | Plan |
|---|---|---|---|---|---|
| 1 | Deterministische Pipeline (`fingerprint` → `ingest` → `score` → `signals`) | ✅ done | `v0.1.0-alpha` (geplant) | [Phase-1-Design](../superpowers/specs/2026-05-20-lodestone-design.md) | [Phase-1-Plan](../superpowers/plans/2026-05-20-lodestone-mvp.md) |
| 2 | Planning + Onboarding (`init`, `plan`, vier neue Quellen, Audit-Log) | ✅ done | — | [Phase-2-Design](../superpowers/specs/2026-05-20-lodestone-phase2-design.md) | [Phase-2-Plan](../superpowers/plans/2026-05-20-lodestone-phase2.md) |
| 3 | Remote-Schnittstellen (`lodestone-mcp`, `memory`, GitHub-Action) | ✅ done | — | [Phase-3-Design](../superpowers/specs/2026-05-20-lodestone-phase3-design.md) | — |
| 4 | Auto-PR-Engine (`apply`, `undo`, `stats`, Safety-Gates) | ✅ done | — | [Phase-4-Design](../superpowers/specs/2026-05-20-lodestone-phase4-design.md) | — |
| 5+ | `recommend`, `calibrate`, `share` (Cross-Repo-Sharing) | 🚧 Stub / geplant | — | [Privacy-Spec](../superpowers/specs/2026-05-20-lodestone-sharing-privacy.md) | — |

## Phase 1 — Deterministische Pipeline ✅

LLM-freie, byte-reproduzierbare Pipeline für AI-Ökosystem-Signale.

**Lieferumfang:**
- Schemas für `Signal`, `Fingerprint`, `Recommendation`, `WorkPackage`.
- File-basierter Store (JSONL + atomic-rename).
- Sources: `github_trending`, `hackernews`.
- Fingerprint für Go und Node (Walker + Framework-Heuristik).
- Scoring: `compatibility` (gewichtete Jaccard), `effort`, `risk`.
- Config-Loader für `.lodestone.yaml`.
- Subkommandos: `fingerprint`, `ingest`, `score`, `signals`.
- E2E-Test `e2e/lodestone_test.sh` mit Determinismus-Diff.

## Phase 2 — Planning + Onboarding ✅

LLM-Integration über die Claude-CLI, ohne neue Go-Dependencies.

**Lieferumfang:**
- Vier neue Sources: `arxiv`, `anthropic_changelog`, `openai_changelog`, `npm_trending`.
- Shared Cache- und Retry-Helper.
- Planning-Engine: `BuildPrompt`, `SplitResponse`, `Persist` (shell-out an `claude --print`).
- Subkommandos: `lodestone init`, `lodestone plan <rec-id>` (mit `--dry-run`, `--model`).
- Audit-Log: `.lodestone/decisions.log` (JSONL, append-only).
- Vier Claude-Skills unter `flavors/lodestone/skills/`, via `go:embed` eingebettet.

## Phase 3 — Remote-Schnittstellen ✅

Lodestone wird per MCP-Server und GitHub-Action ansprechbar.

**Lieferumfang:**
- Zweites Binary `lodestone-mcp` (JSON-RPC 2.0, Protocol-Version `2024-11-05`, stdio-Transport).
- Fünf MCP-Tools: `list_signals`, `query_trends`, `score_repo`, `generate_plan`, `record_decision`.
- Memory-Konsolidierung: `lodestone memory` → `.claude/memory.json`.
- GitHub-Action-Template `.github/workflows/templates/lodestone-weekly.yml` (Sonntag 03:00 UTC, opt-in).
- `make build` und `.goreleaser.yaml` bauen beide Binaries.

## Phase 4 — Auto-PR-Engine ✅

Eine Recommendation wird zum Draft-PR — mit harten Schranken.

**Lieferumfang:**
- `lodestone apply <rec-id>`: vier Safety-Gates (`risk == low`, `effort == XS`,
  `compatibility >= 0.85`, kein Apply in letzten 24 h) plus sauberes `git status`.
- `lodestone undo <branch-or-rec-id>`: PR schließen + Branch entfernen.
- `lodestone stats`: aggregierte Apply-States.
- Pluggable `GitRunner` / `PRRunner` — Real-Implementierungen plus Fakes für Tests.
- Branch-Schema: `lodestone/apply-<rec-suffix>-<date>`, **nie auf `main`**, immer Draft.
- Privacy-Spec für Cross-Repo-Sharing (Phase 5+).

## Phase 5+ — Geplant 🚧

Stubs für `recommend`, `calibrate`, `share`. Privacy-Spec
[liegt bereits](../superpowers/specs/2026-05-20-lodestone-sharing-privacy.md):
k=5-Anonymität für Goals/TechInterests, Opt-In-Flow,
Re-Identifikations-Schutz. Implementierung beginnt erst nach
Beantwortung der offenen Privacy-Fragen.

## Wie der Status aktualisiert wird

Diese Datei ist die kanonische Quelle. Jede Phase-bezogene PR muss
laut [`contributor/docs-maintenance.md`](../contributor/docs-maintenance.md):

1. den Status in dieser Tabelle aktualisieren,
2. den entsprechenden Block in `README.md §Status` anpassen,
3. `CHANGELOG.md` erweitern.

`make docs-status-check` verifiziert Konsistenz zwischen den drei
Stellen vor dem Merge.
