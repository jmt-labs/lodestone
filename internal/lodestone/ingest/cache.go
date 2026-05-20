package ingest

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jmt-labs/lodestone/internal/lodestone/schema"
)

func cachePath(cacheDir, source string, now time.Time) string {
	if cacheDir == "" {
		return ""
	}
	return filepath.Join(cacheDir, fmt.Sprintf("%s-%s.json", source, now.Format("2006-01-02")))
}

func loadCache(path string) ([]schema.Signal, bool, error) {
	if path == "" {
		return nil, false, nil
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, false, nil
		}
		return nil, false, err
	}
	var sigs []schema.Signal
	if err := json.Unmarshal(raw, &sigs); err != nil {
		return nil, false, err
	}
	return sigs, true, nil
}

func saveCache(path string, sigs []schema.Signal) error {
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(sigs, "", "  ")
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
