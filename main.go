package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"html/template"
)

var tmpl = template.Must(template.ParseFiles("templates.html"))

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", landingPageHandler)
	mux.HandleFunc("/health", healthHandler)

	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
	})
}

func landingPageHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "landing-page", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}