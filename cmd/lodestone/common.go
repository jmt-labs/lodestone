package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmt-labs/lodestone/internal/config"
	"github.com/jmt-labs/lodestone/internal/lodestone/audit"
	"github.com/jmt-labs/lodestone/internal/lodestone/store"
)

type paths struct {
	repoRoot   string
	storeRoot  string
	configPath string
}

func resolvePaths(rootFlag string) (paths, error) {
	root := rootFlag
	if root == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return paths{}, fmt.Errorf("cwd: %w", err)
		}
		root = cwd
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return paths{}, fmt.Errorf("abs %s: %w", root, err)
	}
	return paths{
		repoRoot:   abs,
		storeRoot:  filepath.Join(abs, ".lodestone"),
		configPath: filepath.Join(abs, config.DefaultFilename),
	}, nil
}

func openStore(p paths) (*store.FileStore, error) {
	return store.New(p.storeRoot)
}

func loadConfig(p paths) (config.Config, error) {
	return config.Load(p.configPath)
}

func openAudit(p paths) (*audit.Log, error) {
	return audit.New(p.storeRoot)
}

func recordAudit(p paths, entry audit.Entry) {
	log, err := openAudit(p)
	if err != nil {
		return
	}
	_ = log.Record(entry)
}
