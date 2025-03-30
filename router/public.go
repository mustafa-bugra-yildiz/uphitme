package router

import (
	"errors"
	"net/http"
	"net/mail"
	"regexp"
	"strings"

	"github.com/mustafa-bugra-yildiz/uphitme/env"
	"github.com/mustafa-bugra-yildiz/uphitme/repos/user"
	"github.com/mustafa-bugra-yildiz/uphitme/view"
	"golang.org/x/crypto/bcrypt"
)

func landingPageHandler(w http.ResponseWriter, r *http.Request) {
	err := view.Landing(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s state) signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		err := view.SignIn(w, "", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))

	_, err := mail.ParseAddress(email)
	if err != nil {
		view.SignIn(w, email, err)
		return
	}

	if password == "" {
		err = errors.New("password cannot be empty")
		view.SignIn(w, email, err)
		return
	}

	user_, err := s.userRepo.GetByEmail(r.Context(), email)
	if err == user.ErrUserNotFound {
		err = errors.New("invalid credentials")
		view.SignIn(w, email, err)
		return
	}
	if err != nil {
		view.SignIn(w, email, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user_.HashedPassword), []byte(password))
	if err != nil {
		err = errors.New("invalid credentials")
		view.SignIn(w, email, err)
		return
	}

	err = s.auth.Login(w, *user_)
	if err != nil {
		view.SignIn(w, email, err)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (s state) signUpHandler(w http.ResponseWriter, r *http.Request) {
	if env.Env == env.Prod {
		err := view.SignUpWIP(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if r.Method != http.MethodPost {
		err := view.SignUp(w, "", "", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	fullName := strings.TrimSpace(r.FormValue("full-name"))
	email := strings.TrimSpace(r.FormValue("email"))

	if fullName == "" {
		err := errors.New("your full name cannot be empty")
		view.SignUp(w, fullName, email, err)
		return
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		view.SignUp(w, fullName, email, err)
		return
	}

	password := strings.TrimSpace(r.FormValue("password"))
	err = validatePassword(password)
	if err != nil {
		view.SignUp(w, fullName, email, err)
		return
	}

	_, err = s.userRepo.GetByEmail(r.Context(), email)
	if err != user.ErrUserNotFound {
		err = errors.New("there is already an account with that email")
		view.SignUp(w, fullName, email, err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		view.SignUp(w, fullName, email, err)
		return
	}

	err = s.userRepo.Create(r.Context(), fullName, email, string(hashedPassword))
	if err != nil {
		view.SignUp(w, fullName, email, err)
		return
	}

	http.Redirect(w, r, "/sign-up/success", http.StatusSeeOther)
}

func signUpSuccessHandler(w http.ResponseWriter, r *http.Request) {
	err := view.SignUpSuccess(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

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
