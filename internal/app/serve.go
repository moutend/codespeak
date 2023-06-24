package app

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/moutend/codespeak/internal/token"
)

func (a *App) Start(ctx context.Context, u *url.URL) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/speak", a.speakHandler)
	mux.HandleFunc("/v1/pause", a.pauseHandler)

	a.server = &http.Server{
		Addr:    u.String(),
		Handler: mux,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		a.debug.Println("Starting server:", u)

		if err := a.server.ListenAndServe(); err != nil {
			a.debug.Println(err)
		}
	}()

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	a.debug.Println("Stopping server")

	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func (a *App) speakHandler(w http.ResponseWriter, r *http.Request) {
	a.debug.Println("Called /v1/speak")

	data, err := os.ReadFile(a.inputFilePath)

	if err != nil {
		http.Error(w, `{"error": "invalid request"}`, http.StatusBadRequest)

		return
	}

	tokens, err := token.Split(string(data))

	if err != nil {
		http.Error(w, `{"error": "failed to parse request"}`, http.StatusBadRequest)

		return
	}

	a.engine.Speak(context.Background(), tokens, a.processOption)
}

func (a *App) pauseHandler(w http.ResponseWriter, r *http.Request) {
	a.debug.Println("Called /v1/pause")
	a.engine.Pause(a.processOption)
}
