package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// List stores the information of a giftlist.
type List struct {
	ID    string `json:"id"`
	Name  string `json:"name,omitempty"`
	Items []Item `json:"items,omitempty"`
}

// Item represents a single entry in a giftlist.
type Item struct {
	ID       string `json:"id"`
	Name     string `json:"name,omitempty"`
	Link     string `json:"link,omitempty"`
	Assigned bool   `json:"assigned,omitempty"`
}

// Setup the API request routing for list handling and the list storage.
func Setup(r *mux.Router, prefix string) {
	s := r.PathPrefix(prefix).Subrouter()
	s.HandleFunc("/", createList).Methods(http.MethodPost)
	s.HandleFunc("/{lid}/", viewList).Methods(http.MethodGet)
	s.HandleFunc("/{lid}/items/", createItem).Methods(http.MethodPost)
	s.HandleFunc("/{lid}/items/{itemid}/", updateItem).Methods(http.MethodPost)

	store = &easyStore{make(map[string]*List)}
}
