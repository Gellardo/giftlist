package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var lists map[string]*list = make(map[string]*list)

type list struct {
	Id    string `json:"id,omitempty"`
	Items string `json:"items,omitempty"`
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if r.Method == "GET" {
		l := lists[vars["id"]]
		log.Print("GET /", vars["id"], "/ found:", l == nil)
		if l == nil {
			l = &list{}
		}

		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		json.NewEncoder(w).Encode(l)
	} else if r.Method == "POST" {
		var l list
		if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
			log.Print("POST /", vars["id"], "/ json decode error ", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		lists[l.Id] = &l
		log.Print("POST /", vars["id"], "/ added:", l.Id)

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{id}/", ListHandler)
	lists["abc"] = &list{"abc", "some items"}

	log.Fatal(http.ListenAndServe(":8000", r))
}
