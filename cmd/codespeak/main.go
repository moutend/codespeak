package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/moutend/codespeak/pkg/global"
	"github.com/moutend/codespeak/pkg/token"
	"github.com/moutend/codespeak/pkg/tts"
)

var (
	engine     *tts.TTS
	bufferPath string
)

func main() {
	log.SetPrefix("error: ")
	log.SetFlags(0)

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	engine = tts.Open()
	defer engine.Close()

	bufferPath = filepath.Join(global.CodespeakBufferPath, "buffer.txt")

	http.HandleFunc("/v1/speak", speakHandler)
	http.HandleFunc("/v1/pause", pauseHandler)

	return http.ListenAndServe(":8501", nil)
}

func speakHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile(bufferPath)

	if err != nil {
		http.Error(w, `{"error": "invalid request"}`, http.StatusBadRequest)

		return
	}
	tokens, err := token.Split(string(data))

	if err != nil {
		http.Error(w, `{"error": "invalid request"}`, http.StatusBadRequest)

		return
	}

	engine.SpeakContext(context.Background(), tokens)

	return
}

func pauseHandler(w http.ResponseWriter, r *http.Request) {
	engine.Pause()
}
