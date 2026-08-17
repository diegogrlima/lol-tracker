package champion

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ChampionLister interface {
	List(ctx context.Context) ([]Champion, error)
}

type Handler struct {
	champions ChampionLister
	logger    *slog.Logger
}

type listResponse struct {
	Champions []Champion `json:"champions"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewHandler(champions ChampionLister, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}

	return &Handler{champions: champions, logger: logger}
}

func (h *Handler) Routes() http.Handler {
	router := chi.NewRouter()
	router.Get("/", h.List)
	return router
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	champions, err := h.champions.List(r.Context())
	if err != nil {
		h.logger.Error("failed to list champions", "error", err)
		h.respondJSON(w, http.StatusBadGateway, errorResponse{
			Error: "failed to communicate with Data Dragon",
		})
		return
	}

	h.respondJSON(w, http.StatusOK, listResponse{Champions: champions})
}

func (h *Handler) respondJSON(w http.ResponseWriter, status int, value any) {
	data, err := json.Marshal(value)
	if err != nil {
		h.logger.Error("failed to encode HTTP response", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		h.logger.Error("failed to write HTTP response", "error", err)
	}
}
