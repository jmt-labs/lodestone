# Testing

Lodestones Test-Strategie folgt einer schlanken Pyramide. Tests
beschreiben Verhalten, nicht Implementierung. Coverage-Ziel für
`internal/`-Pakete: **≥ 70 %**.

## Test-Pyramide

```
Unit-Tests (internal/lodestone/*)
  └── httptest.Server für Source-Adapter
  └── FakeRunner / FakeGit / FakePR für I/O-lastige Pfade
  └── Golden-Fixtures unter internal/lodestone/fingerprint/testdata/

E2E (e2e/lodestone_test.sh)
  └── Wegwerf-Repo via mktemp, init + fingerprint + ingest --mock +
      score + plan --dry-run + decisions.log-Check + Determinismus-Diff

CI (.github/workflows/ci.yml)
  └── test (vet + race) · lint (golangci-lint v2) · vuln (govulncheck)
      · e2e · shellcheck · readme-coverage
```

## Unit-Tests

- **HTTP-Mocks:** `httptest.Server` für jeden Source-Adapter, mit
  Golden-Response-Fixtures unter `internal/lodestone/ingest/testdata/`.
- **Pluggable Runner:** `FakeRunner` (Planning), `FakeGit` / `FakePR`
  (Apply) erlauben deterministische Tests ohne echte
  `claude`/`git`/`gh`-Aufrufe.
- **Determinismus:** `TestScoreDeterminism` in
  `internal/lodestone/scoring/scoring_test.go` führt drei Score-Läufe
  und vergleicht die JSON-Bytes. Siehe
  [Determinismus](../internals/determinism.md).

## E2E

`e2e/lodestone_test.sh` ist ein Bash-Skript, das ein Wegwerf-Repo via
`mktemp` aufsetzt, die komplette Pipeline durchläuft und alle
Artefakte verifiziert. Offline-fähig über `LODESTONE_MOCK_FIXTURES`.

```sh
make e2e
```

Schritte:

1. `lodestone init` → `.lodestone.yaml` und `.gitignore` anlegen.
2. `lodestone fingerprint` → `.lodestone/fingerprint.json`.
3. `lodestone ingest --mock` → `.lodestone/signals.jsonl` aus Fixtures.
4. `lodestone score` → `.lodestone/recommendations.jsonl`.
5. `lodestone plan --dry-run` → Prompt-Check, kein echter Claude-Call.
6. Determinismus-Diff: zweiter `score`-Lauf, `diff -q` gegen Snapshot.
7. `decisions.log` muss alle Verben in der richtigen Reihenfolge
   enthalten.

## Lokale Test-Targets

```sh
make test          # go test ./... (mit -race)
make lint          # golangci-lint v2
make vuln          # govulncheck
make e2e           # End-to-End-Smoke-Test
```

Alle vier müssen vor jedem PR grün sein. Siehe
[PR-Checkliste](pr-checklist.md).

## Was nicht getestet wird

- **Echte externe APIs.** Sources werden gegen `httptest.Server`
  getestet, nicht gegen Live-GitHub/HackerNews. Spart Rate-Limits und
  hält Tests deterministisch.
- **Echte LLM-Calls.** Planning-Tests laufen gegen `FakeRunner` mit
  hartcodierten Responses; Live-Claude-Aufrufe nur in manuellen
  Smoke-Tests vor Releases.
- **Echte Git/GitHub-Operationen.** Apply-Tests laufen gegen
  `FakeGit` / `FakePR`. Reale Aufrufe nur im E2E gegen ein
  Wegwerf-Repo, ohne `--push`.
