package handler

import (
	"encoding/json"
	"net/http"
)

func GetHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]any{
		"message": "Servidor funcionando League of Legends Tracker",
		"status":  true,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error ao gerar json", http.StatusInternalServerError)
	}
}
