package http

import (
	"context"
	"docker/internal/config"
	"docker/internal/containers"
	"docker/internal/models"
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"strconv"
)

type containersHandlers struct {
	cfg   *config.Config
	conUC containers.UseCase
}

func NewContainerHandlers(cfg *config.Config, conUC containers.UseCase) containers.Handlers {
	return &containersHandlers{cfg: cfg, conUC: conUC}
}

func (h containersHandlers) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		slog.Debug("Handling request for containers list ((cUseCase.GetAll))",
			slog.String("url", r.URL.String()))

		page := 1
		size := 100
		if p := r.URL.Query().Get("page"); p != "" {
			if parsed, err := strconv.Atoi(p); err == nil {
				page = parsed
			}
		}
		if p := r.URL.Query().Get("size"); p != "" {
			if parsed, err := strconv.Atoi(p); err == nil {
				page = parsed
			}
		}

		list, err := h.conUC.GetAll(ctx, page, size)
		if err != nil {
			slog.Error("Failed to fetch containers list (cUseCase.GetAll)",
				slog.String("error", err.Error()),
				slog.String("url", r.URL.String()))
			http.Error(w, "failed to get containers", http.StatusInternalServerError)
			return
		}

		slog.Debug("Fetched containers list (cUseCase.GetAll)",
			slog.String("from", r.URL.String()),
			slog.Any("containers", list),
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(list)
	}
}

func (h containersHandlers) SearchByIP() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		query := r.URL.Query().Get("query")
		if query == "" {
			slog.Warn("Missing query parameter (cUseCase.SearchByIP)",
				slog.String("message", "query parameter is required"),
				slog.String("url", r.URL.String()))
			http.Error(w, "query parameter is required", http.StatusBadRequest)
			return
		}

		slog.Debug("Searching containers (cUseCase.SearchByIP)",
			slog.String("query", query),
			slog.String("url", r.URL.String()))

		var list []*models.Container
		var err error

		if net.ParseIP(query) != nil {
			list, err = h.conUC.GetByIP(ctx, query)
		} else {
			list, err = h.conUC.GetByStatus(ctx, query)
		}

		if err != nil {
			slog.Error("Failed to fetch containers (cUseCase.SearchByIP)",
				slog.String("error", err.Error()),
				slog.String("url", r.URL.String()))
			http.Error(w, "failed to fetch containers", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(list); err != nil {
			slog.Error("Failed to encode response (cUseCase.SearchByIP)",
				slog.String("error", err.Error()),
				slog.String("url", r.URL.String()))
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}

	}
}

func (h containersHandlers) SetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		slog.Debug("Received request to set all containers (cUseCase.SetAll)", slog.String("url", r.URL.String()))

		var list []*models.Container
		if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
			slog.Warn("Invalid request body (cUseCase.SetAll)", slog.String("url", r.URL.String()), slog.String("error", err.Error()), slog.Any("body", r.Body))
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		slog.Debug("Decoded container list (cUseCase.SetAll)", slog.String("url", r.URL.String()), slog.Any("list", list))

		err := h.conUC.SetAll(ctx, list)
		if err != nil {
			slog.Error("Failed to set containers (cUseCase.SetAll)", slog.String("error", err.Error()), slog.String("url", r.URL.String()))
			http.Error(w, "failed to set containers", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
func (h containersHandlers) GetHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		slog.Debug("Received request to fetch containers history (cUseCase.GetHistory)", slog.String("url", r.URL.String()))

		page := 1
		size := 100
		if p := r.URL.Query().Get("page"); p != "" {
			if parsed, err := strconv.Atoi(p); err == nil {
				page = parsed
			}
		}
		if p := r.URL.Query().Get("size"); p != "" {
			if parsed, err := strconv.Atoi(p); err == nil {
				page = parsed
			}
		}

		list, err := h.conUC.GetHistory(ctx, page, size)
		if err != nil {
			slog.Error("Failed to get containers (cUseCase.GetHistory)", slog.String("url", r.URL.String()), slog.String("error", err.Error()))
			http.Error(w, "failed to get containers", http.StatusInternalServerError)
			return
		}
		slog.Info("Successfully fetched container history (cUseCase.GetHistory)", slog.String("url", r.URL.String()), slog.Any("count", list))

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(list); err != nil {
			slog.Error("Failed to encode response (cUseCase.GetHistory)", slog.String("url", r.URL.String()), slog.String("error", err.Error()))
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
