package cli

import (
	"github.com/moutend/codespeak/internal/app"
	"github.com/spf13/cobra"
)

var versionCommand = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "print version",
	RunE:    versionCommandRunE,
}

func versionCommandRunE(cmd *cobra.Command, args []string) error {
	cmd.Println(app.Version())

	return nil
}

func init() {
	RootCommand.AddCommand(versionCommand)
}
