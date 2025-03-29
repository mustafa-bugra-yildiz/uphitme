package router

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"net/url"
	"regexp"
	"strings"

	"github.com/mustafa-bugra-yildiz/uphitme/auth"
	"github.com/mustafa-bugra-yildiz/uphitme/env"
	"github.com/mustafa-bugra-yildiz/uphitme/middleware"
	"github.com/mustafa-bugra-yildiz/uphitme/page"
	"github.com/mustafa-bugra-yildiz/uphitme/repos/task"
	"github.com/mustafa-bugra-yildiz/uphitme/repos/user"
	"github.com/mustafa-bugra-yildiz/uphitme/scheduler"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"
)

type state struct {
	auth      *auth.Auth
	taskRepo  task.Repo
	userRepo  user.Repo
	scheduler *scheduler.Scheduler
}

func New(
	auth *auth.Auth,
	taskRepo task.Repo,
	userRepo user.Repo,
	scheduler *scheduler.Scheduler,
) http.Handler {
	s := state{
		auth:      auth,
		taskRepo:  taskRepo,
		userRepo:  userRepo,
		scheduler: scheduler,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", landingPageHandler)

	mux.HandleFunc("/sign-in", s.signInHandler)
	mux.HandleFunc("/sign-up", s.signUpHandler)
	mux.HandleFunc("/sign-up/success", signUpSuccessHandler)
	mux.HandleFunc("/sign-out", s.signOutHandler)

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
	_, err := s.auth.Verify(r)
	if err != nil {
		http.Redirect(w, r, "/sign-in", http.StatusUnauthorized)
		return
	}

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

	err = g.Wait()
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

func (s state) signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := page.SignIn(w, "", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))

	_, err := mail.ParseAddress(email)
	if err != nil {
		page.SignIn(w, email, err)
		return
	}

	if password == "" {
		err = errors.New("password cannot be empty")
		page.SignIn(w, email, err)
		return
	}

	user_, err := s.userRepo.Get(r.Context(), email)
	if err == user.ErrUserNotFound {
		err = errors.New("invalid credentials")
		page.SignIn(w, email, err)
		return
	}
	if err != nil {
		page.SignIn(w, email, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user_.HashedPassword), []byte(password))
	if err != nil {
		err = errors.New("invalid credentials")
		page.SignIn(w, email, err)
		return
	}

	err = s.auth.Login(w, *user_)
	if err != nil {
		page.SignIn(w, email, err)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (s state) signUpHandler(w http.ResponseWriter, r *http.Request) {
	if env.Env == env.Prod {
		err := page.SignUpWIP(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if r.Method != http.MethodPost {
		err := page.SignUp(w, "", "", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	fullName := strings.TrimSpace(r.FormValue("full-name"))
	email := strings.TrimSpace(r.FormValue("email"))

	if fullName == "" {
		err := errors.New("your full name cannot be empty")
		page.SignUp(w, fullName, email, err)
		return
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		page.SignUp(w, fullName, email, err)
		return
	}

	password := strings.TrimSpace(r.FormValue("password"))
	err = validatePassword(password)
	if err != nil {
		page.SignUp(w, fullName, email, err)
		return
	}

	_, err = s.userRepo.Get(r.Context(), email)
	if err != user.ErrUserNotFound {
		err = errors.New("there is already an account with that email")
		page.SignUp(w, fullName, email, err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		page.SignUp(w, fullName, email, err)
		return
	}

	err = s.userRepo.Create(r.Context(), fullName, email, string(hashedPassword))
	if err != nil {
		page.SignUp(w, fullName, email, err)
		return
	}

	http.Redirect(w, r, "/sign-up/success", http.StatusSeeOther)
}

func signUpSuccessHandler(w http.ResponseWriter, r *http.Request) {
	err := page.SignUpSuccess(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s state) signOutHandler(w http.ResponseWriter, r *http.Request) {
	s.auth.Logout(w)
	http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
}

// TODO: Move this to somewhere else maybe?
type Response[T any] struct {
	Succeeded bool    `json:"succeeded"`
	Data      *T      `json:"data"`
	Error     *string `json:"error"`
}

// TODO: Move this to somewhere else maybe?
func validatePassword(password string) error {
	// Rule 1: Password must be at least 8 characters long
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Rule 2: Password must contain at least one uppercase letter
	upperCaseRegex := regexp.MustCompile(`[A-Z]`)
	if !upperCaseRegex.MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	// Rule 3: Password must contain at least one number
	numberRegex := regexp.MustCompile(`[0-9]`)
	if !numberRegex.MatchString(password) {
		return errors.New("password must contain at least one number")
	}

	return nil // No errors, password is valid
}
