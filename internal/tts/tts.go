package tts

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/moutend/codespeak/internal/token"
)

type ProcessOption struct {
	Debug         *log.Logger
	AudioDir      string
	EnglishVoice  string
	EnglishRate   int
	JapaneseVoice string
	JapaneseRate  int
}

func Process(ctx context.Context, ts []token.Token, option ProcessOption) error {
	var cmd *exec.Cmd

	englishVoice := option.EnglishVoice
	englishRate := fmt.Sprint(option.EnglishRate)
	japaneseVoice := option.JapaneseVoice
	japaneseRate := fmt.Sprint(option.JapaneseRate)

	for _, t := range ts {
		switch t.Kind {
		case token.Number, token.Symbol:
			args := []string{}

			for _, r := range []rune(t.Text) {
				args = append(args, filepath.Join(option.AudioDir, fmt.Sprintf("%03d.wav", r+1)))
			}
			if len(args) == 0 {
				continue
			}

			cmd = exec.CommandContext(ctx, "play", args...)
		case token.Alphabet:
			cmd = exec.CommandContext(ctx, "say", "-v", englishVoice, "-r", englishRate, fmt.Sprintf("%q", t.Text))
		case token.Unicode:
			cmd = exec.CommandContext(ctx, "say", "-v", japaneseVoice, "-r", japaneseRate, fmt.Sprintf("%q", t.Text))
		}

		go option.Debug.Println("Processing:", cmd.String())

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

func ProcessOffline(ctx context.Context, ts []token.Token, option ProcessOption) error {
	commands := make([]*exec.Cmd, 0, len(ts))

	englishVoice := option.EnglishVoice
	englishRate := fmt.Sprint(option.EnglishRate)
	japaneseVoice := option.JapaneseVoice
	japaneseRate := fmt.Sprint(option.JapaneseRate)

	temporaryDir := filepath.Join(option.AudioDir, "temporary")

	os.RemoveAll(temporaryDir)
	os.MkdirAll(temporaryDir, 0755)

	filePaths := []string{}

	for i, t := range ts {
		switch t.Kind {
		case token.Number, token.Symbol:
			for _, r := range []rune(t.Text) {
				filePaths = append(filePaths, filepath.Join(option.AudioDir, fmt.Sprintf("%03d.wav", r+1)))
			}
		case token.Alphabet:
			filePath := filepath.Join(temporaryDir, fmt.Sprintf("text%03d.aiff", i))
			filePaths = append(filePaths, filePath)
			commands = append(commands, exec.CommandContext(ctx, "say", "-o", filePath, "-v", englishVoice, "-r", englishRate, fmt.Sprintf("%q", t.Text)))
		case token.Unicode:
			filePath := filepath.Join(temporaryDir, fmt.Sprintf("text%03d.aiff", i))
			filePaths = append(filePaths, filePath)
			commands = append(commands, exec.CommandContext(ctx, "say", "-o", filePath, "-v", japaneseVoice, "-r", japaneseRate, fmt.Sprintf("%q", t.Text)))
		}
	}

	var wg sync.WaitGroup

	for i := range commands {
		wg.Add(1)

		go func(n int) {
			if err := commands[n].Run(); err != nil {
				option.Debug.Println("Error:", err)
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	if err := exec.CommandContext(ctx, "play", filePaths...).Run(); err != nil {
		return err
	}

	return nil
}

type Engine struct {
	cancelFuncCounter int
	cancelFuncMap     map[int]context.CancelFunc
	cancelFuncMutex   sync.Mutex
}

func (e *Engine) Speak(ctx context.Context, tokens []token.Token, option ProcessOption) {
	e.cancelFuncMutex.Lock()
	defer e.cancelFuncMutex.Unlock()

	for k := range e.cancelFuncMap {
		if cancel, ok := e.cancelFuncMap[k]; ok {
			cancel()
		}

		delete(e.cancelFuncMap, k)
	}

	ctx, cancel := context.WithCancel(ctx)

	e.cancelFuncCounter += 1
	e.cancelFuncMap[e.cancelFuncCounter] = cancel

	go option.Debug.Println("Speak function invoked:", e.cancelFuncCounter)
	go ProcessOffline(ctx, tokens, option)
}

func (e *Engine) Pause(option ProcessOption) {
	e.cancelFuncMutex.Lock()
	defer e.cancelFuncMutex.Unlock()

	for k := range e.cancelFuncMap {
		if cancel, ok := e.cancelFuncMap[k]; ok {
			cancel()
		}

		delete(e.cancelFuncMap, k)
	}

	go option.Debug.Println("Pause function invoked:", e.cancelFuncCounter)
}

func (e *Engine) Close() {
}

func NewEngine() *Engine {
	return &Engine{
		cancelFuncMap: map[int]context.CancelFunc{},
	}
}
