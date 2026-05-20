package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jmt-labs/lodestone/internal/config"
	"github.com/jmt-labs/lodestone/internal/lodestone/ingest"
	"github.com/jmt-labs/lodestone/internal/lodestone/schema"
	"github.com/jmt-labs/lodestone/internal/lodestone/store"
)

const (
	sourceGitHubTrending     = "github_trending"
	sourceHackerNews         = "hackernews"
	sourceArXiv              = "arxiv"
	sourceAnthropicChangelog = "anthropic_changelog"
	sourceOpenAIChangelog    = "openai_changelog"
	sourceNPMTrending        = "npm_trending"
)

const mockFixturesEnv = "LODESTONE_MOCK_FIXTURES"

func newIngestCmd(rootPath *string) *cobra.Command {
	var (
		sources []string
		useMock bool
	)
	cmd := &cobra.Command{
		Use:   "ingest",
		Short: "Externe Signale abrufen (→ .lodestone/signals.jsonl)",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := resolvePaths(*rootPath)
			if err != nil {
				return err
			}
			cfg, err := loadConfig(p)
			if err != nil {
				return err
			}

			if len(sources) == 0 {
				sources = knownSources()
			}

			s, err := openStore(p)
			if err != nil {
				return err
			}

			if useMock {
				return runMockIngest(cmd, sources, s)
			}

			ctx := context.Background()
			cacheDir := filepath.Join(p.storeRoot, "cache")

			var totalFetched, totalNew int
			for _, name := range sources {
				src, err := buildSource(name, cfg, cacheDir)
				if err != nil {
					return err
				}
				signals, err := src.Fetch(ctx)
				if err != nil {
					return fmt.Errorf("fetch %s: %w", name, err)
				}
				added, err := appendNew(s, signals)
				if err != nil {
					return err
				}
				totalFetched += len(signals)
				totalNew += added
				fmt.Fprintf(cmd.OutOrStdout(), "ingest %s: %d fetched, %d new\n", name, len(signals), added)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "total: %d fetched, %d new\n", totalFetched, totalNew)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&sources, "source", nil, "source name (kann mehrfach angegeben werden); Default: alle Phase-1-Quellen")
	cmd.Flags().BoolVar(&useMock, "mock", false, "Mock-Modus: Signale aus $"+mockFixturesEnv+" laden statt HTTP-Fetch")
	return cmd
}

func buildSource(name string, cfg config.Config, cacheDir string) (ingest.Source, error) {
	switch name {
	case sourceGitHubTrending:
		return ingest.NewGithubTrending(
			ingest.WithGithubCacheDir(cacheDir),
			ingest.WithGithubMinStars(cfg.Lodestone.MinStars),
			ingest.WithGithubRecentDays(cfg.Lodestone.MinAgeDays),
		), nil
	case sourceHackerNews:
		return ingest.NewHackerNews(
			ingest.WithHackerNewsCacheDir(cacheDir),
		), nil
	case sourceArXiv:
		return ingest.NewArXiv(
			ingest.WithArXivCacheDir(cacheDir),
		), nil
	case sourceAnthropicChangelog:
		return ingest.NewAnthropicChangelog(
			ingest.WithChangelogCacheDir(cacheDir),
		), nil
	case sourceOpenAIChangelog:
		return ingest.NewOpenAIChangelog(
			ingest.WithChangelogCacheDir(cacheDir),
		), nil
	case sourceNPMTrending:
		return ingest.NewNPMTrending(
			ingest.WithNPMCacheDir(cacheDir),
		), nil
	default:
		return nil, fmt.Errorf("unknown source %q; valid: %s", name, strings.Join(knownSources(), ", "))
	}
}

func knownSources() []string {
	return []string{
		sourceGitHubTrending,
		sourceHackerNews,
		sourceArXiv,
		sourceAnthropicChangelog,
		sourceOpenAIChangelog,
		sourceNPMTrending,
	}
}

func appendNew(s *store.FileStore, signals []schema.Signal) (int, error) {
	added := 0
	for _, sig := range signals {
		ok, err := s.Has(sig.ID)
		if err != nil {
			return added, err
		}
		if ok {
			continue
		}
		if err := s.Append(sig); err != nil {
			return added, fmt.Errorf("append %s: %w", sig.ID, err)
		}
		added++
	}
	return added, nil
}

func runMockIngest(cmd *cobra.Command, sources []string, s *store.FileStore) error {
	dir := os.Getenv(mockFixturesEnv)
	if dir == "" {
		return fmt.Errorf("--mock requires $%s to be set to a directory containing <source>.json fixtures", mockFixturesEnv)
	}
	var totalFetched, totalNew int
	for _, name := range sources {
		path := filepath.Join(dir, name+".json")
		raw, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read fixture %s: %w", path, err)
		}
		var signals []schema.Signal
		if err := json.Unmarshal(raw, &signals); err != nil {
			return fmt.Errorf("parse fixture %s: %w", path, err)
		}
		added, err := appendNew(s, signals)
		if err != nil {
			return err
		}
		totalFetched += len(signals)
		totalNew += added
		fmt.Fprintf(cmd.OutOrStdout(), "ingest %s (mock): %d fetched, %d new\n", name, len(signals), added)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "total: %d fetched, %d new\n", totalFetched, totalNew)
	return nil
}
