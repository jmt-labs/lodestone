package apply

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

const StateFilename = "applies.jsonl"

type Apply struct {
	RecID     string    `json:"rec_id"`
	Branch    string    `json:"branch"`
	PRNumber  int       `json:"pr_number,omitempty"`
	PRURL     string    `json:"pr_url,omitempty"`
	Status    string    `json:"status"`
	AppliedAt time.Time `json:"applied_at"`
}

type State struct {
	path string
}

func NewState(root string) (*State, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, fmt.Errorf("mkdir %s: %w", root, err)
	}
	return &State{path: filepath.Join(root, StateFilename)}, nil
}

func (s *State) Path() string { return s.path }

func (s *State) List() ([]Apply, error) {
	f, err := os.Open(s.path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("open %s: %w", s.path, err)
	}
	defer f.Close()
	var out []Apply
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		var a Apply
		if err := json.Unmarshal(scanner.Bytes(), &a); err != nil {
			return nil, fmt.Errorf("decode apply: %w", err)
		}
		out = append(out, a)
	}
	return out, scanner.Err()
}

func (s *State) Append(a Apply) error {
	if a.AppliedAt.IsZero() {
		a.AppliedAt = time.Now().UTC()
	}
	raw, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf("marshal apply: %w", err)
	}
	f, err := os.OpenFile(s.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open %s: %w", s.path, err)
	}
	defer f.Close()
	_, err = f.Write(append(raw, '\n'))
	return err
}

func (s *State) FindBy(branchOrRec string) (Apply, bool, error) {
	list, err := s.List()
	if err != nil {
		return Apply{}, false, err
	}
	for _, a := range list {
		if a.Branch == branchOrRec || a.RecID == branchOrRec {
			return a, true, nil
		}
	}
	return Apply{}, false, nil
}

func (s *State) Replace(updated Apply) error {
	list, err := s.List()
	if err != nil {
		return err
	}
	for i, a := range list {
		if a.Branch == updated.Branch {
			list[i] = updated
		}
	}
	tmp, err := os.CreateTemp(filepath.Dir(s.path), filepath.Base(s.path)+".tmp.*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	enc := json.NewEncoder(tmp)
	for _, a := range list {
		if err := enc.Encode(a); err != nil {
			tmp.Close()
			os.Remove(tmpName)
			return err
		}
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}
	return os.Rename(tmpName, s.path)
}
