package handler

import (
	"encoding/json"
	"net/http"
)

type Health struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

func (h *Health) GetHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]Health{
		"health": {
			Message: "API funcionando",
			Status:  true,
		},
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error ao gerar json", http.StatusInternalServerError)
	}
}
