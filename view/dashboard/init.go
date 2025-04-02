package dashboard

import (
	"embed"
	"html/template"
	"io"

	"github.com/mustafa-bugra-yildiz/uphitme/repos/task"
	"github.com/mustafa-bugra-yildiz/uphitme/repos/user"
)

//go:embed *.html
var htmlFiles embed.FS

type Dashboard struct {
	home         *template.Template
	scheduler    *template.Template
	billingUsage *template.Template
	profile      *template.Template
}

func New(t *template.Template) *Dashboard {
	t = template.Must(t.Clone())

	// parse layout
	t = template.Must(t.ParseFS(htmlFiles, "layout.html"))

	// parse each page
	home := template.Must(t.Clone())
	home = template.Must(home.ParseFS(htmlFiles, "home.html"))

	scheduler := template.Must(t.Clone())
	scheduler = template.Must(scheduler.ParseFS(htmlFiles, "scheduler.html"))

	billingUsage := template.Must(t.Clone())
	billingUsage = template.Must(billingUsage.ParseFS(htmlFiles, "billing-usage.html"))

	profile := template.Must(t.Clone())
	profile = template.Must(profile.ParseFS(htmlFiles, "profile.html"))

	// done
	return &Dashboard{
		home:         home,
		scheduler:    scheduler,
		billingUsage: billingUsage,
		profile:      profile,
	}
}

func (d *Dashboard) Home(w io.Writer, user *user.User) error {
	return d.home.ExecuteTemplate(w, "layout", map[string]any{
		"user": user,
		"tab":  "home",
	})
}

func (d *Dashboard) Scheduler(w io.Writer, user *user.User, tasks []task.Task) error {
	return d.scheduler.ExecuteTemplate(w, "layout", map[string]any{
		"user":  user,
		"tab":   "scheduler",
		"tasks": tasks,
	})
}

func (d *Dashboard) BillingUsage(w io.Writer, user *user.User) error {
	return d.billingUsage.ExecuteTemplate(w, "layout", map[string]any{
		"user": user,
		"tab":  "billing-usage",
	})
}

func (d *Dashboard) Profile(w io.Writer, user *user.User) error {
	return d.profile.ExecuteTemplate(w, "layout", map[string]any{
		"user": user,
		"tab":  "profile",
	})
}
