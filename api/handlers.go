package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	l, err := store.GetList(vars["id"])
	if err != nil {
		http.Error(w, "{\"error\":\"no list\"}", http.StatusNotFound)
		return
	}

	var i Item
	if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
		log.Print("POST ", r.URL.Path, " jsonerr: ", err)
		http.Error(w, "{\"error\":\"jsondecode\"}", http.StatusInternalServerError)
		return
	}
	i.Id = strconv.Itoa(len(l.Items)) + i.Id
	l.Items = append(l.Items, i)
	store.StoreList(l)
	log.Print("POST ", r.URL.Path, " itemadded: ", i)

	w.WriteHeader(http.StatusCreated)
}

func CreateList(w http.ResponseWriter, r *http.Request) {
	var l List
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		log.Print("POST ", r.URL.Path, " jsonerr: ", err)
		http.Error(w, "{\"error\":\"jsondecode\"}", http.StatusInternalServerError)
		return
	}
	if _, err := store.GetList(l.Id); err == nil {
		log.Print("POST ", r.URL.Path, " exists: ", l.Id)
		http.Error(w, "{\"error\":\"exists\"}", http.StatusInternalServerError)
		return
	}
	store.StoreList(&l)
	log.Print("POST ", r.URL.Path, " listadded:", l.Id)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(List{Id: l.Id})
}

func ViewList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	l, err := store.GetList(vars["id"])
	log.Print("GET ", r.URL.Path, " found:", l != nil)
	if err != nil {
		http.Error(w, "{\"error\":\"no list\"}", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(l)
}
