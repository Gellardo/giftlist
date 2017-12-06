package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Setup the API request routing for list handling and the list storage.
func Setup(r *mux.Router, prefix string) {
	s := r.PathPrefix(prefix).Subrouter()
	s.HandleFunc("/", createList).Methods(http.MethodPost)
	s.HandleFunc("/{lid}/", viewList).Methods(http.MethodGet)
	s.HandleFunc("/{lid}/items/", createItem).Methods(http.MethodPost)
	s.HandleFunc("/{lid}/items/{itemid}/", updateItem).Methods(http.MethodPost)
	s.HandleFunc("/{lid}/items/{itemid}/", deleteItem).Methods(http.MethodDelete)

	store = &easyStore{make(map[string]*List)}
}
