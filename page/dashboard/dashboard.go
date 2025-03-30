package dashboard

import (
	"embed"
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/mustafa-bugra-yildiz/uphitme/repos/task"
	"github.com/mustafa-bugra-yildiz/uphitme/repos/user"
	"github.com/mustafa-bugra-yildiz/uphitme/auth"
)

//go:embed *.html
var templates embed.FS

var home *template.Template
var scheduler *template.Template
var billingUsage *template.Template
var profile *template.Template

func init() {
	funcMap := template.FuncMap{
		"httpStatusText": http.StatusText,
		"jsonMarshalIndent": func(v any) string {
			bytes, err := json.MarshalIndent(v, "", "\t")
			if err != nil {
				return err.Error()
			}
			return string(bytes)
		},
		"isActive": func(expect, got string) string {
			if expect == got {
				return "active"
			}
			return ""
		},
	}

	home = template.Must(template.New("").Funcs(funcMap).ParseFS(templates, "layout.html", "home.html"))
	scheduler = template.Must(template.New("").Funcs(funcMap).ParseFS(templates, "layout.html", "scheduler.html"))
	billingUsage = template.Must(template.New("").Funcs(funcMap).ParseFS(templates, "layout.html", "billing-usage.html"))
	profile = template.Must(template.New("").Funcs(funcMap).ParseFS(templates, "layout.html", "profile.html"))
}

type Dashboard struct {
	auth *auth.Auth
	taskRepo task.Repo
	userRepo user.Repo
}

func New(
	auth *auth.Auth,
	taskRepo task.Repo,
	userRepo user.Repo,
) *Dashboard {
	return &Dashboard{
		auth: auth,
		taskRepo: taskRepo,
		userRepo: userRepo,
	}
}

func (d *Dashboard) Home(w http.ResponseWriter, r *http.Request) {
	user, err := d.auth.User(r)
	if err != nil {
		http.Redirect(w, r, "/sign-in", http.StatusUnauthorized)
		return
	}

	data := map[string]any{
		"user": user,
		"tab": "home",
	}

	err = home.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (d *Dashboard) Scheduler(w http.ResponseWriter, r *http.Request) {
	user, err := d.auth.User(r)
	if err != nil {
		http.Redirect(w, r, "/sign-in", http.StatusUnauthorized)
		return
	}

	data := map[string]any{
		"user": user,
		"tab": "scheduler",
	}

	err = scheduler.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (d *Dashboard) BillingUsage(w http.ResponseWriter, r *http.Request) {
	user, err := d.auth.User(r)
	if err != nil {
		http.Redirect(w, r, "/sign-in", http.StatusUnauthorized)
		return
	}

	data := map[string]any{
		"user": user,
		"tab": "billing-usage",
	}

	err = billingUsage.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (d *Dashboard) Profile(w http.ResponseWriter, r *http.Request) {
	user, err := d.auth.User(r)
	if err != nil {
		http.Redirect(w, r, "/sign-in", http.StatusUnauthorized)
		return
	}

	data := map[string]any{
		"user": user,
		"tab": "profile",
	}

	err = profile.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
