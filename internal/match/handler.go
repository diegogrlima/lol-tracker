package match

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type MatchLister interface {
	ListIDsByPUUID(
		ctx context.Context,
		puuid string,
		options ListOptions,
	) ([]string, error)
}

type Handler struct {
	matches MatchLister
	logger  *slog.Logger
}

type listResponse struct {
	MatchIDs   []string   `json:"matchIds"`
	Pagination pagination `json:"pagination"`
}

type pagination struct {
	Start int `json:"start"`
	Count int `json:"count"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewHandler(
	matches MatchLister,
	logger *slog.Logger,
) *Handler {
	if logger == nil {
		logger = slog.Default()
	}

	return &Handler{
		matches: matches,
		logger:  logger,
	}
}

func (h *Handler) Routes() http.Handler {
	router := chi.NewRouter()

	router.Get(
		"/by-puuid/{puuid}",
		h.ListIDsByPUUID,
	)

	return router
}

func (h *Handler) ListIDsByPUUID(
	w http.ResponseWriter,
	r *http.Request,
) {
	puuid := strings.TrimSpace(
		chi.URLParam(r, "puuid"),
	)

	options, err := parseListOptions(r)
	if err != nil {
		h.respondJSON(
			w,
			http.StatusBadRequest,
			errorResponse{
				Error: "invalid pagination parameters",
			},
		)
		return
	}

	matchIDs, err := h.matches.ListIDsByPUUID(
		r.Context(),
		puuid,
		options,
	)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(
		w,
		http.StatusOK,
		listResponse{
			MatchIDs: matchIDs,
			Pagination: pagination{
				Start: options.Start,
				Count: options.Count,
			},
		},
	)
}

func parseListOptions(r *http.Request) (ListOptions, error) {
	start, err := parseQueryInteger(
		r,
		"start",
		0,
	)
	if err != nil {
		return ListOptions{}, err
	}

	count, err := parseQueryInteger(
		r,
		"count",
		DefaultCount,
	)
	if err != nil {
		return ListOptions{}, err
	}

	return ListOptions{
		Start: start,
		Count: count,
	}, nil
}

func parseQueryInteger(
	r *http.Request,
	name string,
	defaultValue int,
) (int, error) {
	value := strings.TrimSpace(
		r.URL.Query().Get(name),
	)

	if value == "" {
		return defaultValue, nil
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (h *Handler) handleError(
	w http.ResponseWriter,
	err error,
) {
	switch {
	case errors.Is(err, ErrInvalidPUUID):
		h.respondJSON(
			w,
			http.StatusBadRequest,
			errorResponse{
				Error: "PUUID is required",
			},
		)

	case errors.Is(err, ErrInvalidPagination):
		h.respondJSON(
			w,
			http.StatusBadRequest,
			errorResponse{
				Error: err.Error(),
			},
		)

	default:
		h.logger.Error(
			"failed to list match IDs",
			"error",
			err,
		)

		h.respondJSON(
			w,
			http.StatusBadGateway,
			errorResponse{
				Error: "failed to communicate with Riot API",
			},
		)
	}
}

func (h *Handler) respondJSON(
	w http.ResponseWriter,
	status int,
	value any,
) {
	data, err := json.Marshal(value)
	if err != nil {
		h.logger.Error(
			"failed to encode HTTP response",
			"error",
			err,
		)

		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)
	w.WriteHeader(status)

	if _, err := w.Write(data); err != nil {
		h.logger.Error(
			"failed to write HTTP response",
			"error",
			err,
		)
	}
}
