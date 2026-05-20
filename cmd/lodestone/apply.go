package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jmt-labs/lodestone/internal/lodestone/apply"
	"github.com/jmt-labs/lodestone/internal/lodestone/audit"
	"github.com/jmt-labs/lodestone/internal/lodestone/planning"
)

func newApplyCmd(rootPath *string) *cobra.Command {
	var model string
	cmd := &cobra.Command{
		Use:   "apply <rec-id>",
		Short: "Recommendation als Draft-PR aufsetzen (harte Safety-Gates)",
		Args:  cobra.ExactArgs(1),
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
				return fmt.Errorf("read fingerprint: %w", err)
			}
			recs, err := s.List()
			if err != nil {
				return err
			}
			rec, ok := findRec(recs, recID)
			if !ok {
				return fmt.Errorf("recommendation %q not found", recID)
			}

			state, err := apply.NewState(p.storeRoot)
			if err != nil {
				return err
			}
			planOpts := []planning.Option{}
			if model != "" {
				planOpts = append(planOpts, planning.WithModel(model))
			}
			eng := apply.New(p.repoRoot, state,
				apply.WithPlanning(planning.New(planOpts...)),
			)
			res, err := eng.Apply(context.Background(), fp, rec)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "applied: branch=%s pr=%d url=%s status=%s\n",
				res.Apply.Branch, res.Apply.PRNumber, res.Apply.PRURL, res.Apply.Status)
			recordAudit(p, audit.Entry{
				Verb:    "apply",
				Args:    map[string]string{"rec_id": rec.ID, "branch": res.Apply.Branch},
				Outcome: res.Apply.Status,
				Detail:  res.Apply.PRURL,
			})
			return nil
		},
	}
	cmd.Flags().StringVar(&model, "model", "", "Claude-Modell-Override")
	return cmd
}

func newUndoCmd(rootPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "undo <branch-or-rec-id>",
		Short: "Letzten apply rückgängig machen (PR schließen + Branch löschen)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := resolvePaths(*rootPath)
			if err != nil {
				return err
			}
			state, err := apply.NewState(p.storeRoot)
			if err != nil {
				return err
			}
			eng := apply.New(p.repoRoot, state)
			res, err := eng.Undo(context.Background(), args[0])
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "undone: rec=%s branch=%s status=%s\n", res.RecID, res.Branch, res.Status)
			recordAudit(p, audit.Entry{
				Verb:    "undo",
				Args:    map[string]string{"branch": res.Branch},
				Outcome: res.Status,
			})
			return nil
		},
	}
}

func newStatsCmd(rootPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "Apply-Statistiken aus .lodestone/applies.jsonl",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := resolvePaths(*rootPath)
			if err != nil {
				return err
			}
			state, err := apply.NewState(p.storeRoot)
			if err != nil {
				return err
			}
			list, err := state.List()
			if err != nil {
				return err
			}
			counts := map[string]int{}
			for _, a := range list {
				counts[a.Status]++
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Apply-Stats über %d Einträge:\n", len(list))
			for status, n := range counts {
				fmt.Fprintf(cmd.OutOrStdout(), "  %-22s %d\n", status, n)
			}
			return nil
		},
	}
}
