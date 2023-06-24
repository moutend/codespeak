package tts

import (
	"context"
	"testing"
	"time"

	"github.com/moutend/codespeak/pkg/token"
)

func TestTTS(t *testing.T) {
	tts := Open()
	defer tts.Close()

	tts.Pause()

	tts.SpeakContext(context.Background(), []token.Token{
		{
			Kind: token.Unicode,
			Text: "おはようございます、私の名前は強固です。日本語の音声をお届けします。",
		},
	})
	time.Sleep(10 * time.Second)

	tts.SpeakContext(context.Background(), []token.Token{
		{
			Kind: token.Unicode,
			Text: "こんにちは、私の名前は強固です。日本語の音声をお届けします。",
		},
	})
	time.Sleep(1 * time.Second)

	tts.Pause()
	time.Sleep(1 * time.Second)

	tts.SpeakContext(context.Background(), []token.Token{
		{
			Kind: token.Unicode,
			Text: "こんばんは、私の名前は強固です。日本語の音声をお届けします。",
		},
	})
}
