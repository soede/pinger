package main

import (
	"docker/internal/config"
	"docker/internal/server"
	"docker/pkg/db/postgres"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.MustEnvConfig()

	db, err := postgres.NewPsqlDB(&cfg)
	if err != nil {
		slog.Error("failed DB init", "Error: ", err)
		os.Exit(1)
	}
	defer db.Close()

	s := server.NewServer(&cfg, db, *logger)
	if err = s.Run(); err != nil {
		logger.Info("not ok")
	}
}
