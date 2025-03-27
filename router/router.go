package router

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/mustafa-bugra-yildiz/uphitme/middleware"
	"github.com/mustafa-bugra-yildiz/uphitme/page"
	"github.com/mustafa-bugra-yildiz/uphitme/repos/task"
)

type state struct {
	taskRepo task.Repo
}

func New(taskRepo task.Repo) http.Handler {
	s := state{taskRepo: taskRepo}
	mux := http.NewServeMux()

	mux.HandleFunc("/", landingPageHandler)
	mux.HandleFunc("/dashboard", s.dashboardPageHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/schedule", s.scheduleHandler)

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

func (s state) dashboardPageHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := s.taskRepo.List(r.Context(), 1, 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = page.Dashboard(w, tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s state) scheduleHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title   string       `json:"title"`
		Target  string       `json:"target"`
		Payload task.Payload `json:"payload"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	target, err := url.Parse(req.Target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	taskID, err := s.taskRepo.Create(
		r.Context(),
		strings.TrimSpace(req.Title),
		*target,
		req.Payload,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go func() error {
		ctx := context.Background()

		payload, _ := json.Marshal(req.Payload)
		res, err := http.Post(req.Target, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			return s.taskRepo.Fail(ctx, taskID, err)
		}
		defer res.Body.Close()

		var body task.Payload
		err = json.NewDecoder(res.Body).Decode(&body)
		if err != nil {
			return s.taskRepo.Fail(ctx, taskID, err)
		}

		return s.taskRepo.Succeed(ctx, taskID, res.StatusCode, body)
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
