package page

import (
	_ "embed"
	"html/template"
	"io"
)

//go:embed templates.html
var templates string

var tmpl = template.Must(template.New("").Parse(templates))

func Landing(w io.Writer) error {
	return tmpl.ExecuteTemplate(w, "landing-page", nil)
}
