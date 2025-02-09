package http

import (
	"docker/internal/containers"
	"net/http"
)

func MapContainersRoutes(mux *http.ServeMux, h containers.Handlers) {
	mux.HandleFunc("GET /api/v1/getAll", enableCORS(h.GetAll()))
	mux.HandleFunc("GET /api/v1/history", enableCORS(h.GetHistory()))
	mux.HandleFunc("GET /api/v1/search", enableCORS(h.SearchByIP()))
	mux.HandleFunc("POST /api/v1/setAll", enableCORS(h.SetAll()))
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
