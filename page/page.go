package page

import (
	_ "embed"
	"encoding/json"
	"html/template"
	"io"

	"github.com/mustafa-bugra-yildiz/uphitme/repos/task"
)

//go:embed templates.html
var templates string

var tmpl = template.Must(template.New("").
	Funcs(template.FuncMap{
		"jsonMarshalIndent": func(v any) string {
			bytes, err := json.MarshalIndent(v, "", "\t")
			if err != nil {
				return err.Error()
			}
			return string(bytes)
		},
	}).
	Parse(templates),
)

func Landing(w io.Writer) error {
	return tmpl.ExecuteTemplate(w, "landing-page", nil)
}

func Dashboard(w io.Writer, tasks []task.Task) error {
	return tmpl.ExecuteTemplate(w, "dashboard-page", map[string]any{
		"tasks": tasks,
	})
}
