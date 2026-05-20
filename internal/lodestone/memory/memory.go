package memory

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/jmt-labs/lodestone/internal/lodestone/audit"
)

const DefaultRelPath = ".claude/memory.json"

type File struct {
	Decisions []Entry `json:"decisions"`
}

type Entry struct {
	Date    string `json:"date"`
	Verb    string `json:"verb"`
	Summary string `json:"summary"`
}

type Options struct {
	Since time.Time
	Now   func() time.Time
}

func defaultOptions() Options {
	return Options{
		Since: time.Now().UTC().AddDate(0, 0, -90),
		Now:   func() time.Time { return time.Now().UTC() },
	}
}

func Consolidate(decisionsPath, memoryPath string, opts ...func(*Options)) (int, error) {
	o := defaultOptions()
	for _, fn := range opts {
		fn(&o)
	}

	entries, err := readAudit(decisionsPath, o.Since)
	if err != nil {
		return 0, err
	}

	existing, err := readMemory(memoryPath)
	if err != nil {
		return 0, fmt.Errorf("read existing memory: %w", err)
	}

	seen := map[string]struct{}{}
	for _, e := range existing.Decisions {
		seen[e.Date+"|"+e.Verb+"|"+e.Summary] = struct{}{}
	}

	added := 0
	for _, ae := range entries {
		date := ae.Timestamp.Format("2006-01-02")
		summary := ae.Detail
		if summary == "" {
			summary = ae.Outcome
		}
		key := date + "|" + ae.Verb + "|" + summary
		if _, ok := seen[key]; ok {
			continue
		}
		existing.Decisions = append(existing.Decisions, Entry{
			Date:    date,
			Verb:    ae.Verb,
			Summary: summary,
		})
		seen[key] = struct{}{}
		added++
	}

	sort.SliceStable(existing.Decisions, func(i, j int) bool {
		if existing.Decisions[i].Date != existing.Decisions[j].Date {
			return existing.Decisions[i].Date < existing.Decisions[j].Date
		}
		return existing.Decisions[i].Verb < existing.Decisions[j].Verb
	})

	if err := writeMemory(memoryPath, existing); err != nil {
		return added, err
	}
	return added, nil
}

func WithSince(t time.Time) func(*Options) {
	return func(o *Options) { o.Since = t }
}

func WithNow(fn func() time.Time) func(*Options) {
	return func(o *Options) { o.Now = fn }
}

func readAudit(path string, since time.Time) ([]audit.Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("open decisions.log: %w", err)
	}
	defer f.Close()

	var out []audit.Entry
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		var e audit.Entry
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			continue
		}
		if !e.Timestamp.IsZero() && e.Timestamp.Before(since) {
			continue
		}
		out = append(out, e)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan decisions.log: %w", err)
	}
	return out, nil
}

func readMemory(path string) (File, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return File{}, nil
		}
		return File{}, err
	}
	if len(raw) == 0 {
		return File{}, nil
	}
	var f File
	if err := json.Unmarshal(raw, &f); err != nil {
		return File{}, fmt.Errorf("decode memory: %w", err)
	}
	return f, nil
}

func writeMemory(path string, f File) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".tmp.*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(raw); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}
	return os.Rename(tmpName, path)
}
