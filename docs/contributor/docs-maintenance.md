# Doku-Wartung

Die Doku ist nicht „fertig", sondern Teil des Codes — sie muss mit
jeder Verhaltens- oder Architektur-Änderung mitwachsen. Diese Regeln
sorgen strukturell dafür, dass sie das tut.

## Regel: Doku updaten ist Teil der PR

Eine PR, die User-facing Verhalten ändert, muss die zugehörige Doku
mit aktualisieren. Welche Datei betroffen ist, hängt von der Änderung
ab:

| Geändert wurde | Doku-Update |
|---|---|
| Ein CLI-Verb (neues Flag, anderes Verhalten) | [`docs/user/commands/<verb>.md`](../user/commands/README.md) |
| Konfig-Schema (`.lodestone.yaml`, ENV) | [`docs/user/configuration.md`](../user/configuration.md) |
| Ingest-Source dazugekommen oder weg | [`commands/ingest.md`](../user/commands/ingest.md), [`commands/README.md`](../user/commands/README.md) |
| MCP-Tool dazugekommen oder weg | [`docs/user/mcp-server.md`](../user/mcp-server.md), [`commands/README.md`](../user/commands/README.md) |
| Skill dazugekommen | [`docs/user/skills.md`](../user/skills.md) (Index), Skill-File unter `flavors/lodestone/skills/` |
| Architektur, Datenmodell, Scoring geändert | passender `docs/internals/*.md` |
| Neue Design-Entscheidung | neuer [ADR](../internals/adr/README.md) |
| Neuer Begriff eingeführt | [`docs/user/glossary.md`](../user/glossary.md) |
| Phase abgeschlossen oder Stand geändert | [Roadmap](../internals/roadmap.md), `README.md § Status`, `CHANGELOG.md` (alle drei!) |

## Single Sources of Truth

| Thema | Datei |
|---|---|
| Phasen-Status | [`docs/internals/roadmap.md`](../internals/roadmap.md) |
| Skill-Inhalt | [`flavors/lodestone/skills/<name>.md`](../../flavors/lodestone/skills/) |
| Befehle-Liste | [`docs/user/commands/README.md`](../user/commands/README.md) |
| ADRs | [`docs/internals/adr/`](../internals/adr/README.md) |
| Glossar | [`docs/user/glossary.md`](../user/glossary.md) |

Andere Dateien zitieren diese — nie duplizieren.

## Strukturelle Absicherung

Make-Targets, die vor Merge laufen müssen:

| Target | Prüfung |
|---|---|
| `make docs-status-check` | Phasen-Status konsistent zwischen `README.md`, `docs/internals/roadmap.md` und `CHANGELOG.md`. |
| `make docs-cmd-coverage` | Für jedes Cobra-Verb existiert `docs/user/commands/<verb>.md`. |
| `make skills-coverage` | Jeder `flavors/lodestone/skills/*.md` wird in `docs/user/skills.md` verlinkt; gleichzeitig sind `flavors/` und `internal/lodestone/skills/data/` byte-identisch. |
| `make docs-links` | Relative Markdown-Links innerhalb von `docs/` und vom Root nach `docs/` zeigen auf existierende Pfade. |

Diese Targets sind Teil des `docs`-Jobs in
`.github/workflows/ci.yml` und schlagen bei Drift fehl. Sie ersetzen
manuelle Disziplin nicht — sie machen Drift sichtbar.

## Wann ADR statt Spec?

Eine Spec beantwortet „Was bauen wir als Nächstes?". Ein ADR
beantwortet „Warum haben wir uns für X entschieden — und werden diese
Entscheidung wahrscheinlich noch öfter brauchen?". Wenn eine
Spec-Entscheidung wiederkehrend relevant ist, gehört sie als ADR
hinterher festgehalten (auch rückwirkend) — siehe
[ADR-Index § Wann ADR statt Spec?](../internals/adr/README.md#wann-adr-statt-spec).

## Wann nicht doku-updaten

- Bei rein kosmetischen Refactorings ohne Verhaltensänderung.
- Bei internen Performance-Optimierungen, die nur Latenzen ändern.
- Bei Tests-Only-Changes ohne neuen Test-Hook.

Im Zweifel: lieber einen Satz im PR-Body erwähnen („keine Doku-Änderung
nötig, weil reines internes Refactoring") als sich später fragen, ob
es vergessen wurde.
