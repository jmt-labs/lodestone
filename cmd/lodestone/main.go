package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	root := newRootCmd()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	var rootPath string
	cmd := &cobra.Command{
		Use:           "lodestone",
		Short:         "Liest das AI-Ökosystem für dein Repo.",
		Long:          "lodestone sammelt externe AI-Signale, scort sie gegen einen Repo-Fingerprint und erzeugt reproduzierbare Empfehlungen.",
		SilenceUsage:  true,
		SilenceErrors: false,
	}
	cmd.PersistentFlags().StringVar(&rootPath, "root", "", "Projekt-Wurzel (Default: aktuelles Verzeichnis)")

	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newFingerprintCmd(&rootPath))
	cmd.AddCommand(newIngestCmd(&rootPath))
	cmd.AddCommand(newScoreCmd(&rootPath))
	cmd.AddCommand(newSignalsCmd(&rootPath))
	cmd.AddCommand(newPlanCmd(&rootPath))

	for _, verb := range laterPhaseVerbs {
		cmd.AddCommand(newStubCmd(verb.name, verb.short))
	}

	return cmd
}

type verbSpec struct {
	name  string
	short string
}

var laterPhaseVerbs = []verbSpec{
	{"recommend", "Empfehlungen interaktiv durchgehen (Phase 2)"},
	{"apply", "Auto-PR-Engine: Recommendation als PR umsetzen (Phase 4)"},
	{"undo", "Letzten apply-Vorgang rückgängig machen (Phase 4)"},
	{"stats", "Erfolgs-Statistiken angewandter Empfehlungen (Phase 3)"},
	{"calibrate", "Scoring-Gewichte gegen Decision-Log nachjustieren (Phase 3)"},
	{"share", "Decisions anonymisiert teilen (Phase 4)"},
}

func newStubCmd(name, short string) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStderr(), "lodestone %s: not yet implemented in Phase 1\n", name)
			return nil
		},
	}
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Version, Commit und Build-Datum anzeigen",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "lodestone %s (commit %s, built %s)\n", version, commit, date)
		},
	}
}
