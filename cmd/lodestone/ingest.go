package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jmt-labs/lodestone/internal/config"
	"github.com/jmt-labs/lodestone/internal/lodestone/ingest"
)

const (
	sourceGitHubTrending = "github_trending"
	sourceHackerNews     = "hackernews"
)

func newIngestCmd(rootPath *string) *cobra.Command {
	var sources []string
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
				sources = []string{sourceGitHubTrending, sourceHackerNews}
			}

			s, err := openStore(p)
			if err != nil {
				return err
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
				added := 0
				for _, sig := range signals {
					ok, err := s.Has(sig.ID)
					if err != nil {
						return err
					}
					if ok {
						continue
					}
					if err := s.Append(sig); err != nil {
						return fmt.Errorf("append %s: %w", sig.ID, err)
					}
					added++
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
	default:
		return nil, fmt.Errorf("unknown source %q; valid: %s", name, strings.Join(knownSources(), ", "))
	}
}

func knownSources() []string {
	return []string{sourceGitHubTrending, sourceHackerNews}
}
