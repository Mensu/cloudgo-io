package service

import (
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {

	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	n := negroni.Classic()
	mx := mux.NewRouter()

	initRoutes(mx, formatter)

	n.Use(NewIconvHandler())
	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
	webRoot := os.Getenv("WEBROOT")
	if len(webRoot) == 0 {
		if root, err := os.Getwd(); err != nil {
			panic("Could not retrive working directory")
		} else {
			webRoot = root
		}
	}

	// static file access
	mx.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(webRoot+"/assets/"))))
	// js access
	mx.HandleFunc("/js", jsHandler(formatter)).Methods("GET")
	// index
	mx.HandleFunc("/", indexHandler()).Methods("GET")
	// post a form and output a table
	mx.HandleFunc("/table", tableHandler()).Methods("POST")
	// not implemented
	mx.NotFoundHandler = NotImplementedHandler()
}
