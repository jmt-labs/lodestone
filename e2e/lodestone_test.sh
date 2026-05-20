#!/usr/bin/env bash
# End-to-End-Smoke-Test für lodestone Phase 1.
#
# Erzeugt ein wegwerfbares Test-Repo, läuft fingerprint → ingest → score
# durch und prüft die erzeugten Artefakte unter .lodestone/.
#
# Erwartet $LODESTONE_BIN als Pfad zum gebauten Binary (default:
# bin/lodestone relativ zum Repo-Root).

set -euo pipefail

readonly REPO_ROOT="${REPO_ROOT:-$(cd "$(dirname "$0")/.." && pwd)}"
readonly LODESTONE_BIN="${LODESTONE_BIN:-$REPO_ROOT/bin/lodestone}"
readonly FIXTURE_DIR="${FIXTURE_DIR:-$REPO_ROOT/e2e/fixtures/signals}"

if [[ ! -x "$LODESTONE_BIN" ]]; then
  echo "lodestone-Binary nicht ausführbar: $LODESTONE_BIN" >&2
  echo "Tipp: 'make build' lokal ausführen." >&2
  exit 1
fi
if [[ ! -d "$FIXTURE_DIR" ]]; then
  echo "Fixture-Verzeichnis fehlt: $FIXTURE_DIR" >&2
  exit 1
fi

tmpdir="$(mktemp -d -t lodestone-e2e.XXXXXX)"
trap 'rm -rf "$tmpdir"' EXIT
echo "==> Test-Repo: $tmpdir"

cd "$tmpdir"

git init --quiet
git config user.email "e2e@lodestone.test"
git config user.name "E2E Bot"

cat > go.mod <<'EOF'
module example.com/e2e-fixture

go 1.24

require github.com/spf13/cobra v1.10.2
EOF

mkdir -p src
cat > src/main.go <<'EOF'
package main

import "github.com/spf13/cobra"

func main() { _ = &cobra.Command{} }
EOF

cat > src/main_test.go <<'EOF'
package main

import "testing"

func TestNoop(t *testing.T) {}
EOF

cat > .lodestone.yaml <<'EOF'
goals:
  - reliability
  - shipping
tech_interests:
  - mcp
lodestone:
  min_stars: 50
EOF

echo "==> lodestone fingerprint"
"$LODESTONE_BIN" fingerprint
test -f .lodestone/fingerprint.json
grep -q '"languages"' .lodestone/fingerprint.json
grep -q '"goals"' .lodestone/fingerprint.json

echo "==> lodestone ingest --mock"
export LODESTONE_MOCK_FIXTURES="$FIXTURE_DIR"
"$LODESTONE_BIN" ingest --mock
unset LODESTONE_MOCK_FIXTURES
test -f .lodestone/signals.jsonl
sig_count=$(wc -l < .lodestone/signals.jsonl)
if [[ "$sig_count" -lt 2 ]]; then
  echo "expected >= 2 signals, got $sig_count" >&2
  exit 1
fi

echo "==> lodestone score"
"$LODESTONE_BIN" score
test -f .lodestone/recommendations.jsonl
rec_count=$(wc -l < .lodestone/recommendations.jsonl)
if [[ "$rec_count" -lt 1 ]]; then
  echo "expected >= 1 recommendation, got $rec_count" >&2
  exit 1
fi

echo "==> lodestone signals --top 1 --json"
"$LODESTONE_BIN" signals --top 1 --json | grep -q '"id"'

echo "==> Determinismus (zweiter score-Lauf identisch?)"
cp .lodestone/recommendations.jsonl /tmp/recs.1.jsonl
"$LODESTONE_BIN" score >/dev/null
if ! diff -q /tmp/recs.1.jsonl .lodestone/recommendations.jsonl; then
  echo "score nicht deterministisch!" >&2
  exit 1
fi
rm -f /tmp/recs.1.jsonl

echo "==> OK"
