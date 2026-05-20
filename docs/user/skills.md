# Lodestone-Skills

Vier Claude-Skills bündeln typische Lodestone-Workflows in
LLM-ausführbarer Form. Sie werden vom Binary über `go:embed`
mitgeliefert und durch `lodestone init` ins User-Projekt installiert.

## Was ist ein Skill?

Ein Skill ist ein Markdown-File mit YAML-Frontmatter (`name`,
`description`), das Claude-Skills automatisch lädt und über
Trigger-Phrasen aktiviert. Skills bauen auf der Lodestone-CLI auf —
sie ersetzen sie nicht.

## Installation

```sh
cd dein-projekt && lodestone init
```

Das legt die Skill-Dateien unter `.claude/skills/` ab. Anderer Ziel-Pfad
über `lodestone init --skills-dir <pfad>`. Details:
[`init`-Befehl](commands/init.md).

## Verfügbare Skills

Quelle ist immer [`flavors/lodestone/skills/`](../../flavors/lodestone/skills/)
(kanonisch, einziger editierbarer Pfad — siehe
[ADR-0007](../internals/adr/0007-skill-embed-strategie.md)).

| Skill | Trigger | Use-Case |
|---|---|---|
| [`lodestone-scout`](../../flavors/lodestone/skills/lodestone-scout.md) | „scout", „neue Signale", „was läuft in AI", „trends fürs Repo" | Ingestion + Triage frischer AI-Ökosystem-Signale |
| [`lodestone-recommend`](../../flavors/lodestone/skills/lodestone-recommend.md) | „was sollte ich nächstes einbauen", „Empfehlungen", „was lohnt sich" | Interaktive Top-N-Priorisierung der Recommendations |
| [`lodestone-plan`](../../flavors/lodestone/skills/lodestone-plan.md) | „plan für diesen Vorschlag", „spec für rec-id X" | Generiert Spec + Plan aus einer Recommendation |
| [`lodestone-review-trends`](../../flavors/lodestone/skills/lodestone-review-trends.md) | „review trends", „was haben wir verpasst", „lodestone-stats" | Periodischer Review-Report |

## Eigene Skills hinzufügen

Lodestone-Skills sind kein geschlossenes Set. Eigene Skills, die die
CLI-Verben nutzen, gehören direkt nach `.claude/skills/<eigener-name>.md`
und folgen dem gleichen Frontmatter-Schema. Wenn dein Skill für andere
Lodestone-User nützlich ist, schlage ihn als Spec im
[Contributor-Workflow](../contributor/workflow.md) vor.

## Beziehung zur CLI und zum MCP-Server

Skills, CLI und MCP-Server sind drei Frontends desselben Kerns
([ADR-0001](../internals/adr/0001-three-layer-model.md)):

- **CLI** für Menschen am Terminal — siehe [Befehle](commands/README.md).
- **Skills** für Claude-Konversationen.
- **MCP-Server** für programmatischen Zugriff aus Claude Desktop,
  Claude Code oder IDEs — siehe [MCP-Server](mcp-server.md).

Alle drei rufen `internal/lodestone/*` auf — keine Doppelung.
