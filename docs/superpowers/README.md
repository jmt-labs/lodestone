# Superpowers — Specs, Pläne, Notes

Historische Designs und Implementierungspläne pro Phase. Die Dateien
sind nach Datum benannt, was beim Scrollen wenig hilft; diese Tabelle
gibt jedem Artefakt einen logischen Alias.

## Specs

| Datei | Alias | Phase | Status | Inhalt |
|---|---|---|---|---|
| [`specs/2026-05-20-lodestone-design.md`](specs/2026-05-20-lodestone-design.md) | **Phase-1-Design** | 1 | ✅ umgesetzt | MVP-Architektur: Schemas, Store, Sources, Fingerprint, Scoring |
| [`specs/2026-05-20-lodestone-phase2-design.md`](specs/2026-05-20-lodestone-phase2-design.md) | **Phase-2-Design** | 2 | ✅ umgesetzt | Planning-Engine, Skills, Audit-Log, neue Sources |
| [`specs/2026-05-20-lodestone-phase3-design.md`](specs/2026-05-20-lodestone-phase3-design.md) | **Phase-3-Design** | 3 | ✅ umgesetzt | `lodestone-mcp`, Memory, GitHub-Action |
| [`specs/2026-05-20-lodestone-phase4-design.md`](specs/2026-05-20-lodestone-phase4-design.md) | **Phase-4-Design** | 4 | ✅ umgesetzt | Auto-PR-Engine mit Safety-Gates |
| [`specs/2026-05-20-lodestone-sharing-privacy.md`](specs/2026-05-20-lodestone-sharing-privacy.md) | **Privacy-Spec** | 5+ | 🚧 noch nicht umgesetzt | k=5-Anonymität, Opt-In-Flow für Cross-Repo-Sharing |

## Pläne

| Datei | Alias | Phase | Status |
|---|---|---|---|
| [`plans/2026-05-20-lodestone-mvp.md`](plans/2026-05-20-lodestone-mvp.md) | **Phase-1-Plan** | 1 | ✅ alle Tasks erledigt |
| [`plans/2026-05-20-lodestone-phase2.md`](plans/2026-05-20-lodestone-phase2.md) | **Phase-2-Plan** | 2 | ✅ alle Tasks erledigt |

Phase 3 und Phase 4 wurden ohne dedizierten Plan umgesetzt (kleinere
Scope, direkt aus der Spec abgearbeitet).

## Notes

| Datei | Alias | Zweck |
|---|---|---|
| [`notes/2026-05-20-phase1-verification.md`](notes/2026-05-20-phase1-verification.md) | **Phase-1-Verifikation** | Hands-On-Verifikation der Phase-1-Akzeptanzkriterien |

## Wann neue Specs/Pläne anlegen?

Sobald eine neue Phase oder ein größeres Feature beginnt, das nicht
mehr in einen einzelnen PR passt. Format:
[Spec-Format](../contributor/spec-format.md),
[Plan-Format](../contributor/plan-format.md).

## Spec vs. ADR

- **Spec** (hier in `superpowers/specs/`) = forward-looking, „was bauen
  wir als Nächstes?".
- **ADR** (in [`docs/internals/adr/`](../internals/adr/README.md)) =
  backward-looking, „warum haben wir uns für X entschieden?".

Beide Artefakte verlinken bei Bedarf aufeinander.
