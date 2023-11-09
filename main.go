package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/mrsad/oidc/server/config"
	"github.com/mrsad/oidc/server/op"
	"github.com/mrsad/oidc/server/storage"

	"github.com/joho/godotenv"
	"golang.org/x/exp/slog"
)

var logger *slog.Logger
var cfg *config.Config

func init() {
	godotenv.Load()
	cfg = config.New()
	logLevel := slog.LevelInfo
	if cfg.DebugMode {
		logLevel = slog.LevelDebug
	}
	logger = slog.New(
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
			Level:     logLevel,
		}),
	)
}

func main() {
	issuer := fmt.Sprintf("%s://%s:%s/", cfg.Server.Proto, cfg.Server.Host, cfg.Server.Port)
	storage := storage.NewStorage(storage.NewUserStore(issuer))
	router := op.SetupServer(issuer, storage, logger, false)
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}
	logger.Info("server listening, press ctrl+c to stop", "addr", issuer)
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		logger.Error("server terminated", "error", err)
		os.Exit(1)
	}
}
