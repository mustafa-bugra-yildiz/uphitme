package view

import (
	"embed"
	"encoding/json"
	"html/template"
	"io"
	"net/http"

	"github.com/mustafa-bugra-yildiz/uphitme/view/dashboard"
	"github.com/mustafa-bugra-yildiz/uphitme/view/docs"
)

//go:embed *
var htmlFiles embed.FS

var (
	landing *template.Template

	signIn *template.Template

	signUp        *template.Template
	signUpSuccess *template.Template
	signUpWIP     *template.Template

	Dashboard *dashboard.Dashboard
	Docs      *docs.Docs
)

func Setup(tailwind string) {
	t := template.Must(template.New("tailwind.css").Parse(tailwind))

	// include common functions
	t = t.Funcs(template.FuncMap{
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
	})

	// parse layout
	t = template.Must(t.ParseFS(htmlFiles, "layout.html"))

	// parse each page
	landing = template.Must(t.Clone())
	landing = template.Must(landing.ParseFS(htmlFiles, "landing.html"))

	signIn = template.Must(t.Clone())
	signIn = template.Must(signIn.ParseFS(htmlFiles, "sign-in.html"))

	signUp = template.Must(t.Clone())
	signUp = template.Must(signUp.ParseFS(htmlFiles, "sign-up.html"))

	signUpSuccess = template.Must(t.Clone())
	signUpSuccess = template.Must(signUpSuccess.ParseFS(htmlFiles, "sign-up-success.html"))

	signUpWIP = template.Must(t.Clone())
	signUpWIP = template.Must(signUpWIP.ParseFS(htmlFiles, "sign-up-wip.html"))

	// setup dashboard
	Dashboard = dashboard.New(t)

	// setup docs
	Docs = docs.New(t)
}

func Landing(w io.Writer) error {
	return landing.ExecuteTemplate(w, "layout", nil)
}

func SignIn(w io.Writer, email string, err error) error {
	var errMsg *string
	if err != nil {
		value := err.Error()
		errMsg = &value
	}

	return signIn.ExecuteTemplate(w, "layout", map[string]any{
		"email": email,
		"error": errMsg,
	})
}

func SignUp(w io.Writer, fullName, email string, err error) error {
	var errMsg *string
	if err != nil {
		value := err.Error()
		errMsg = &value
	}

	return signUp.ExecuteTemplate(w, "layout", map[string]any{
		"fullName": fullName,
		"email":    email,
		"error":    errMsg,
	})
}

func SignUpSuccess(w io.Writer) error {
	return signUpSuccess.ExecuteTemplate(w, "layout", nil)
}

func SignUpWIP(w io.Writer) error {
	return signUpWIP.ExecuteTemplate(w, "layout", nil)
}
