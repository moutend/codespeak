package cli

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/moutend/codespeak/internal/app"
	"github.com/spf13/cobra"
)

var serveCommand = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"s"},
	Short:   "start the codespeak server",
	RunE:    serveCommandRunE,
}

func serveCommandRunE(cmd *cobra.Command, args []string) error {
	a, err := app.New()

	if err != nil {
		return err
	}
	if yes, _ := cmd.Flags().GetBool("debug"); yes {
		a.SetDebug(log.New(cmd.ErrOrStderr(), "Debug: ", 0))
	}

	host, _ := cmd.Flags().GetString("host")
	port, _ := cmd.Flags().GetInt("port")

	u, err := url.Parse(fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		return err
	}

	englishVoice, _ := cmd.Flags().GetString("english-voice")
	englishRate, _ := cmd.Flags().GetInt("english-rate")
	japaneseVoice, _ := cmd.Flags().GetString("japanese-voice")
	japaneseRate, _ := cmd.Flags().GetInt("japanese-rate")

	a.SetEnglishVoice(englishVoice)
	a.SetEnglishRate(englishRate)
	a.SetJapaneseVoice(japaneseVoice)
	a.SetJapaneseRate(japaneseRate)

	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	if err := a.Start(ctx, u); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	if err := a.Stop(ctx); err != nil {
		return err
	}

	return nil
}

func init() {
	RootCommand.AddCommand(serveCommand)

	serveCommand.PersistentFlags().StringP("host", "", "localhost", "host name")
	serveCommand.PersistentFlags().IntP("port", "", 8501, "port")
	serveCommand.PersistentFlags().StringP("english-voice", "e", "Alex", "English voice")
	serveCommand.PersistentFlags().StringP("japanese-voice", "j", "Kyoko", "Japanese voice")
	serveCommand.PersistentFlags().IntP("english-rate", "", 280, "English rate")
	serveCommand.PersistentFlags().IntP("japanese-rate", "", 480, "Japanese rate")
}
