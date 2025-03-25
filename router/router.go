package router

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/mustafa-bugra-yildiz/uphitme/middleware"
	"github.com/mustafa-bugra-yildiz/uphitme/page"
)

func New() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", landingPageHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/schedule", scheduleHandler)

	return middleware.Logger(mux)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
	})
}

func landingPageHandler(w http.ResponseWriter, r *http.Request) {
	err := page.Landing(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func scheduleHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title   string         `json:"title"`
		Target  string         `json:"target"`
		Payload map[string]any `json:"payload"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	go func() {
		payload, _ := json.Marshal(req.Payload)
		http.Post(req.Target, "application/json", bytes.NewBuffer(payload))
	}()

	msg := "task scheduled, we will call you back"
	json.NewEncoder(w).Encode(Response[string]{
		Succeeded: true,
		Data:      &msg,
		Error:     nil,
	})
}

// TODO: Move this to somewhere else maybe?
type Response[T any] struct {
	Succeeded bool    `json:"succeeded"`
	Data      *T      `json:"data"`
	Error     *string `json:"error"`
}
