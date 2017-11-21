package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func renderPage(w http.ResponseWriter, name string, data interface{}) {
	t := template.Must(template.ParseFiles("web/templates/basics.template", "web/templates/index.html", "web/templates/show.html"))
	t.ExecuteTemplate(w, name, data)
}

func Run(r *mux.Router, prefix string) {
	s := r.PathPrefix(prefix).Subrouter()
	s.PathPrefix("/static/").Handler(
		http.StripPrefix(prefix+"static", http.FileServer(http.Dir("./web/static"))),
	)

	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print(r.Method + " " + r.URL.Path)
		renderPage(w, "index.html", nil)
	})
	s.HandleFunc("/show/", func(w http.ResponseWriter, r *http.Request) {
		log.Print(r.Method + " " + r.URL.Path)
		renderPage(w, "show.html", nil)
	})
}
