package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/jmt-labs/lodestone/internal/lodestone/audit"
	"github.com/jmt-labs/lodestone/internal/lodestone/memory"
)

func newMemoryCmd(rootPath *string) *cobra.Command {
	var (
		days    int
		memPath string
	)
	cmd := &cobra.Command{
		Use:   "memory",
		Short: "Decisions aus .lodestone/decisions.log nach .claude/memory.json konsolidieren",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := resolvePaths(*rootPath)
			if err != nil {
				return err
			}
			decisionsPath := filepath.Join(p.storeRoot, audit.DefaultFilename)
			out := memPath
			if out == "" {
				out = filepath.Join(p.repoRoot, memory.DefaultRelPath)
			}
			since := time.Now().UTC().AddDate(0, 0, -days)
			added, err := memory.Consolidate(decisionsPath, out, memory.WithSince(since))
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "memory aktualisiert → %s (+%d Einträge, Fenster %d Tage)\n", out, added, days)
			return nil
		},
	}
	cmd.Flags().IntVar(&days, "days", 90, "wieviele Tage rückwärts berücksichtigen")
	cmd.Flags().StringVar(&memPath, "out", "", "Ziel-Pfad (Default: .claude/memory.json relativ zum --root)")
	return cmd
}
