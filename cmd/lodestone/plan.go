package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jmt-labs/lodestone/internal/lodestone/audit"
	"github.com/jmt-labs/lodestone/internal/lodestone/planning"
	"github.com/jmt-labs/lodestone/internal/lodestone/schema"
)

func newPlanCmd(rootPath *string) *cobra.Command {
	var (
		model  string
		dryRun bool
	)
	cmd := &cobra.Command{
		Use:   "plan <rec-id>",
		Short: "Spec/Plan/Tasks aus Recommendation generieren (ruft Claude)",
		Long: `Lädt die Recommendation per ID und den Fingerprint aus .lodestone/,
ruft das Claude-CLI mit einem strukturierten Prompt und persistiert die
Antwort als Spec + Plan unter docs/superpowers/.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			recID := args[0]
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
			recs, err := s.List()
			if err != nil {
				return err
			}
			rec, ok := findRec(recs, recID)
			if !ok {
				return fmt.Errorf("recommendation %q not found (run `lodestone score`)", recID)
			}

			if model == "" {
				model = planning.DefaultModel
			}

			engine := planning.New(planning.WithModel(model))
			if dryRun {
				prompt, err := planning.BuildPrompt(fp, rec)
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), prompt)
				return nil
			}

			res, err := engine.Plan(context.Background(), fp, rec)
			if err != nil {
				return err
			}
			if err := res.Persist(p.repoRoot); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "spec written → %s\nplan written → %s\nmodel: %s\n", res.SpecPath, res.PlanPath, res.Model)
			recordAudit(p, audit.Entry{
				Verb:    "plan",
				Args:    map[string]string{"rec_id": recID, "model": res.Model},
				Outcome: "ok",
				Detail:  fmt.Sprintf("spec=%s plan=%s", res.SpecPath, res.PlanPath),
			})
			return nil
		},
	}
	cmd.Flags().StringVar(&model, "model", "", "Claude-Modell-Override (Default: claude-opus-4-7)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "nur Prompt anzeigen, keinen Claude-Aufruf machen")
	return cmd
}

func findRec(recs []schema.Recommendation, id string) (schema.Recommendation, bool) {
	for _, r := range recs {
		if r.ID == id || r.SignalID == id {
			return r, true
		}
	}
	return schema.Recommendation{}, false
}
