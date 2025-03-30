package router

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/mustafa-bugra-yildiz/uphitme/repos/task"
)

type Response[T any] struct {
	Succeeded bool    `json:"succeeded"`
	Data      *T      `json:"data"`
	Error     *string `json:"error"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"status": "ok",
	}

	json.NewEncoder(w).Encode(Response[map[string]any]{
		Succeeded: true,
		Data: &data,
		Error: nil,
	})
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

	go s.scheduler.Schedule(context.Background(), taskID)

	msg := "task scheduled, we will call you back"
	json.NewEncoder(w).Encode(Response[string]{
		Succeeded: true,
		Data:      &msg,
		Error:     nil,
	})
}
