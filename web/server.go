package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Run(r *mux.Router, prefix string) {
	s := r.PathPrefix(prefix).Subrouter()
	s.PathPrefix("/static/").Handler(
		http.StripPrefix(prefix+"static", http.FileServer(http.Dir("./web/static"))),
	)

	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "hello world", http.StatusOK)
	})
}
