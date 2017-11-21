package service

import (
	"html/template"
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type indexBinding struct {
	Token string
}

var indexTmpl *template.Template

func indexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		token := uuid.NewV4().String()
		// FIXME: store token to session
		// ...
		w.WriteHeader(200)
		if err := indexTmpl.Execute(w, indexBinding{Token: token}); err != nil {
			panic(err)
		}
	}
}

type tableBinding struct {
	Header   map[string][]string
	Username string
	Password string
	Token    string
}

var tableTmpl *template.Template

func tableHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			panic(err)
		}
		// FIXME: validate token from session
		// ...
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)
		err := tableTmpl.Execute(w, tableBinding{
			Header:   req.Header,
			Username: req.PostForm.Get("username"),
			Password: req.PostForm.Get("password"),
			Token:    req.PostForm.Get("token"),
		})
		if err != nil {
			panic(err)
		}
	}
}

func init() {
	indexTmpl = template.Must(template.New("index.html").ParseFiles("templates/index.html"))

	funcMap := template.FuncMap{"join": strings.Join}
	tableTmpl = template.Must(template.New("table.html").Funcs(funcMap).ParseFiles("templates/table.html"))
}
