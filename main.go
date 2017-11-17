package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var api listApi

type list struct {
	Id    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Items []item `json:"items,omitempty"`
}

type item struct {
	Name string `json:"name"`
}

type listApi struct {
	Router *mux.Router
	Store  Storage
}

func listAPIinit() *listApi {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/", CreateList).Methods(http.MethodPost)
	r.HandleFunc("/{id}/", ViewList).Methods(http.MethodGet)
	r.HandleFunc("/{id}/", CreateItem).Methods(http.MethodPost)

	return &listApi{r, &easyStore{make(map[string]*list)}}
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	l, err := api.Store.GetList(vars["id"])
	if err != nil {
		http.Error(w, "{\"error\":\"no list\"}", http.StatusNotFound)
		return
	}

	var i item
	if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
		log.Print("POST ", r.RequestURI, " jsonerr: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	l.Items = append(l.Items, i)
	api.Store.StoreList(l)
	log.Print("POST ", r.RequestURI, " itemadded: ", i)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
}

func CreateList(w http.ResponseWriter, r *http.Request) {
	var l list
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		log.Print("POST / jsonerr: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err := api.Store.GetList(l.Id); err == nil {
		log.Print("POST / exists: ", l.Id)
		http.Error(w, "{\"error\":\"exists\"}", http.StatusInternalServerError)
		return
	}
	api.Store.StoreList(&l)
	log.Print("POST / listadded:", l.Id)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(list{Id: l.Id})
}

func ViewList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	l, err := api.Store.GetList(vars["id"])
	log.Print("GET /", vars["id"], "/ found:", l != nil)
	if err != nil {
		http.Error(w, "{\"error\":\"no list\"}", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(l)
}

func main() {
	api = *listAPIinit() //TODO feels really dirty to use the state
	api.Store.StoreList(&list{"abc", "some name", []item{}})

	log.Fatal(http.ListenAndServe(":8000", api.Router))
}
