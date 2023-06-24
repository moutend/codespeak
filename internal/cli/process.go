package cli

import (
	"log"
	"os"

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
	a, err := app.New()

	if err != nil {
		return err
	}
	if yes, _ := cmd.Flags().GetBool("debug"); yes {
		a.SetDebug(log.New(cmd.ErrOrStderr(), "Debug: ", 0))
	}

	englishVoice, _ := cmd.Flags().GetString("english-voice")
	englishRate, _ := cmd.Flags().GetInt("english-rate")
	japaneseVoice, _ := cmd.Flags().GetString("japanese-voice")
	japaneseRate, _ := cmd.Flags().GetInt("japanese-rate")

	a.SetEnglishVoice(englishVoice)
	a.SetEnglishRate(englishRate)
	a.SetJapaneseVoice(japaneseVoice)
	a.SetJapaneseRate(japaneseRate)

	var input string

	if filePath, _ := cmd.Flags().GetString("file"); filePath != "" {
		data, err := os.ReadFile(filePath)

		if err != nil {
			return err
		}

		input = string(data)
	} else {
		if len(args) < 1 {
			return nil
		}

		input = args[0]
	}
	if err := a.Process(cmd.Context(), input); err != nil {
		return err
	}

	return nil
}

func init() {
	RootCommand.AddCommand(processCommand)

	processCommand.PersistentFlags().StringP("file", "f", "", "read from file")
	processCommand.PersistentFlags().StringP("english-voice", "e", "Alex", "English voice")
	processCommand.PersistentFlags().StringP("japanese-voice", "j", "Kyoko", "Japanese voice")
	processCommand.PersistentFlags().IntP("english-rate", "", 280, "English rate")
	processCommand.PersistentFlags().IntP("japanese-rate", "", 480, "Japanese rate")
}
