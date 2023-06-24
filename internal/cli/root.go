package cli

import (
	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Use:   "codespeak",
	Short: "codespeak -- speech synthsis server",
}

func init() {
	RootCommand.PersistentFlags().BoolP("debug", "d", false, "print debug messages")
}
