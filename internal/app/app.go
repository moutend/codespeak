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
		inputFilePath: inputFilePath,
		audioDir:      audioDir,
		engine: tts.NewEngine(tts.ProcessOption{
			Debug:    log.New(io.Discard, "", 0),
			AudioDir: audioDir,
		}),
		debug: log.New(io.Discard, "", 0),
	}

	return a, nil
}

func (a *App) SetDebug(logger *log.Logger) {
	if logger == nil {
		return
	}
	if a.engine != nil {
		a.engine.ProcessOption.Debug = logger
	}
	a.debug = logger
}
