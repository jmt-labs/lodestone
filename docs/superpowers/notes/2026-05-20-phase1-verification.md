# Phase-1 Verifikations-Report

**Datum:** 2026-05-20
**Scope:** Abschluss-Checkliste für T11 aus dem Phase-1-Plan
([`../plans/2026-05-20-lodestone-mvp.md`](../plans/2026-05-20-lodestone-mvp.md)).
**Branch zum Zeitpunkt der Verifikation:** `feat/p1-t11-determinism`
(stacked auf T1–T10).

## Zusammenfassung

Alle Phase-1-Exit-Kriterien erfüllt. Pipeline `fingerprint → ingest →
score → signals` läuft offline mit Mock-Fixtures grün, ist
deterministisch und besteht alle statischen sowie dynamischen Checks
(modulo `govulncheck`, das nur im CI-Sandbox-Netzwerk erreichbar ist).

## Toolchain

- Go: `go1.24.7 linux/amd64` (CI nutzt Go 1.24 via `actions/setup-go@v5`)
- golangci-lint: `v2`-Schema, fünf aktive Linter (`errcheck`, `govet`,
  `ineffassign`, `staticcheck`, `unused`) plus `gofmt` als Formatter
- govulncheck: jeweils latest via `go install`

## Ergebnis-Matrix

| Check | Befehl | Status |
|---|---|---|
| Vet | `go vet ./...` | ✅ 0 Findings |
| Unit + Integration | `go test -race -count=1 ./...` | ✅ 6 Pakete grün, 0 Failures |
| Lint | `golangci-lint run` | ✅ 0 Issues (nach Errcheck-Exclude-Liste in `.golangci.yml`) |
| Vuln | `govulncheck ./...` | ⚠ lokal blockiert (Sandbox: HTTP 403 auf `vuln.go.dev`); CI führt es aus |
| Build | `go build ./...` | ✅ kompiliert sauber, keine Vendor-Reste |
| E2E | `make e2e` | ✅ fingerprint → ingest (mock) → score → signals → Determinismus-Diff |

### Determinismus

Zwei separate Mechanismen verifizieren Byte-Identität:

1. **Unit-Test** `internal/lodestone/scoring/scoring_test.go:TestScoreDeterminism`
   führt `scoring.Score` dreimal mit identischem Input aus,
   `json.Marshal`-t das Ergebnis und vergleicht die Bytes paarweise.
2. **E2E-Schritt** in `e2e/lodestone_test.sh`: nach dem ersten
   `lodestone score` wird `.lodestone/recommendations.jsonl` als
   Snapshot gesichert, ein zweiter `score`-Lauf folgt, dann `diff -q`
   gegen den Snapshot. Aktueller Lauf: identisch.

### Errcheck-Exclude-Liste

In `.golangci.yml` für Phase 1 ignoriert (idiomatische CLI-Muster,
keine Bug-Risiken):

- `fmt.Fprint{,f,ln}` — CLI-Ausgabe, Fehler praktisch nicht relevant
- `(io.Closer).Close` und `(*os.File).Close` — defer-Cleanup
- `os.Remove` — Best-Effort-Cleanup in Atomic-Rename-Pfaden
- `encoding/json.Encoder.Encode` — Test-Server-Antworten

Alle anderen Error-Returns werden geprüft.

## Backward-Compat

Nicht anwendbar in Phase 1: lodestone ist ein neues, eigenständiges
Projekt ohne vorhergehende Releases. Ab Phase 2 gilt: keine breaking
Changes am `.lodestone/`-Datei-Layout (`signals.jsonl`,
`fingerprint.json`, `recommendations.jsonl`) ohne Major-Bump.

## Bekannte Phase-2-Items

Keine Blocker für Phase-1-Abschluss, aber Notizen für später:

- `lodestone fingerprint` auf diesem Repo selbst zählt `testdata/`-
  Inhalte mit. Real-User-Repos haben üblicherweise keine
  `testdata/`-Verzeichnisse auf Root-Ebene, daher nicht
  exit-blockierend. Fix-Option: `testdata` in `skipDirs` aufnehmen
  (`internal/lodestone/fingerprint/fingerprint.go`).
- `lodestone signals` zeigt aktuell nur Stars/Source/Lang/Title.
  Eine spätere `--format`-Option für CSV/Markdown wäre nützlich
  für Reports.

## Sign-off

Phase 1 ist bereit für Tag `v0.1.0-alpha` (T12). Die GoReleaser-
Snapshot-Verifikation erfolgt in T12 separat.
