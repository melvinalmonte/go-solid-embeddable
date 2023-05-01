package main

import (
	"context"
	"errors"
	"go-solid-embed/frontend"
	"go-solid-embed/handlers"
	"go-solid-embed/utils"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	utils.InitLogger()
	zap.S().Info("Starting server...")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	contents, err := frontend.Contents()
	if err != nil {
		zap.S().Errorw("failed to find assets", "error", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", handlers.FileServer(contents))
	mux.HandleFunc("/api", handlers.RootResource)
	mux.HandleFunc("/api/todos", handlers.HandleTodos)

	server := &http.Server{
		Addr:    ":9090",
		Handler: mux,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			zap.S().Errorw("failed to shutdown server", "error", err)
		}
	}()

	zap.S().Info("Server is ready to handle requests at http://localhost:9090/")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		zap.S().Errorw("server error", "error", err)
	}
}
