package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/jmt-labs/lodestone/internal/lodestone/audit"
	"github.com/jmt-labs/lodestone/internal/lodestone/scoring"
)

func newScoreCmd(rootPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "score",
		Short: "Signale × Fingerprint scoren (→ .lodestone/recommendations.jsonl)",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := resolvePaths(*rootPath)
			if err != nil {
				return err
			}
			s, err := openStore(p)
			if err != nil {
				return err
			}

			fp, err := s.Read()
			if err != nil {
				return fmt.Errorf("read fingerprint (run `lodestone fingerprint` first): %w", err)
			}
			signals, err := s.ListSince(time.Time{})
			if err != nil {
				return err
			}
			if len(signals) == 0 {
				return fmt.Errorf("no signals in store (run `lodestone ingest` first)")
			}

			recs, err := scoring.Score(fp, signals)
			if err != nil {
				return fmt.Errorf("score: %w", err)
			}
			if err := s.Replace(recs); err != nil {
				return fmt.Errorf("replace recommendations: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "scored %d signals → %d recommendations\n", len(signals), len(recs))
			recordAudit(p, audit.Entry{
				Verb:    "score",
				Outcome: "ok",
				Detail:  fmt.Sprintf("signals=%d recommendations=%d", len(signals), len(recs)),
			})
			return nil
		},
	}
}
