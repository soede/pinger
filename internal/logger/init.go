package logger

import (
	"docker/internal/config"
	"log/slog"
	"os"
)

func NewLogger(cfg *config.Config) *slog.Logger {
	opts := &slog.HandlerOptions{}
	switch cfg.App.EnvName {
	case "local":
		opts.Level = slog.LevelDebug
	case "staging":
		opts.Level = slog.LevelWarn
	case "prod":
		opts.Level = slog.LevelError
	default:
		opts.Level = slog.LevelInfo
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	return logger
}
