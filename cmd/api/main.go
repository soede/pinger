package main

import (
	"docker/internal/config"
	log "docker/internal/logger"
	"docker/internal/server"
	"docker/pkg/db/postgres"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	cfg := config.LoadConfig()
	slog.SetDefault(log.NewLogger(cfg))
	slog.Debug(fmt.Sprintf("ReadTimeout: %d, WriteTimeout: %d", cfg.ReadTimeout, cfg.WriteTimeout))
	db, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		slog.Error("failed DB init", "Error: ", err)
		os.Exit(1)
	}
	defer db.Close()

	s := server.NewServer(cfg, db)
	if err = s.Run(); err != nil {
		slog.Error("not ok")
	}
}
