package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/api/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Peran 3 dan 4 menyambungkan Store, Commander, Broadcaster, dan route di sini.
	log.Println("backend berjalan di :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

