package tts

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"sync"

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
	var cmd *exec.Cmd

	for _, t := range ts {
		switch t.Kind {
		case token.Number, token.Symbol:
			args := []string{}

			for _, r := range []rune(t.Text) {
				args = append(args, filepath.Join(global.CodespeakAudioPath, fmt.Sprintf("%03d.wav", r+1)))
			}
			if len(args) == 0 {
				continue
			}
			cmd = exec.CommandContext(ctx, "play", args...)
		case token.Alphabet:
			cmd = exec.CommandContext(ctx, "say", "-v", "Alex", "-r", "272", fmt.Sprintf("%q", "[[ pbas 42 ]]"+t.Text))
		case token.Unicode:
			cmd = exec.CommandContext(ctx, "say", "-v", "Kyoko", "-r", "480", fmt.Sprintf("%q", t.Text))
		}
		if err := cmd.Run(); err != nil {
			return err
		}
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
