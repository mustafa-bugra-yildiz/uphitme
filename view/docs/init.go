package docs

import (
	"embed"
	"html/template"
	"io"
)

//go:embed *.html
var htmlFiles embed.FS

type Docs struct {
	schedulerApi *template.Template
}

func New(t *template.Template) *Docs {
	t = template.Must(t.Clone())

	// parse layout
	t = template.Must(t.ParseFS(htmlFiles, "layout.html"))

	// parse each page
	schedulerApi := template.Must(t.Clone())
	schedulerApi = template.Must(schedulerApi.ParseFS(htmlFiles, "scheduler-api.html"))

	return &Docs{
		schedulerApi: schedulerApi,
	}
}

func (d *Docs) SchedulerApi(w io.Writer) error {
	return d.schedulerApi.ExecuteTemplate(w, "layout", nil)
}
