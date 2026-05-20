package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const DefaultFilename = "decisions.log"

type Entry struct {
	Timestamp time.Time         `json:"ts"`
	Verb      string            `json:"verb"`
	Args      map[string]string `json:"args,omitempty"`
	Outcome   string            `json:"outcome"`
	Detail    string            `json:"detail,omitempty"`
}

type Log struct {
	path string
	mu   sync.Mutex
	now  func() time.Time
}

func New(root string) (*Log, error) {
	if root == "" {
		return nil, fmt.Errorf("audit: empty root")
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, fmt.Errorf("audit mkdir: %w", err)
	}
	return &Log{
		path: filepath.Join(root, DefaultFilename),
		now:  func() time.Time { return time.Now().UTC() },
	}, nil
}

func (l *Log) WithNow(fn func() time.Time) *Log {
	l.now = fn
	return l
}

func (l *Log) Record(e Entry) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if e.Timestamp.IsZero() {
		e.Timestamp = l.now()
	}
	raw, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("marshal entry: %w", err)
	}
	f, err := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open %s: %w", l.path, err)
	}
	defer f.Close()
	if _, err := f.Write(append(raw, '\n')); err != nil {
		return fmt.Errorf("write entry: %w", err)
	}
	return nil
}

func (l *Log) Path() string {
	return l.path
}
