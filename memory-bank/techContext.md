# Tech Context — lodestone

## Stack

- **Sprache:** Go 1.24.7
- **CLI-Framework:** `cobra`
- **Config-Format:** YAML (`gopkg.in/yaml.v3`)
- Keine weiteren Produktions-Dependencies (YAGNI / Dependency-Footprint-Prinzip)

## Tools & Infrastruktur

| Tool | Zweck |
|---|---|
| `make build` | Binary bauen |
| `make test` | Unit-Tests ausführen |
| `make lint` | golangci-lint |
| `make e2e` | End-to-End-Tests in `e2e/` |
| `make vuln` | govulncheck |
| `gh` CLI | Issues, PRs, Releases |
| `codegraph` MCP | Semantische Code-Suche (7 Tools: search, node, callers, callees, explore, impact, files) |

## Externe Abhängigkeiten

- **GITHUB_TOKEN** (env) — GitHub-API für Signal-Abfragen und gh-Kommandos
- **Claude CLI** — für `lodestone plan` / `lodestone apply` (Phase 2, nur dort)

## Constraints

- LLM-Aufrufe **niemals** im Score-Pfad (ADR-0006, Determinismus-Garantie)
- Neue Dependencies brauchen Spec-Diskussion
- Go-Standardbibliothek bevorzugen
- Coverage-Ziel `internal/`: ≥ 70 %
