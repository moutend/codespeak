package tts

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/moutend/codespeak/pkg/global"
	"github.com/moutend/codespeak/pkg/token"
)

type TTS struct {
	ctx        context.Context
	cancel     context.CancelFunc
	api        sync.Mutex
	wg         sync.WaitGroup
	quitChan   chan struct{}
	tokensChan chan struct {
		ctx    context.Context
		tokens []token.Token
	}
}

func sayAlex(ctx context.Context, wg *sync.WaitGroup, text, outputPath string) {
	defer wg.Done()

	if err := exec.CommandContext(ctx, "say", "-v", "Alex", "-r", "272", fmt.Sprintf("%q", "[[ pbas 42 ]]"+text), "-o", outputPath).Run(); err != nil {
		go log.Println("tts: failed to generate audio:", outputPath)

		return
	}

	stat, err := os.Stat(outputPath)

	if err != nil {
		go log.Println("tts: failed to stat file:", outputPath)

		os.Remove(outputPath)

		return
	}
	if stat.Size() <= 4096 {
		go log.Println("tts: generated file is broken:", outputPath)

		os.Remove(outputPath)

		return
	}

	return
}
func sayKyoko(ctx context.Context, wg *sync.WaitGroup, text, outputPath string) {
	defer wg.Done()

	if err := exec.CommandContext(ctx, "say", "-v", "Kyoko", "-r", "480", fmt.Sprintf("%q", text), "-o", outputPath).Run(); err != nil {
		go log.Println("tts: failed to generate audio:", outputPath)

		return
	}

	stat, err := os.Stat(outputPath)

	if err != nil {
		go log.Println("tts: failed to stat file: ", outputPath)

		os.Remove(outputPath)

		return
	}
	if stat.Size() <= 4096 {
		go log.Println("tts: generated file is broken:", outputPath)

		os.Remove(outputPath)

		return
	}

	return
}

func process(ctx context.Context, ts []token.Token) error {
	tempPath := filepath.Join(global.CodespeakAudioPath, "tmp")
	os.MkdirAll(tempPath, 0755)

	var wg sync.WaitGroup

	removePaths := []string{}
	audioPaths := []string{}

	for _, t := range ts {
		switch t.Kind {
		case token.Number, token.Symbol:
			for _, r := range []rune(t.Text) {
				audioPaths = append(audioPaths, filepath.Join(global.CodespeakAudioPath, fmt.Sprintf("%03d.wav", r+1)))
			}
		case token.Alphabet:
			hash := make([]byte, 16, 16)

			if _, err := rand.Read(hash); err != nil {
				continue
			}

			aiffPath := filepath.Join(tempPath, fmt.Sprintf("%s_voice.aiff", hex.EncodeToString(hash)))
			removePaths = append(removePaths, aiffPath)
			audioPaths = append(audioPaths, aiffPath)
			audioPaths = append(audioPaths, filepath.Join(global.CodespeakAudioPath, `mute.wav`))
			audioPaths = append(audioPaths, filepath.Join(global.CodespeakAudioPath, `mute.wav`))
			audioPaths = append(audioPaths, filepath.Join(global.CodespeakAudioPath, `mute.wav`))

			wg.Add(1)

			go sayAlex(ctx, &wg, t.Text, aiffPath)
		case token.Unicode:
			hash := make([]byte, 16, 16)

			if _, err := rand.Read(hash); err != nil {
				continue
			}

			aiffPath := filepath.Join(tempPath, fmt.Sprintf("%s_voice.aiff", hex.EncodeToString(hash)))
			removePaths = append(removePaths, aiffPath)
			audioPaths = append(audioPaths, aiffPath)
			audioPaths = append(audioPaths, filepath.Join(global.CodespeakAudioPath, `mute.wav`))
			audioPaths = append(audioPaths, filepath.Join(global.CodespeakAudioPath, `mute.wav`))
			audioPaths = append(audioPaths, filepath.Join(global.CodespeakAudioPath, `mute.wav`))

			wg.Add(1)

			go sayKyoko(ctx, &wg, t.Text, aiffPath)
		}
	}

	wg.Wait()

	defer func() {
		for _, removePath := range removePaths {
			os.Remove(removePath)
		}
	}()
	for i, audioPath := range audioPaths {
		if !strings.HasSuffix(audioPath, "_voice.aiff") {
			continue
		}
		if _, err := os.Stat(audioPath); err != nil {
			audioPaths = append(audioPaths[:i], audioPaths[i+1:]...)
		}
	}
	if err := exec.CommandContext(ctx, "play", audioPaths...).Run(); err != nil {
		return err
	}

	return nil
}

func (t *TTS) SpeakContext(ctx context.Context, tokens []token.Token) {
	t.api.Lock()
	defer t.api.Unlock()

	if t.cancel != nil {
		t.cancel()
	}
	if t.tokensChan != nil {
		t.ctx, t.cancel = context.WithCancel(ctx)

		t.tokensChan <- struct {
			ctx    context.Context
			tokens []token.Token
		}{
			ctx:    t.ctx,
			tokens: tokens,
		}
	}

	return
}

func (t *TTS) Pause() {
	t.api.Lock()
	defer t.api.Unlock()

	if t.cancel != nil {
		t.cancel()
	}

	return
}

func (t *TTS) Close() {
	t.api.Lock()
	defer t.api.Unlock()

	t.quitChan <- struct{}{}

	t.wg.Wait()

	close(t.quitChan)
	close(t.tokensChan)

	return
}

func Open() *TTS {
	t := &TTS{
		quitChan: make(chan struct{}),
		tokensChan: make(chan struct {
			ctx    context.Context
			tokens []token.Token
		}),
	}

	t.wg.Add(1)

	go func() {
		defer t.wg.Done()

		for {
			select {
			case <-t.quitChan:
				return
			case v := <-t.tokensChan:
				process(v.ctx, v.tokens)
			}
		}
	}()

	return t
}

func init() {
	rand.Seed(time.Now().Unix())
}
