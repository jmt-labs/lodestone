package memory

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jmt-labs/lodestone/internal/lodestone/audit"
)

func writeAuditFile(t *testing.T, dir string, entries []audit.Entry) string {
	t.Helper()
	log, err := audit.New(dir)
	if err != nil {
		t.Fatalf("audit.New: %v", err)
	}
	for _, e := range entries {
		if err := log.Record(e); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}
	return log.Path()
}

func TestConsolidateAppendsAndDedups(t *testing.T) {
	dir := t.TempDir()
	t0 := time.Date(2026, 5, 18, 0, 0, 0, 0, time.UTC)
	auditPath := writeAuditFile(t, dir, []audit.Entry{
		{Timestamp: t0, Verb: "fingerprint", Outcome: "ok", Detail: "languages=[Go]"},
		{Timestamp: t0.AddDate(0, 0, 1), Verb: "ingest", Outcome: "ok", Detail: "fetched=3"},
	})
	memPath := filepath.Join(dir, "memory.json")

	added, err := Consolidate(auditPath, memPath, WithSince(time.Time{}))
	if err != nil {
		t.Fatalf("Consolidate: %v", err)
	}
	if added != 2 {
		t.Fatalf("added = %d, want 2", added)
	}

	added2, err := Consolidate(auditPath, memPath, WithSince(time.Time{}))
	if err != nil {
		t.Fatalf("Consolidate 2: %v", err)
	}
	if added2 != 0 {
		t.Errorf("expected 0 added on second run (dedup), got %d", added2)
	}

	raw, _ := os.ReadFile(memPath)
	var f File
	if err := json.Unmarshal(raw, &f); err != nil {
		t.Fatalf("decode memory: %v", err)
	}
	if len(f.Decisions) != 2 {
		t.Errorf("got %d decisions, want 2", len(f.Decisions))
	}
}

func TestConsolidateRespectsSince(t *testing.T) {
	dir := t.TempDir()
	now := time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC)
	auditPath := writeAuditFile(t, dir, []audit.Entry{
		{Timestamp: now.AddDate(0, 0, -100), Verb: "old", Outcome: "ok"},
		{Timestamp: now.AddDate(0, 0, -1), Verb: "fresh", Outcome: "ok"},
	})
	memPath := filepath.Join(dir, "memory.json")
	added, err := Consolidate(auditPath, memPath, WithSince(now.AddDate(0, 0, -30)))
	if err != nil {
		t.Fatalf("Consolidate: %v", err)
	}
	if added != 1 {
		t.Errorf("added = %d, want 1 (only fresh)", added)
	}
}

func TestConsolidateMissingFiles(t *testing.T) {
	dir := t.TempDir()
	added, err := Consolidate(filepath.Join(dir, "absent.log"), filepath.Join(dir, "out.json"), WithSince(time.Time{}))
	if err != nil {
		t.Fatalf("Consolidate: %v", err)
	}
	if added != 0 {
		t.Errorf("added = %d, want 0", added)
	}
}
