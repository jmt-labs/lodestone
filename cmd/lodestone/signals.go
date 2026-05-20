package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/spf13/cobra"

	"github.com/jmt-labs/lodestone/internal/lodestone/schema"
)

func newSignalsCmd(rootPath *string) *cobra.Command {
	var (
		since     string
		source    string
		top       int
		asJSON    bool
	)
	cmd := &cobra.Command{
		Use:   "signals",
		Short: "Gespeicherte Signale anzeigen / filtern",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := resolvePaths(*rootPath)
			if err != nil {
				return err
			}
			s, err := openStore(p)
			if err != nil {
				return err
			}

			var cutoff time.Time
			if since != "" {
				cutoff, err = time.Parse("2006-01-02", since)
				if err != nil {
					return fmt.Errorf("--since must be YYYY-MM-DD: %w", err)
				}
			}

			signals, err := s.ListSince(cutoff)
			if err != nil {
				return err
			}
			if source != "" {
				signals = filterBySource(signals, source)
			}

			sort.Slice(signals, func(i, j int) bool {
				if signals[i].Stars != signals[j].Stars {
					return signals[i].Stars > signals[j].Stars
				}
				return signals[i].ID < signals[j].ID
			})

			if top > 0 && len(signals) > top {
				signals = signals[:top]
			}

			if asJSON {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(signals)
			}
			renderSignals(cmd, signals)
			return nil
		},
	}
	cmd.Flags().StringVar(&since, "since", "", "nur Signale ab Datum (YYYY-MM-DD)")
	cmd.Flags().StringVar(&source, "source", "", "nur Signale dieser Quelle")
	cmd.Flags().IntVar(&top, "top", 0, "auf die Top-N nach Stars beschränken (0 = alle)")
	cmd.Flags().BoolVar(&asJSON, "json", false, "JSON-Ausgabe statt Tabelle")
	return cmd
}

func filterBySource(sigs []schema.Signal, source string) []schema.Signal {
	out := make([]schema.Signal, 0, len(sigs))
	for _, s := range sigs {
		if s.Source == source {
			out = append(out, s)
		}
	}
	return out
}

func renderSignals(cmd *cobra.Command, sigs []schema.Signal) {
	if len(sigs) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "(keine Signale)")
		return
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%-6s  %-18s  %-8s  %s\n", "STARS", "SOURCE", "LANG", "TITLE")
	for _, s := range sigs {
		title := s.Title
		if len(title) > 60 {
			title = title[:57] + "..."
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%-6d  %-18s  %-8s  %s\n", s.Stars, s.Source, s.Language, title)
	}
}
