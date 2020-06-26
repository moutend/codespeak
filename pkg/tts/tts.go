package tts

import (
	"context"
	"math/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
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

func process(ctx context.Context, ts []token.Token) error {
	tempPath, err := ioutil.TempDir("", "codespeak")

	if err != nil {
		return err
	}

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

			aiffPath := filepath.Join(tempPath, fmt.Sprintf("%s.aiff", hex.EncodeToString(hash)))
			cmd := exec.CommandContext(ctx, "say", "-v", "Alex", "-r", "272", fmt.Sprintf("%q", "[[ pbas 42 ]]"+t.Text), "-o", aiffPath)

			if err := cmd.Run(); err != nil {
				continue
			}

			audioPaths = append(audioPaths, aiffPath)
		case token.Unicode:
			hash := make([]byte, 16, 16)

			if _, err := rand.Read(hash); err != nil {
				continue
			}

			aiffPath := filepath.Join(tempPath, fmt.Sprintf("%s.aiff", hex.EncodeToString(hash)))
			cmd := exec.CommandContext(ctx, "say", "-v", "Kyoko", "-r", "480", fmt.Sprintf("%q", t.Text), "-o", aiffPath)

			if err := cmd.Run(); err != nil {
				continue
			}

			audioPaths = append(audioPaths, aiffPath)
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
