package player

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Finder interface {
	GetByRiotID(ctx context.Context, gameName string, tagLine string) (*Player, error)
}

type Handler struct {
	players Finder
	logger  *slog.Logger
}

type playerResponse struct {
	Player *Player `json:"player"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewHandler(players Finder, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}

	return &Handler{players: players, logger: logger}
}

func (h *Handler) Routes() http.Handler {
	router := chi.NewRouter()
	router.Get("/{gameName}/{tagLine}", h.GetByRiotID)
	return router
}

func (h *Handler) GetByRiotID(w http.ResponseWriter, r *http.Request) {
	gameName := strings.TrimSpace(chi.URLParam(r, "gameName"))
	tagLine := strings.TrimSpace(chi.URLParam(r, "tagLine"))

	if gameName == "" || tagLine == "" {
		h.respondJSON(w, http.StatusBadRequest, errorResponse{
			Error: "gameName and tagLine are required",
		})
		return
	}

	result, err := h.players.GetByRiotID(r.Context(), gameName, tagLine)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, playerResponse{Player: result})
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		h.respondJSON(w, http.StatusNotFound, errorResponse{Error: "player not found"})
	case errors.Is(err, ErrRateLimited):
		h.respondJSON(w, http.StatusServiceUnavailable, errorResponse{
			Error: "Riot API temporarily unavailable",
		})
	case errors.Is(err, ErrCacheUnavailable):
		h.logger.Warn("player cache unavailable", "error", err)
		h.respondJSON(w, http.StatusServiceUnavailable, errorResponse{
			Error: "service temporarily unavailable",
		})
	default:
		h.logger.Error("failed to get player", "error", err)
		h.respondJSON(w, http.StatusBadGateway, errorResponse{
			Error: "failed to communicate with Riot API",
		})
	}
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
