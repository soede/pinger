package main

import (
	"log/slog"
	"os"
)

func NewLogger(cfg *Config) *slog.Logger {
	opts := &slog.HandlerOptions{}
	switch cfg.EnvName {
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
