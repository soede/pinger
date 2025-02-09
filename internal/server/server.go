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
	mux    *http.ServeMux
	cfg    *config.Config
	db     *sqlx.DB
	logger slog.Logger
}

const (
	maxHeaderBytes = 1 << 20
	readTimeout    = 5
	writeTimeout   = 5
	ctxTimeout     = 5
)

func NewServer(cfg *config.Config, db *sqlx.DB, logger slog.Logger) *Server {
	return &Server{http.NewServeMux(), cfg, db, logger}
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:           ":" + s.cfg.App.Port,
		Handler:        s.mux,
		ReadTimeout:    readTimeout * time.Second,
		WriteTimeout:   writeTimeout * time.Second,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		s.logger.Info("Server is listening on PORT", slog.String("port", s.cfg.App.Port))
		if err := server.ListenAndServe(); err != nil {
			s.logger.Error("Error starting Server: ", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	/*go func() {
		s.logger.Info("Starting Debug Server on PORT: ", s.cfg.App.PprofPort)
	}()*/

	if err := s.MapHandlers(); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()
	s.logger.Info("Server Exited Properly")

	return server.Shutdown(ctx)

}
