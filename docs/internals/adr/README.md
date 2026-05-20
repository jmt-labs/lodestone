# Architecture Decision Records

> Rückwirkend extrahierte Design-Entscheidungen. Specs sind
> forward-looking („was kommt in Phase X?"), ADRs sind backward-looking
> („warum haben wir uns für Y entschieden?").

Format: MADR-light — Status, Kontext, Entscheidung, Konsequenzen,
Alternativen, Quelle. Neue ADRs werden mit der nächsten freien
Nummer angelegt; alte ADRs werden nicht umgeschrieben, sondern bei
Bedarf von einem neuen ADR „Supersedes ADR-NNNN" abgelöst.

## Index

| Nr. | Titel | Status | Phase |
|---|---|---|---|
| [0001](0001-three-layer-model.md) | Drei-Ebenen-Modell (CLI + Skills + MCP) | Accepted | 1 |
| [0002](0002-jsonl-statt-sqlite.md) | JSONL statt SQLite als Persistierung | Accepted | 1 |
| [0003](0003-claude-cli-shellout.md) | Planning via `claude --print`, kein SDK | Accepted | 2 |
| [0004](0004-cobra-als-cli-framework.md) | Cobra als CLI-Framework | Accepted | 1 |
| [0005](0005-anti-hype-defaults.md) | Konservative Anti-Hype-Defaults | Accepted | 1 |
| [0006](0006-deterministisches-scoring.md) | Deterministisches Scoring, keine LLM im Score-Pfad | Accepted | 1 |
| [0007](0007-skill-embed-strategie.md) | Skill-Embed: `flavors/` kanonisch, `internal/.../skills/data/` als Copy | Accepted | 2 |
| [0008](0008-apply-safety-gates.md) | Vier Safety-Gates für Auto-Apply | Accepted | 4 |

## Wann ADR statt Spec?

| Situation | Artefakt |
|---|---|
| Du planst, was als Nächstes gebaut wird | Spec (`docs/superpowers/specs/`) |
| Du dokumentierst, warum eine Design-Entscheidung getroffen wurde | ADR |
| Du beschreibst, wie ein laufendes System funktioniert | `docs/internals/*.md` |

ADRs sind kurz (1–2 Seiten), nüchtern und antworten auf eine
spezifische Frage. Querverweis von Spec auf ADR ist üblich, sobald
eine Spec-Entscheidung wiederkehrend relevant wird.
