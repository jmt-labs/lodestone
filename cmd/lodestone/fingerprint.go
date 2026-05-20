package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jmt-labs/lodestone/internal/lodestone/audit"
	"github.com/jmt-labs/lodestone/internal/lodestone/fingerprint"
)

func newFingerprintCmd(rootPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "fingerprint",
		Short: "Aktuelles Repo analysieren (→ .lodestone/fingerprint.json)",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := resolvePaths(*rootPath)
			if err != nil {
				return err
			}
			cfg, err := loadConfig(p)
			if err != nil {
				return err
			}

			fp, err := fingerprint.New(p.repoRoot).Analyze()
			if err != nil {
				return fmt.Errorf("analyze: %w", err)
			}
			if len(cfg.Goals) > 0 {
				fp.Goals = cfg.Goals
			}
			if len(cfg.TechInterests) > 0 {
				fp.TechInterests = cfg.TechInterests
			}

			s, err := openStore(p)
			if err != nil {
				return err
			}
			if err := s.Write(fp); err != nil {
				return fmt.Errorf("write fingerprint: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "fingerprint written → %s\n", p.storeRoot+"/fingerprint.json")
			fmt.Fprintf(cmd.OutOrStdout(), "  languages : %v\n", fp.Languages)
			fmt.Fprintf(cmd.OutOrStdout(), "  frameworks: %v\n", fp.Frameworks)
			fmt.Fprintf(cmd.OutOrStdout(), "  deps      : %d\n", len(fp.Deps))
			fmt.Fprintf(cmd.OutOrStdout(), "  has_ci    : %v (%s)\n", fp.HasCI, fp.CIProvider)
			recordAudit(p, audit.Entry{
				Verb:    "fingerprint",
				Outcome: "ok",
				Detail:  fmt.Sprintf("languages=%v frameworks=%v deps=%d", fp.Languages, fp.Frameworks, len(fp.Deps)),
			})
			return nil
		},
	}
}
