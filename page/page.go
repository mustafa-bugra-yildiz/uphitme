package page

import (
	"embed"
	"encoding/json"
	"html/template"
	"io"
	"net/http"

	"github.com/mustafa-bugra-yildiz/uphitme/repos/task"
)

//go:embed *
var templates embed.FS

var tmpl = template.Must(template.New("").
	Funcs(template.FuncMap{
		"jsonMarshalIndent": func(v any) string {
			bytes, err := json.MarshalIndent(v, "", "\t")
			if err != nil {
				return err.Error()
			}
			return string(bytes)
		},
		"httpStatusText": http.StatusText,
	}).
	ParseFS(templates, "*.html"),
)

func Landing(w io.Writer) error {
	return tmpl.ExecuteTemplate(w, "landing-page", nil)
}

func Dashboard(w io.Writer, page, pageSize, count int, tasks []task.Task) error {
	pageCount := count / pageSize

	var prev *int
	if page > 1 {
		value := page - 1
		prev = &value
	}

	var next *int
	if page < pageCount {
		value := page + 1
		next = &value
	}

	return tmpl.ExecuteTemplate(w, "dashboard-page", map[string]any{
		"tasks": tasks,

		"prev": prev,
		"page": page,
		"next": next,
	})
}

func SignUp(w io.Writer, err error) error {
	var errMsg *string
	if err != nil {
		value := err.Error()
		errMsg = &value
	}

	return tmpl.ExecuteTemplate(w, "sign-up-page", map[string]any{
		"error": errMsg,
	})
}

func SignIn(w io.Writer) error {
	return tmpl.ExecuteTemplate(w, "sign-in-page", nil)
}
