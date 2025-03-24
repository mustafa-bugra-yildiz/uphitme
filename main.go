package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/mustafa-bugra-yildiz/uphitme/env"

	_ "github.com/lib/pq"
)

var tmpl = template.Must(template.ParseFiles("templates.html"))

func main() {
	log.Println("Connecting the database")
	db, err := sql.Open("postgres", env.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Connected the database")

	mux := http.NewServeMux()
	mux.HandleFunc("/", landingPageHandler)
	mux.HandleFunc("/health", healthHandler)

	log.Printf("Starting server on port %s", env.Port)
	log.Fatal(http.ListenAndServe(":"+env.Port, mux))
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
