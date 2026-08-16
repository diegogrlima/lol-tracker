package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

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

type playerResponse struct {
	Player *riot.Account `json:"player"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewPlayer(riotClient RiotAccountClient) *Player {
	return &Player{
		riotClient: riotClient,
	}
}

func (p *Player) GetByRiotID(w http.ResponseWriter, r *http.Request) {
	gameName := strings.TrimSpace(chi.URLParam(r, "gameName"))
	tagLine := strings.TrimSpace(chi.URLParam(r, "tagLine"))

	if gameName == "" || tagLine == "" {
		respondJSON(w, http.StatusBadRequest, errorResponse{
			Error: "gameName and tagLine are required",
		})
		return
	}

	account, err := p.riotClient.GetAccountByRiotID(
		r.Context(),
		gameName,
		tagLine,
	)
	if err != nil {
		p.handleRiotError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, playerResponse{
		Player: account,
	})
}

func (p *Player) handleRiotError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, riot.ErrAccountNotFound):
		respondJSON(w, http.StatusNotFound, errorResponse{
			Error: "player not found",
		})

	case errors.Is(err, riot.ErrRateLimited):
		respondJSON(w, http.StatusServiceUnavailable, errorResponse{
			Error: "Riot API temporarily unavailable",
		})

	default:
		log.Printf("get Riot account: %v", err)

		respondJSON(w, http.StatusBadGateway, errorResponse{
			Error: "failed to communicate with Riot API",
		})
	}
}

func respondJSON(w http.ResponseWriter, status int, value any) {
	data, err := json.Marshal(value)
	if err != nil {
		log.Printf("encode HTTP response: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(data); err != nil {
		log.Printf("write HTTP response: %v", err)
	}
}
