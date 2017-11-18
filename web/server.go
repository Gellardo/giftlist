package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Run(r *mux.Router, prefix string) {
	s := r.PathPrefix(prefix).Subrouter()
	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "hello world", http.StatusOK)
	})
	fmt.Println("vim-go")
}
