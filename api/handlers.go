package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func createItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	l, err := store.GetList(vars["lid"])
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
	i.ID = strconv.Itoa(len(l.Items)) + i.ID
	if err := l.addItem(i); err != nil {
		log.Print("POST ", r.URL.Path, " addItem: ", err)
		http.Error(w, "{\"error\":\"exists\"}", http.StatusInternalServerError)
		return
	}
	store.StoreList(l)
	log.Print("POST ", r.URL.Path, " itemadded: ", i)

	w.WriteHeader(http.StatusCreated)
}

func createList(w http.ResponseWriter, r *http.Request) {
	var l List
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		log.Print("POST ", r.URL.Path, " jsonerr: ", err)
		http.Error(w, "{\"error\":\"jsondecode\"}", http.StatusBadRequest)
		return
	}
	if _, err := store.GetList(l.ID); err == nil {
		log.Print("POST ", r.URL.Path, " exists: ", l.ID)
		http.Error(w, "{\"error\":\"exists\"}", http.StatusInternalServerError)
		return
	}
	store.StoreList(&l)
	log.Print("POST ", r.URL.Path, " listadded:", l.ID)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(List{ID: l.ID})
}

func viewList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	l, err := store.GetList(vars["lid"])
	log.Print("GET ", r.URL.Path, " found:", l != nil)
	if err != nil {
		http.Error(w, "{\"error\":\"no list\"}", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(l)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	l, err := store.GetList(vars["lid"])
	if err != nil {
		log.Print("POST ", r.URL.Path, " updateItem failed: list not found")
		http.Error(w, "{\"error\":\"no list\"}", http.StatusNotFound)
		return
	}

	var oldItem *Item
	if oldItem, err = l.getItem(vars["itemid"]); err != nil {
		log.Print("POST ", r.URL.Path, " updateItem failed: item not found")
		http.Error(w, "{\"error\":\"no item\"}", http.StatusNotFound)
		return
	}

	var upItem Item
	if err := json.NewDecoder(r.Body).Decode(&upItem); err != nil {
		log.Print("POST ", r.URL.Path, " jsonerr: ", err)
		http.Error(w, "{\"error\":\"jsondecode\"}", http.StatusBadRequest)
		return
	} else if upItem.ID != "" {
		log.Print("POST ", r.URL.Path, " update item: has id-field")
		http.Error(w, "{\"error\":\"has ID field\"}", http.StatusBadRequest)
		return
	}

	if err := oldItem.merge(upItem); err != nil {
		log.Print("POST ", r.URL.Path, " itemmerge error: ", err)
		http.Error(w, "{\"error\":\"mergeerror\"}", http.StatusInternalServerError)
		return
	}
	l.replaceItem(oldItem)
	log.Print("POST ", r.URL.Path, " updateItem successful")
	w.WriteHeader(http.StatusOK)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	l, err := store.GetList(vars["lid"])
	if err != nil {
		log.Print("POST ", r.URL.Path, " deleteItem failed: list not found")
		http.Error(w, "{\"error\":\"no list\"}", http.StatusNotFound)
		return
	}

	var oldItem *Item
	if oldItem, err = l.getItem(vars["itemid"]); err != nil {
		log.Print("POST ", r.URL.Path, " deleteItem failed: item not found")
		http.Error(w, "{\"error\":\"no item\"}", http.StatusNotFound)
		return
	}

	l.deleteItem(oldItem.ID)
	w.WriteHeader(http.StatusOK)
}
