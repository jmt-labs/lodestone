<!-- GENERATED:BEGIN -->
# Claude-Konfiguration

Dieses Repository nutzt eine reproduzierbare forgecrate-Konfiguration. Die hier
beschriebenen Regeln gelten für alle Agenten (Claude Code, Codex, …) die im Repo
arbeiten. Die generierten Abschnitte werden bei `forgecrate update` überschrieben —
eigene Anpassungen gehören in den CUSTOM-Abschnitt der Root-`CLAUDE.md`.

## Pflicht-Skills

| Situation | Skill | Verhalten |
|---|---|---|
| Neues Feature / Bug-Fix | `superpowers:brainstorming` | MUSS vor Code aufgerufen werden |
| Nach Brainstorming | `forgecrate-roadmap-triage` | MUSS aufgerufen werden — entscheidet ob jetzt oder Future Feature |
| Implementierung | `superpowers:test-driven-development` | MUSS vor Code aufgerufen werden |
| Vor jeder nicht-trivialen Änderung | `forgecrate-research` | MUSS aufgerufen werden |
| Vor Commit/PR | `superpowers:verification-before-completion` | MUSS ausgeführt werden |
| Debug | `superpowers:systematic-debugging` | MUSS vor Fix aufgerufen werden |
| Bug gefunden (nach Debug) | `superpowers:test-driven-development` | Regressionstest schreiben, BEVOR der Fix committed wird |
| Session-Start | `mcp__memory__read_graph` | Projektübergreifendes Wissen laden |
| Architekturentscheidung / Debugging-Ergebnis | `mcp__memory__add_observations` | In memory MCP schreiben |

**Codegraph-Pflicht** (wenn codegraph-Flavor aktiv): Vor jeder nicht-trivialen Änderung `codegraph_node` + `codegraph_callers` für betroffene Symbole ausführen — kein Edit/Write ohne vorherige Codegraph-Abfrage.

## Recherche-Pflicht

**Alle** Rollen MÜSSEN vor jeder nicht-trivialen Code-Änderung mindestens ein
Recherche-Tool nutzen — statt aus gelerntem Wissen zu arbeiten. Raten ist verboten;
Quellen werden referenziert. Der `pre-tool.sh`-Hook **warnt** bei fehlender Recherche,
blockiert aber nicht.

| Frage-Typ | Tool | Beispiele |
|---|---|---|
| Library-/Framework-Doku | `context7` | API-Syntax, Migrationen, Versions-Updates |
| Spezifische URL aus Issue/Ticket | `fetch` MCP | RFCs, MDN, Changelogs |
| Allgemeine Web-Recherche | `WebSearch` | Best Practices, Vergleiche, aktuelle Probleme |

**Regeln:**

- Mindestens eine Quelle pro nicht-trivialer Entscheidung; eine Recherche pro Session
  schaltet weitere Warnungen für die Session ab
- Quellen im Plan-Dokument (`docs/superpowers/plans/*.md`) referenzieren
- Deaktivierbar via Flavor `no-research`

## Entwicklungs-Workflow

Für alle Features, Bugfixes und Änderungen:

1. **Brainstorming** — `superpowers:brainstorming` aufrufen, Design abstimmen
2. **Spec** — Branch anlegen (`git checkout -b feat/<thema>`); Spec in
   `docs/superpowers/specs/YYYY-MM-DD-<thema>-design.md` schreiben und committen;
   GitHub-Issue anlegen oder verlinken; Branch-Name im Issue vermerken; Kommentar
   im Issue: "Spec fertig"
3. **Plan** — in `docs/superpowers/plans/YYYY-MM-DD-<thema>.md` schreiben und
   committen; Plan-Pfad im Issue ergänzen; Kommentar: "Plan fertig"
4. **Implementierung** — nach jedem Task kurzer Kommentar im Issue
5. **PR & Abschluss** — Vor `gh pr create` diese Sequenz vollständig ausführen:
   1. `forgecrate-doc-sync` — Doku mit Code abgleichen
   2. `forgecrate-handoff` — memory-bank aktualisieren (`activeContext.md`, `progress.md`)
   3. `forgecrate-db-migration` — Migrations-Review
   4. `accessibility-audit` — A11y-Prüfung
   5. `ui-ux-audit` — UX-Review
   6. `forgecrate-pr-checklist` — Abschluss-Checkliste

   Dann PR erstellen, Issue im PR-Body verlinken ("Closes #N").
   Issue wird nach Merge automatisch geschlossen.

Ticket-Kommentare immer kurz (ein Satz): Fortschritt, Pfad oder Ergebnis.

## Session-Start

Beim Session-Start: aktuellen Projektkontext aus der memory-bank lesen.
**Pflicht:** `mcp__memory-bank__memory_bank_read` verwenden — direktes Lesen via
Read-Tool auf `memory-bank/`-Dateien ist verboten.

## Verhalten

- Antworte auf Deutsch
- Keine unnötigen Kommentare im Code
- YAGNI: keine ungefragten Features
- Änderungen immer über Branch + PR, nie direkt auf `main`

## Hook-Schutz: Hinweis

Der `pre-tool.sh`-Hook **warnt** bei destruktiven Bash-Befehlen und fehlender
Recherche — er blockiert nie. Die Verantwortung liegt beim Agenten: Warnungen
bewusst wahrnehmen, einschätzen und eine informierte Entscheidung treffen.

Für serverseitigen Schutz auf `main`: GitHub Branch Protection Rules konfigurieren.

Bei fehlender Binary, fehlendem oder kaputtem Transcript verhält sich der Hook
**fail-open** (keine Warnung).

## Team-Rollen & Subagent-Konfiguration

Der Hauptagent koordiniert als Team-Lead. Subagenten übernehmen Rollen
entsprechend ihrer Aufgabe. Der Hauptagent kann bei Bedarf eigenständig von
diesen Empfehlungen abweichen.

Das Hauptmodell der Session ist global (in `.claude/settings.json`). Die
`Modell`-Spalte nennt den empfohlenen Wert für den `model`-Parameter beim
Dispatch eines Subagenten über das Agent-Tool — gültig sind nur die Family-Aliase
`opus`/`sonnet`/`haiku`.

| Rolle | Superpowers-Skill | Modell | Recherche |
|---|---|---|---|
| Analyst / Product Owner | `superpowers:brainstorming` | `opus` | Pflicht |
| Tech Lead / Architekt | `superpowers:writing-plans` | `opus` | Pflicht |
| Entwickler | `superpowers:test-driven-development` | `sonnet` | Pflicht |
| Implementierer (mechanisch) | `superpowers:subagent-driven-development` | `haiku` | Pflicht |
| Reviewer | `superpowers:requesting-code-review` | `sonnet` | Pflicht |
| QA / Abschluss | `superpowers:verification-before-completion` | `sonnet` | Pflicht |
| Debugger | `superpowers:systematic-debugging` | `sonnet` | Pflicht |

## Parallelisierung & Isolation

Subagenten werden proaktiv parallelisiert und isoliert — ohne explizite
Aufforderung.

| Situation | Mechanismus | Anleitung |
|---|---|---|
| Task dauert >1 min oder Ergebnis nicht sofort nötig | `run_in_background: true` | `superpowers:dispatching-parallel-agents` |
| Feature-Branch, Multi-File-Änderung, langer Plan | `isolation: "worktree"` | `superpowers:using-git-worktrees` |
| Mehrere unabhängige Tasks gleichzeitig | beide kombinieren | beide Skills |

Im Zweifelsfall Background nutzen — warten ist kein Default.

### Agenten-Identität

Jeder Subagent bekommt eindeutige Identifikation:

- **Eindeutigen Namen** — via `description`-Parameter im Agent-Tool-Aufruf
  (3–5 Wörter, Rolle + Aufgabe)
- **Eindeutige Farbe** — dynamisch durch FleetView-Dashboard zugewiesen; keine
  zwei gleichzeitig laufenden Agenten teilen eine Farbe

Dies ermöglicht einfaches Tracking und verhindert Verwechslungen bei parallelen
Läufen.

## MCP-Server

Sechs MCP-Server stehen automatisch zur Verfügung. `.mcp.json` wird von forgecrate
generiert — nicht von Hand editieren; MCP-Server-Änderungen über einen erneuten
forgecrate-Lauf.

| Server | Transport | Zweck |
|---|---|---|
| `github` | stdio (`npx`) | Issues, PRs, Code-Suche, Branches, Labels |
| `fetch` | stdio (`npx`) | Externe Webinhalte: Docs, RFCs, Changelogs |
| `memory` | stdio (`npx`) | Projektübergreifende Architektur-Entscheidungen |
| `memory-bank` | stdio (`npx`) | Repo-spezifischer Projektkontext (laufender Stand) |
| `context-mode` | stdio (`npx`) | Automatisches Context-Budget und Session-History-Suche |
| `context7` | stdio (`npx`) | Aktuelle Bibliotheks-Dokumentation aus Source-Repos |

Routing-Grenzen (verhindern Falsch-Aufrufe):

- **`github`** — alle GitHub-Operationen (Issues, PRs, Code-Suche, Labels). NICHT für
  lokale Datei-/Git-Kommandos (→ Read/Edit/Bash). Voraussetzung:
  `GITHUB_PERSONAL_ACCESS_TOKEN`.
- **`fetch`** — externe Webinhalte (Docs, MDN, RFCs, Changelogs). NICHT für
  GitHub-Inhalte (→ `github`) oder lokale Dateien (→ Read).
- **`context-mode`** — sandboxt Tool-Output automatisch (kein Aufruf nötig). Explizit:
  `ctx_search` (History-Suche nach Kompaktierung), `ctx_stats`, `ctx_doctor`.
- **`context7`** — aktuelle Bibliotheks-Doku aus Source-Repos. NICHT für GitHub-Inhalte
  (→ `github`), lokale Dateien (→ Read) oder allgemeine Programmierkonzepte.

`memory` und `memory-bank` haben eigene Pflicht-Regeln — siehe unten.

## Claude Plugins

Vier Plugins werden automatisch via `forgecrate deploy` installiert (`claude plugin install --scope project`).

| Plugin | Zweck |
|---|---|
| `superpowers` | Skill-System: Workflows für TDD, Brainstorming, Debugging, Reviews |
| `commit-commands` | Slash-Commands für standardisierte Commits und PRs |
| `security-guidance` | Sicherheitshinweise und Best-Practices für Code-Reviews |
| `claude-md-management` | Verwaltung und Verbesserung von CLAUDE.md-Dateien |

Plugins stellen Slash-Commands und Skills bereit — sie sind nicht über MCP aufrufbar.

### Memory (`memory`)

Projektübergreifendes Wissen persistent speichern. Datei: `.claude/memory.json`
(versioniert).

**Schreiben nach:** Architekturentscheidungen, Begründungen für nicht-
offensichtliche Lösungen, Debugging-Ergebnisse, Brainstorming-Ergebnisse.

**Lesen am:** Sessionbeginn, nach Context-Kompaktierung, wenn unklar warum etwas
so gebaut wurde.

**Niemals speichern:** API-Keys, Tokens, Passwörter, temporären Zwischenstand,
Code-Details die direkt aus dem Code lesbar sind.

### Memory-Bank (`memory-bank`)

Repo-spezifischer, strukturierter Projektkontext im Verzeichnis `memory-bank/`
(versioniert, committed). Persistiert kontextuelles Wissen über Sessions hinweg.

**Dateien:**

- `projectbrief.md` — Projektziel und Scope
- `techContext.md` — Stack, Tools, technische Constraints
- `systemPatterns.md` — Architektur-Entscheidungen, ADRs, Anti-Patterns
- `activeContext.md` — Aktueller Fokus, offene Fragen, Blocker
- `progress.md` — Was fertig ist, was läuft, was als nächstes kommt

**Lesen** am Session-Start und bei Bedarf — **ausschließlich** via
`mcp__memory-bank__memory_bank_read`.

**Schreiben** wenn sich Fokus, Fortschritt oder Architektur-Kontext ändert —
**ausschließlich** via `mcp__memory-bank__memory_bank_write` oder
`mcp__memory-bank__memory_bank_update`.

> **Direkte Datei-Tools (Read/Write/Edit) auf `memory-bank/`-Dateien sind
> verboten.**

**Abgrenzung zu `memory`:** `memory-bank` ist repo-spezifisch und dateibasiert —
ideal für laufenden Projekt-Kontext. `memory` (`.claude/memory.json`) ist
graph-basiert und projektübergreifend — ideal für zeitlose
Architektur-Entscheidungen mit Begründung.

## Backend-Profil

- API-Design: REST-First, klare Fehlercodes, keine unnötige Abstraktion
- Datenbankzugriffe: typsicher, keine Raw-Queries ohne Parametrisierung
- Tests: Integrationstests bevorzugt gegenüber reinen Unit-Tests mit Mocks
- Kein ORM-Magic: explizite Queries sind verständlicher

## Codegraph-Flavor

Dieses Repo nutzt **codegraph** — einen semantischen Code-Wissensgraphen als MCP-Server.

### Was codegraph bietet

Der MCP-Server läuft lokal (`codegraph serve --mcp`) und stellt folgende Tools bereit:

| Tool | Zweck |
|---|---|
| `codegraph_search` | Semantische Code-Suche ohne exakte Schlüsselwörter |
| `codegraph_node` | Definition eines Symbols (Funktion, Typ, Variable) abrufen |
| `codegraph_callers` / `codegraph_callees` | Alle Aufrufer / Aufgerufenen eines Symbols |
| `codegraph_explore` | Abhängigkeiten und Nachbarn eines Symbols erkunden |
| `codegraph_impact` | Blast-Radius einer Änderung ermitteln |
| `codegraph_files` | Dateien im Index auflisten |
| `codegraph_status` | Index-Status prüfen |

### Pflicht-Regeln

- **Vor jeder nicht-trivialen Änderung MUSS** `codegraph_node` + `codegraph_callers` für betroffene Symbole aufgerufen werden — kein Edit/Write ohne vorherige Codegraph-Abfrage
- **Beim Debuggen MUSS** `codegraph_explore` die Aufrufkette aufzeigen, bevor ein Fix versucht wird
- **Bei Refactoring MUSS** `codegraph_callers` für Call-Sites + `codegraph_search` für Type-/Import-Referenzen geprüft werden
- **Code-Suche**: `codegraph_search` statt grep — grep ist nur erlaubt, wenn codegraph das Ergebnis nicht liefert
- **Impact-Analyse MUSS** `codegraph_impact` vor größeren Umbauten ausgeführt werden

### Index-Aktualisierung

Der Index wird automatisch bei Session-Start im Hintergrund aktualisiert (einmal pro Commit-Stand).
Manuell: `codegraph index` im Repo-Root. Erstmalige Initialisierung: `codegraph init -i`.

### Voraussetzung

Installation (einmalig, kein Node.js erforderlich):

```bash
# macOS / Linux
curl -fsSL https://raw.githubusercontent.com/colbymchenry/codegraph/main/install.sh | sh

# Windows (PowerShell)
irm https://raw.githubusercontent.com/colbymchenry/codegraph/main/install.ps1 | iex

# Alternativ via npm
npm i -g @colbymchenry/codegraph
```

Danach im Repo initialisieren:

```bash
codegraph init -i
```

Der MCP-Server wird über `.mcp.json` automatisch konfiguriert.

## Force-Research-Flavor (verschärfte Recherche-Pflicht)

Dieser Flavor verschärft die Recherche-Empfehlung des base layer:

- Die Recherche-Warnung (kein Edit/Write/MultiEdit ohne vorherige Recherche
  einmal pro Session) gilt **zusätzlich für schreibende Bash-Befehle** — auch
  Datei-Schreibzugriffe via Shell (`sed -i`, `tee`, `dd of=`, Redirects außerhalb
  `/tmp`) erzeugen ohne vorherige Recherche eine Warnung. Damit ist die Umgehung
  „Datei per Shell schreiben statt Edit/Write" geschlossen.
- Kein impliziter Ausnahmefall. Bewusster Verzicht ausschließlich über den Flavor
  `no-research`.

Die Durchsetzung liegt vollständig im base-Hook (`pre-tool.sh` →
`forgecrate hook require-research`); dieser Flavor aktiviert lediglich die
zusätzliche Bash-Prüfung über die aktive Konfiguration.

## GETBETTER-Flavor

Kontinuierliche Verbesserung durch Festhalten von Erkenntnissen aus jeder Session.

- Am Session-Start: `mcp__memory__read_graph` aufrufen, Entities vom Typ `session-reflection` lesen.
- Am Sessionende: `/forgecrate-getbetter` aufrufen um Erkenntnisse zu speichern.

**Was gespeichert wird (memory MCP, Entity `session-reflection`):**
- Wiederkehrende Fehler und deren Ursachen
- Patterns die gut funktioniert haben
- Entscheidungen die sich im Nachhinein als falsch erwiesen haben
- Projektspezifische Gotchas die nicht aus dem Code ersichtlich sind

**Format pro Erkenntnis:** `[YYYY-MM-DD] <Kategorie>: <Erkenntnis in einem Satz>`
Kategorien: `workflow`, `tooling`, `pattern`, `mistake`, `decision`.

## GitHub-Flavor

- Releases über `gh release create` veröffentlichen (nach `release`-Skill)
- PR-Templates in `.github/pull_request_template.md` pflegen
- CI-Status mit `gh run list` prüfen bevor ein Release getaggt wird

## Multiagent & Subagenten

Parallelisierung/Isolation gemäß Base-Layer-Tabelle gelten auch hier — gerade bei
Issue-Batches proaktiv Background-Mode und Worktrees nutzen.

## Strict-Review-Flavor

- Vor jedem Commit: `superpowers:requesting-code-review` aufrufen
- Keine direkten Commits auf main/master
- PR-Beschreibung enthält: Was, Warum, Wie getestet
- Breaking Changes werden explizit kommuniziert

## TDD-Flavor

- Test schreiben → ausführen (muss fehlschlagen) → implementieren → ausführen (muss bestehen) → committen
- Kein Produktionscode ohne vorherigen Test
- Test-Namen beschreiben Verhalten, nicht Implementierung
- Mocks nur an Systemgrenzen (externe APIs, Datenbanken)
- Für jeden gefundenen Bug: Regressionstest vor dem Fix
<!-- GENERATED:END -->

<!-- CUSTOM:BEGIN -->
# CLAUDE.md — Vorgaben für Claude in diesem Repo

> Diese Datei wird von Claude Code automatisch gelesen. Sie definiert
> die nicht-verhandelbaren Konventionen für jegliche KI-getriebene
> Arbeit an `lodestone`.

## Sprache

**Antworte auf Deutsch.** Commit-Messages, PR-Bodies, Doku, Specs und
Pläne auf Deutsch. Code-Identifier, API-Felder, Log-Messages und
Datei-Inhalte technischer Natur (z. B. `go.mod`) bleiben englisch.

## Pflicht-Skills

| Situation | Skill | Anforderung |
|---|---|---|
| Feature / Bug-Fix | `superpowers:brainstorming` | **MUSS vor Code aufgerufen werden** |
| Implementierung | `superpowers:test-driven-development` | **MUSS vor Code aufgerufen werden** |
| Vor Commit / PR | `superpowers:verification-before-completion` | **MUSS ausgeführt werden** |
| Debugging | `superpowers:systematic-debugging` | **MUSS vor Fix aufgerufen werden** |

**Bug-Regel:** Bei jedem Bug schreibst du den **Regressionstest ZUERST**,
committest ihn separat, und committest dann den Fix.

## Branch- & PR-Workflow

1. Brainstorming → Design abgestimmt.
2. **Spec** in `docs/superpowers/specs/YYYY-MM-DD-<thema>-design.md`.
3. **Plan** in `docs/superpowers/plans/YYYY-MM-DD-<thema>.md` mit
   Checkbox-Tasks (`- [ ] Tx: …`).
4. **Branch** anlegen — Schema `feat/p<phase>-t<task>-<slug>` für
   Plan-Tasks, `fix/<slug>` für Bugs, `chore/<slug>` für Wartung.
5. **TDD-Implementierung** mit kurzen, atomaren Commits.
6. **PR gegen `main`** mit Body, der die Spec/Plan-Datei verlinkt und
   das Epic-Issue mit `Updates #N` referenziert.

### Niemals direkt auf `main`

**Ausnahme:** Der erste Bootstrap-Commit dieses Repos wurde explizit auf
`main` autorisiert. Danach gilt: **alles über Branch + PR**, nie direkt
auf `main`. Kein force-push auf `main`.

### PR-Erstellung

PRs werden **nur auf explizite Aufforderung des Users erstellt**.
Push auf Feature-Branch ist erlaubt, PR-Eröffnung erfordert
expliziten Auftrag.

## YAGNI & Code-Qualität

- Keine spekulativen Features.
- Keine Abstraktionen ohne konkreten zweiten Aufrufer.
- Keine Error-Handling-Pfade für Szenarien, die nicht eintreten können.
- Keine Backward-Compat-Shims für Code, der noch nicht released ist.
- **Keine Kommentare im Code**, außer das WARUM ist non-obvious
  (versteckte Constraint, subtile Invariante, Workaround für Bug).

## Historische Invarianten (Phase 1) und aktuelle Geltung

Diese vier Invarianten wurden für Phase 1 aufgestellt. Punkte 1 und 3
sind ab Phase 2 explizit gelockert; Punkte 2 und 4 gelten unverändert
weiter und sind als ADRs festgehalten.

1. **~~Keine LLM-Aufrufe~~** — ab Phase 2 erlaubt in `lodestone plan`
   und `lodestone apply` (über die `claude`-CLI). **Niemals** im
   Score-Pfad — siehe
   [ADR-0006](docs/internals/adr/0006-deterministisches-scoring.md).
2. **Deterministische Pipeline** — zwei Score-Läufe mit identischem
   Input müssen byte-identische sortierte Outputs liefern. Gilt
   weiter, siehe
   [ADR-0006](docs/internals/adr/0006-deterministisches-scoring.md)
   und [Determinismus](docs/internals/determinism.md).
3. **~~Minimaler Dependency-Footprint~~** — ab Phase 2 gelockert auf
   „Standardbibliothek bevorzugen; neue Deps brauchen
   Spec-Diskussion". Aktueller Stand: weiter nur `cobra` und
   `yaml.v3`.
4. **Anti-Hype-Defaults konservativ** — `min_stars: 50`,
   `min_age_days: 30`, `max_last_commit_age_days: 180`,
   `require_license: true`. Gilt weiter, siehe
   [ADR-0005](docs/internals/adr/0005-anti-hype-defaults.md).

Detaillierte Contributor-Vorgaben:
[`docs/contributor/workflow.md`](docs/contributor/workflow.md),
[`docs/contributor/skills-policy.md`](docs/contributor/skills-policy.md).
Phasen-Status: [`docs/internals/roadmap.md`](docs/internals/roadmap.md).

## Modell-Routing

Siehe `base/models.yaml`. Kurzfassung:

- `planning` → `claude-opus-4-7` (Specs, Pläne, Architektur)
- `default` → `claude-sonnet-4-6` (Implementierung, Reviews)
- `mechanical` → `claude-haiku-4-5-20251001` (Roundtrips,
  Format-Konvertierung, Rationale-Generierung)
- `review` → `claude-sonnet-4-6` (PR-Review, Spec-Critique)

## Testing

- **Coverage-Ziel** für `internal/`-Pakete: ≥ 70 %.
- Tests beschreiben Verhalten, nicht Implementierung.
- E2E in `e2e/`, ausführbar via `make e2e`.
- Vor PR: `make test lint vuln` muss grün sein.

## Was lodestone NICHT ist

- Kein Daemon. Keine Hintergrund-Prozesse.
- Kein Telemetrie-Sender. Lokale Artefakte bleiben lokal.
- Kein Auto-Editor für `main`-Branches. Auto-PRs (Phase 4) immer auf
  Feature-Branch, immer als Draft, immer mit harten Schranken.

<!-- CUSTOM:END -->
