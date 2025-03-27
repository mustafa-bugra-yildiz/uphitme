package router

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/mustafa-bugra-yildiz/uphitme/middleware"
	"github.com/mustafa-bugra-yildiz/uphitme/page"
	"github.com/mustafa-bugra-yildiz/uphitme/repos/task"
	"github.com/mustafa-bugra-yildiz/uphitme/repos/user"
	"github.com/mustafa-bugra-yildiz/uphitme/scheduler"

	"golang.org/x/sync/errgroup"
)

type state struct {
	taskRepo  task.Repo
	userRepo  user.Repo
	scheduler *scheduler.Scheduler
}

func New(
	taskRepo task.Repo,
	userRepo user.Repo,
	scheduler *scheduler.Scheduler,
) http.Handler {
	s := state{
		taskRepo:  taskRepo,
		userRepo:  userRepo,
		scheduler: scheduler,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", landingPageHandler)
	mux.HandleFunc("/sign-in", signInHandler)
	mux.HandleFunc("/sign-up", signUpHandler)
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
	g := new(errgroup.Group)

	var tasks []task.Task
	g.Go(func() error {
		var err error
		tasks, err = s.taskRepo.List(r.Context(), 1, 10)
		return err
	})

	var count int
	g.Go(func() error {
		var err error
		count, err = s.taskRepo.Count(r.Context())
		return err
	})

	err := g.Wait()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = page.Dashboard(w, 1, 10, count, tasks)
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

	go s.scheduler.Schedule(context.Background(), taskID)

	msg := "task scheduled, we will call you back"
	json.NewEncoder(w).Encode(Response[string]{
		Succeeded: true,
		Data:      &msg,
		Error:     nil,
	})
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	err := page.SignIn(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	err := page.SignUp(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// TODO: Move this to somewhere else maybe?
type Response[T any] struct {
	Succeeded bool    `json:"succeeded"`
	Data      *T      `json:"data"`
	Error     *string `json:"error"`
}
