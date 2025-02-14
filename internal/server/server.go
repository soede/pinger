package server

import (
	"context"
	"docker/internal/config"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	mux *http.ServeMux
	cfg *config.Config
	db  *sqlx.DB
}

const (
	maxHeaderBytes = 1 << 16
	ctxTimeout     = 5
)

func NewServer(cfg *config.Config, db *sqlx.DB) *Server {
	return &Server{http.NewServeMux(), cfg, db}
}

func (s *Server) Run() error {
	var writeTimeout, readTimeout = 5, 5
	if s.cfg.WriteTimeout != 0 {
		writeTimeout = s.cfg.WriteTimeout
	}
	if s.cfg.App.ReadTimeout != 0 {
		writeTimeout = s.cfg.ReadTimeout
	}

	server := &http.Server{
		Addr:           ":" + s.cfg.App.Port,
		Handler:        s.mux,
		ReadTimeout:    time.Duration(readTimeout) * time.Second,
		WriteTimeout:   time.Duration(writeTimeout) * time.Second,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		slog.Info("Server is listening on PORT", slog.String("port", s.cfg.App.Port))
		if err := server.ListenAndServe(); err != nil {
			slog.Error("Error starting Server: ", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	if err := s.MapHandlers(); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()
	slog.Info("Server Exited Properly")

	return server.Shutdown(ctx)

}
