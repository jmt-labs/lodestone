package store

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/jmt-labs/lodestone/internal/lodestone/schema"
)

const DefaultRoot = ".lodestone"

const (
	signalsFile         = "signals.jsonl"
	fingerprintFile     = "fingerprint.json"
	recommendationsFile = "recommendations.jsonl"
	cacheDir            = "cache"
)

type FileStore struct {
	root string

	mu     sync.Mutex
	idx    map[string]struct{}
	loaded bool
}

func New(root string) (*FileStore, error) {
	if root == "" {
		root = DefaultRoot
	}
	if err := os.MkdirAll(filepath.Join(root, cacheDir), 0o755); err != nil {
		return nil, fmt.Errorf("create root: %w", err)
	}
	return &FileStore{root: root}, nil
}

func (s *FileStore) signalsPath() string         { return filepath.Join(s.root, signalsFile) }
func (s *FileStore) fingerprintPath() string     { return filepath.Join(s.root, fingerprintFile) }
func (s *FileStore) recommendationsPath() string { return filepath.Join(s.root, recommendationsFile) }

func (s *FileStore) ensureIndex() error {
	if s.loaded {
		return nil
	}
	s.idx = make(map[string]struct{})
	f, err := os.Open(s.signalsPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			s.loaded = true
			return nil
		}
		return fmt.Errorf("open signals: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var sig schema.Signal
		if err := json.Unmarshal(line, &sig); err != nil {
			return fmt.Errorf("decode signal: %w", err)
		}
		s.idx[sig.ID] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan signals: %w", err)
	}
	s.loaded = true
	return nil
}

func (s *FileStore) Append(sig schema.Signal) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.ensureIndex(); err != nil {
		return err
	}
	if _, ok := s.idx[sig.ID]; ok {
		return nil
	}

	f, err := os.OpenFile(s.signalsPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open signals for append: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(sig); err != nil {
		return fmt.Errorf("encode signal: %w", err)
	}
	s.idx[sig.ID] = struct{}{}
	return nil
}

func (s *FileStore) Has(id string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.ensureIndex(); err != nil {
		return false, err
	}
	_, ok := s.idx[id]
	return ok, nil
}

func (s *FileStore) ListSince(t time.Time) ([]schema.Signal, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.Open(s.signalsPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("open signals: %w", err)
	}
	defer f.Close()

	var out []schema.Signal
	dec := json.NewDecoder(f)
	for {
		var sig schema.Signal
		if err := dec.Decode(&sig); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("decode signal: %w", err)
		}
		if sig.CapturedAt.Before(t) {
			continue
		}
		out = append(out, sig)
	}
	return out, nil
}

func (s *FileStore) Write(fp schema.Fingerprint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return atomicWriteJSON(s.fingerprintPath(), fp)
}

func (s *FileStore) Read() (schema.Fingerprint, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var fp schema.Fingerprint
	raw, err := os.ReadFile(s.fingerprintPath())
	if err != nil {
		return fp, fmt.Errorf("read fingerprint: %w", err)
	}
	if err := json.Unmarshal(raw, &fp); err != nil {
		return fp, fmt.Errorf("decode fingerprint: %w", err)
	}
	return fp, nil
}

func (s *FileStore) Replace(recs []schema.Recommendation) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := s.recommendationsPath()
	tmp, err := os.CreateTemp(s.root, "recommendations-*.tmp")
	if err != nil {
		return fmt.Errorf("create tmp: %w", err)
	}
	tmpName := tmp.Name()
	enc := json.NewEncoder(tmp)
	for _, rec := range recs {
		if err := enc.Encode(rec); err != nil {
			tmp.Close()
			os.Remove(tmpName)
			return fmt.Errorf("encode rec: %w", err)
		}
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("close tmp: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("rename: %w", err)
	}
	return nil
}

func (s *FileStore) List() ([]schema.Recommendation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.Open(s.recommendationsPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("open recommendations: %w", err)
	}
	defer f.Close()

	var out []schema.Recommendation
	dec := json.NewDecoder(f)
	for {
		var rec schema.Recommendation
		if err := dec.Decode(&rec); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("decode rec: %w", err)
		}
		out = append(out, rec)
	}
	return out, nil
}

func atomicWriteJSON(path string, v any) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, filepath.Base(path)+".tmp.*")
	if err != nil {
		return fmt.Errorf("create tmp: %w", err)
	}
	tmpName := tmp.Name()
	enc := json.NewEncoder(tmp)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return fmt.Errorf("encode: %w", err)
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("close tmp: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("rename: %w", err)
	}
	return nil
}

var _ SignalStore = (*FileStore)(nil)
var _ FingerprintStore = (*FileStore)(nil)
var _ RecommendationStore = (*FileStore)(nil)
