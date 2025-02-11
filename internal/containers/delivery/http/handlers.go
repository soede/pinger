package http

import (
	"context"
	"docker/internal/config"
	"docker/internal/containers"
	"docker/internal/models"
	"encoding/json"
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

		page := 1
		size := 10
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
			http.Error(w, "failed to get containers", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(list)
	}
}

func (h containersHandlers) SearchByIP() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		query := request.URL.Query().Get("query")
		if query == "" {
			http.Error(writer, "query parameter is required", http.StatusBadRequest)
			return
		}

		var list []*models.Container
		var err error

		if net.ParseIP(query) != nil {
			list, err = h.conUC.GetByIP(ctx, query)
		} else {
			list, err = h.conUC.GetByStatus(ctx, query)
		}

		if err != nil {
			http.Error(writer, "failed to fetch containers", http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(writer).Encode(list); err != nil {
			http.Error(writer, "failed to encode response", http.StatusInternalServerError)
			return
		}

	}
}

func (h containersHandlers) SetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var list []*models.Container
		if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)

			return
		}

		err := h.conUC.SetAll(ctx, list)
		if err != nil {
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

		page := 1
		size := 10
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
			http.Error(w, "failed to get containers", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(list); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
