package router

import (
	"net/http"

	"github.com/mustafa-bugra-yildiz/uphitme/view"
)

func (s state) dashboardHomeHandler(w http.ResponseWriter, r *http.Request) {
	user, err := s.auth.User(r)
	if err != nil {
		http.Redirect(w, r, "/sign-in", http.StatusUnauthorized)
		return
	}

	err = view.Dashboard.Home(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s state) dashboardSchedulerHandler(w http.ResponseWriter, r *http.Request) {
	user, err := s.auth.User(r)
	if err != nil {
		http.Redirect(w, r, "/sign-in", http.StatusUnauthorized)
		return
	}

	err = view.Dashboard.Scheduler(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s state) dashboardBillingUsageHandler(w http.ResponseWriter, r *http.Request) {
	user, err := s.auth.User(r)
	if err != nil {
		http.Redirect(w, r, "/sign-in", http.StatusUnauthorized)
		return
	}

	err = view.Dashboard.BillingUsage(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s state) dashboardProfileHandler(w http.ResponseWriter, r *http.Request) {
	user, err := s.auth.User(r)
	if err != nil {
		http.Redirect(w, r, "/sign-in", http.StatusUnauthorized)
		return
	}

	err = view.Dashboard.Profile(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
