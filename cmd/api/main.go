package main

import (
	"docker/internal/config"
	log "docker/internal/logger"
	"docker/internal/server"
	"docker/pkg/db/postgres"
	"log/slog"
	"os"
)

func main() {
	cfg := config.LoadConfig()
	slog.SetDefault(log.NewLogger(cfg))
	slog.Info("",
		slog.String("env_name", cfg.App.EnvName),
		slog.String("SSL", cfg.SSL),
		slog.Int("ReadTimeout(s)", cfg.ReadTimeout),
		slog.Int("WriteTimeout(s)", cfg.WriteTimeout),
	)
	db, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		slog.Error("Postgresql init: %s", err)
		os.Exit(1)
	}
	slog.Info("Postgres connected", "Status", db.Stats())

	defer db.Close()

	s := server.NewServer(cfg, db)
	if err = s.Run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
