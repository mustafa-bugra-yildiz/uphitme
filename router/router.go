package router

import (
	"net/http"

	"github.com/mustafa-bugra-yildiz/uphitme/auth"
	"github.com/mustafa-bugra-yildiz/uphitme/middleware"
	"github.com/mustafa-bugra-yildiz/uphitme/repos/task"
	"github.com/mustafa-bugra-yildiz/uphitme/repos/user"
	"github.com/mustafa-bugra-yildiz/uphitme/scheduler"
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
	mux := http.NewServeMux()

	// shared state
	s := state{
		auth:      auth,
		taskRepo:  taskRepo,
		userRepo:  userRepo,
		scheduler: scheduler,
	}

	// api
	mux.HandleFunc("/api/health", healthHandler)
	mux.HandleFunc("/api/schedule", s.scheduleHandler)

	// public
	mux.HandleFunc("/", landingPageHandler)
	mux.HandleFunc("/sign-in", s.signInHandler)
	mux.HandleFunc("/sign-up", s.signUpHandler)
	mux.HandleFunc("/sign-up/success", signUpSuccessHandler)
	mux.HandleFunc("/sign-out", s.signOutHandler)

	// protected
	mux.HandleFunc("/dashboard", s.dashboardHomeHandler)
	mux.HandleFunc("/dashboard/scheduler", s.dashboardSchedulerHandler)
	mux.HandleFunc("/dashboard/billing-usage", s.dashboardBillingUsageHandler)
	mux.HandleFunc("/dashboard/profile", s.dashboardProfileHandler)

	// middleware
	handler := middleware.Logger(mux)

	// done
	return handler
}

func (s state) signOutHandler(w http.ResponseWriter, r *http.Request) {
	s.auth.Logout(w)
	http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
}
