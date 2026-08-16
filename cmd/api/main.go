package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]any{
			"message": "Servidor funcionando League of Legends Tracker",
			"status":  true,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error ao gerar json", http.StatusInternalServerError)
		}

	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {

		fmt.Println("Erro ao iniciar o servidor", err)
	}

}
