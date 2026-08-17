package item

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type ItemQuery interface {
	ListStoreItems(ctx context.Context) ([]Item, error)
	GetByID(ctx context.Context, itemID string) (*Item, error)
}

type Handler struct {
	items  ItemQuery
	logger *slog.Logger
}

type listResponse struct {
	Items []Item `json:"items"`
}

type itemResponse struct {
	Item *Item `json:"item"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewHandler(items ItemQuery, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}

	return &Handler{items: items, logger: logger}
}

func (h *Handler) Routes() http.Handler {
	router := chi.NewRouter()
	router.Get("/", h.List)
	router.Get("/{itemID}", h.GetByID)
	return router
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.items.ListStoreItems(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, listResponse{Items: items})
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	itemID := strings.TrimSpace(chi.URLParam(r, "itemID"))

	storeItem, err := h.items.GetByID(r.Context(), itemID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, itemResponse{Item: storeItem})
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidItemID):
		h.respondJSON(w, http.StatusBadRequest, errorResponse{Error: "item ID is required"})
	case errors.Is(err, ErrItemNotFound):
		h.respondJSON(w, http.StatusNotFound, errorResponse{Error: "item not found"})
	default:
		h.logger.Error("failed to query items", "error", err)
		h.respondJSON(w, http.StatusBadGateway, errorResponse{
			Error: "failed to communicate with Data Dragon",
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
