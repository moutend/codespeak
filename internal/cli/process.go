package cli

import (
	"log"

	"github.com/moutend/codespeak/internal/app"
	"github.com/spf13/cobra"
)

var processCommand = &cobra.Command{
	Use:     "process",
	Aliases: []string{"p"},
	Short:   "process the description of the given character",
	RunE:    processCommandRunE,
}

func processCommandRunE(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return nil
	}

	a, err := app.New()

	if err != nil {
		return err
	}
	if yes, _ := cmd.Flags().GetBool("debug"); yes {
		a.SetDebug(log.New(cmd.ErrOrStderr(), "Debug: ", 0))
	}
	if err := a.Process(cmd.Context(), args[0]); err != nil {
		return err
	}

	return nil
}

func init() {
	RootCommand.AddCommand(processCommand)

	processCommand.PersistentFlags().StringP("file", "f", "", "read from file")
	processCommand.PersistentFlags().StringSliceP("languages", "l", []string{}, "target languages (e.g. 'ja') (default: empty)")
}
