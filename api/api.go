package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type List struct {
	Id    string `json:"id"`
	Name  string `json:"name,omitempty"`
	Items []Item `json:"items,omitempty"`
}

type Item struct {
	Id       string `json:"id"`
	Name     string `json:"name,omitempty"`
	Link     string `json:"link,omitempty"`
	Assigned bool   `json:"assigned,omitempty"`
}

func Setup(r *mux.Router, prefix string) {
	s := r.PathPrefix(prefix).Subrouter()
	s.HandleFunc("/", CreateList).Methods(http.MethodPost)
	s.HandleFunc("/{id}/", ViewList).Methods(http.MethodGet)
	s.HandleFunc("/{id}/", CreateItem).Methods(http.MethodPost)

	store = &easyStore{make(map[string]*List)}
}
