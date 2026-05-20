package audit

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRecordAppends(t *testing.T) {
	dir := t.TempDir()
	log, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	now := time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
	log.WithNow(func() time.Time { return now })

	for i, verb := range []string{"ingest", "score", "plan"} {
		if err := log.Record(Entry{
			Verb:    verb,
			Args:    map[string]string{"i": "x"},
			Outcome: "ok",
			Detail:  "test entry",
		}); err != nil {
			t.Fatalf("Record[%d]: %v", i, err)
		}
	}

	raw, err := os.ReadFile(filepath.Join(dir, DefaultFilename))
	if err != nil {
		t.Fatalf("read decisions.log: %v", err)
	}
	lines := strings.Split(strings.TrimRight(string(raw), "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("got %d lines, want 3", len(lines))
	}
	for _, line := range lines {
		var e Entry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			t.Fatalf("unmarshal: %v (line: %s)", err, line)
		}
		if !e.Timestamp.Equal(now) {
			t.Errorf("Timestamp = %v, want %v", e.Timestamp, now)
		}
		if e.Outcome != "ok" {
			t.Errorf("Outcome = %q", e.Outcome)
		}
	}
}

func TestRecordPreservesExplicitTimestamp(t *testing.T) {
	dir := t.TempDir()
	log, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	ts := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	if err := log.Record(Entry{Verb: "test", Outcome: "ok", Timestamp: ts}); err != nil {
		t.Fatalf("Record: %v", err)
	}
	f, _ := os.Open(filepath.Join(dir, DefaultFilename))
	defer f.Close()
	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		t.Fatal("no line")
	}
	var e Entry
	if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !e.Timestamp.Equal(ts) {
		t.Errorf("Timestamp = %v, want %v", e.Timestamp, ts)
	}
}

func TestNewEmptyRoot(t *testing.T) {
	if _, err := New(""); err == nil {
		t.Fatal("expected error on empty root")
	}
}
