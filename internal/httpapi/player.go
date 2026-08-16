package httpapi

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/diegogrlima/lol-tracker/internal/player"
	"github.com/go-chi/chi/v5"
)

type PlayerFinder interface {
	GetByRiotID(
		ctx context.Context,
		gameName string,
		tagLine string,
	) (*player.Player, error)
}

type PlayerHandler struct {
	players PlayerFinder
}

type playerResponse struct {
	Player *player.Player `json:"player"`
}

func NewPlayerHandler(players PlayerFinder) *PlayerHandler {
	return &PlayerHandler{players: players}
}

func (h *PlayerHandler) GetByRiotID(w http.ResponseWriter, r *http.Request) {
	gameName := strings.TrimSpace(chi.URLParam(r, "gameName"))
	tagLine := strings.TrimSpace(chi.URLParam(r, "tagLine"))

	if gameName == "" || tagLine == "" {
		respondJSON(w, http.StatusBadRequest, errorResponse{
			Error: "gameName and tagLine are required",
		})
		return
	}

	result, err := h.players.GetByRiotID(r.Context(), gameName, tagLine)
	if err != nil {
		handlePlayerError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, playerResponse{Player: result})
}

func handlePlayerError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, player.ErrNotFound):
		respondJSON(w, http.StatusNotFound, errorResponse{Error: "player not found"})
	case errors.Is(err, player.ErrRateLimited):
		respondJSON(w, http.StatusServiceUnavailable, errorResponse{
			Error: "Riot API temporarily unavailable",
		})
	default:
		log.Printf("get player: %v", err)
		respondJSON(w, http.StatusBadGateway, errorResponse{
			Error: "failed to communicate with Riot API",
		})
	}
}
