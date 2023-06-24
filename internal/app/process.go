package app

import (
	"context"

	"github.com/moutend/codespeak/internal/token"
	"github.com/moutend/codespeak/internal/tts"
)

func (a *App) Process(ctx context.Context, input string) error {
	tokens, err := token.Split(input)

	if err != nil {
		return err
	}

	option := tts.ProcessOption{
		AudioDir: a.audioDir,
	}

	if err := tts.Process(ctx, tokens, option); err != nil {
		return err
	}

	return nil
}
