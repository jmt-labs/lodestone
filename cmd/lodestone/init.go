package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jmt-labs/lodestone/internal/config"
	"github.com/jmt-labs/lodestone/internal/lodestone/skills"
)

const gitignoreSnippet = `
# Lodestone (Phase 1+) — lokale Artefakte ignorieren, decisions.log ausnehmen
.lodestone/
!.lodestone/decisions.log
`

const defaultConfigYAML = `# Lodestone-Konfiguration. Alle Felder sind optional.
#
# Goals und tech_interests fließen in den Fingerprint ein und beeinflussen
# das Scoring + die Planning-Engine ab Phase 2.

goals: []
tech_interests: []

lodestone:
  min_stars: 50
  min_age_days: 30
  max_last_commit_age_days: 180
  require_license: true
`

func newInitCmd(rootPath *string) *cobra.Command {
	var (
		writeConfig    bool
		writeGitignore bool
		writeSkills    bool
		skillsDir      string
		force          bool
	)
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Lodestone in einem Repo einrichten (Config, .gitignore, Skills)",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := resolvePaths(*rootPath)
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()

			if writeConfig {
				if err := writeConfigYAML(p.configPath, force, out); err != nil {
					return err
				}
			}
			if writeGitignore {
				if err := appendGitignore(filepath.Join(p.repoRoot, ".gitignore"), out); err != nil {
					return err
				}
			}
			if writeSkills {
				target := skillsDir
				if target == "" {
					target = filepath.Join(p.repoRoot, ".claude", "skills")
				}
				if err := installSkills(target, force, out); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&writeConfig, "config", true, "lege .lodestone.yaml mit Defaults an")
	cmd.Flags().BoolVar(&writeGitignore, "gitignore", true, "hänge .gitignore-Snippet für .lodestone/ an")
	cmd.Flags().BoolVar(&writeSkills, "skills", true, "installiere die vier Lodestone-Skills nach .claude/skills/")
	cmd.Flags().StringVar(&skillsDir, "skills-dir", "", "Ziel-Verzeichnis für Skills (Default: .claude/skills)")
	cmd.Flags().BoolVar(&force, "force", false, "vorhandene Dateien überschreiben")
	return cmd
}

func writeConfigYAML(path string, force bool, out io.Writer) error {
	return writeFileIfMissing(path, defaultConfigYAML, force, out, "config")
}

func writeFileIfMissing(path, content string, force bool, _ io.Writer, label string) error {
	if !force {
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("[%s] %s existiert bereits, skip (--force für overwrite)\n", label, path)
			return nil
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	fmt.Printf("[%s] geschrieben: %s\n", label, path)
	return nil
}

func appendGitignore(path string, _ io.Writer) error {
	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read .gitignore: %w", err)
	}
	if strings.Contains(string(existing), ".lodestone/") {
		fmt.Printf("[gitignore] %s enthält bereits .lodestone/-Eintrag, skip\n", path)
		return nil
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open .gitignore: %w", err)
	}
	defer f.Close()
	if len(existing) > 0 && !strings.HasSuffix(string(existing), "\n") {
		if _, err := f.WriteString("\n"); err != nil {
			return err
		}
	}
	if _, err := f.WriteString(gitignoreSnippet); err != nil {
		return err
	}
	fmt.Printf("[gitignore] Snippet angehängt an %s\n", path)
	return nil
}

func installSkills(dir string, force bool, _ io.Writer) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}
	names, err := skills.List()
	if err != nil {
		return fmt.Errorf("list embedded skills: %w", err)
	}
	for _, name := range names {
		if !strings.HasSuffix(name, ".md") {
			continue
		}
		raw, err := skills.Read(name)
		if err != nil {
			return fmt.Errorf("read embedded %s: %w", name, err)
		}
		target := filepath.Join(dir, name)
		if !force {
			if _, err := os.Stat(target); err == nil {
				fmt.Printf("[skills] %s existiert, skip (--force für overwrite)\n", target)
				continue
			}
		}
		if err := os.WriteFile(target, raw, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", target, err)
		}
		fmt.Printf("[skills] installiert: %s\n", target)
	}
	return nil
}

var _ = config.DefaultFilename
