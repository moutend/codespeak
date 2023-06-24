package app

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/moutend/codespeak/internal/tts"
)

type App struct {
	HomeDir       string
	audioDir      string
	inputFilePath string
	engine        *tts.Engine
	server        *http.Server
	processOption tts.ProcessOption
	debug         *log.Logger
}

func New() (*App, error) {
	dir, err := os.UserHomeDir()

	if err != nil {
		return nil, err
	}

	homeDir := filepath.Join(dir, ".codespeak")
	inputFilePath := filepath.Join(homeDir, "text", "buffer.txt")
	audioDir := filepath.Join(homeDir, "sound")

	a := &App{
		HomeDir:       homeDir,
		audioDir:      audioDir,
		inputFilePath: inputFilePath,
		debug:         log.New(io.Discard, "", 0),
		engine:        tts.NewEngine(),
		processOption: tts.ProcessOption{
			AudioDir: audioDir,
			Debug:    log.New(io.Discard, "", 0),
		},
	}

	return a, nil
}

func (a *App) SetDebug(logger *log.Logger) {
	if logger == nil {
		return
	}

	a.debug = logger
	a.processOption.Debug = logger
}

func (a *App) SetEnglishVoice(s string) {
	a.processOption.EnglishVoice = s
}

func (a *App) SetEnglishRate(i int) {
	a.processOption.EnglishRate = i
}

func (a *App) SetJapaneseVoice(s string) {
	a.processOption.JapaneseVoice = s
}

func (a *App) SetJapaneseRate(i int) {
	a.processOption.JapaneseRate = i
}
