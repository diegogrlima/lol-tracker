package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/diegogrlima/lol-tracker/internal/riot"
	"github.com/go-chi/chi/v5"
)

type RiotAccountClient interface {
	GetAccountByRiotID(
		ctx context.Context,
		gameName string,
		tagLine string,
	) (*riot.Account, error)
}

type Player struct {
	riotClient RiotAccountClient
}

func NewPlayer(riotClient RiotAccountClient) *Player {
	return &Player{riotClient: riotClient}
}

func (p *Player) GetPlayer(w http.ResponseWriter, r *http.Request) {
	gameName := chi.URLParam(r, "gameName")
	tagLine := chi.URLParam(r, "tagLine")

	account, err := p.riotClient.GetAccountByRiotID(r.Context(), gameName, tagLine)
	if err != nil {
		switch {
		case errors.Is(err, riot.ErrAccountNotFound):
			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": "player not found",
			})

		case errors.Is(err, riot.ErrRateLimited):
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{
				"error": "Riot API temporarily unavailable",
			})

		default:
			writeJSON(w, http.StatusBadGateway, map[string]string{
				"error": "failed to communicate with Riot API",
			})
		}

		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"player": account,
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
