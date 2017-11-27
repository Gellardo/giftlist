package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/Gellardo/giftlist/web"
)

var api listApi

type list struct {
	Id    string `json:"id"`
	Name  string `json:"name,omitempty"`
	Items []item `json:"items,omitempty"`
}

type item struct {
	Id       string `json:"id"`
	Name     string `json:"name,omitempty"`
	Link     string `json:"link,omitempty"`
	Assigned bool   `json:"assigned,omitempty"`
}

type listApi struct {
	Router *mux.Router
	Store  Storage
}

func listAPIinit(r *mux.Router, prefix string) *listApi {
	s := r.PathPrefix(prefix).Subrouter()
	s.HandleFunc("/", CreateList).Methods(http.MethodPost)
	s.HandleFunc("/{id}/", ViewList).Methods(http.MethodGet)
	s.HandleFunc("/{id}/", CreateItem).Methods(http.MethodPost)

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
		log.Print("POST ", r.URL.Path, " jsonerr: ", err)
		http.Error(w, "{\"error\":\"jsondecode\"}", http.StatusInternalServerError)
		return
	}
	i.Id = strconv.Itoa(len(l.Items)) + i.Id
	l.Items = append(l.Items, i)
	api.Store.StoreList(l)
	log.Print("POST ", r.URL.Path, " itemadded: ", i)

	w.WriteHeader(http.StatusCreated)
}

func CreateList(w http.ResponseWriter, r *http.Request) {
	var l list
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		log.Print("POST ", r.URL.Path, " jsonerr: ", err)
		http.Error(w, "{\"error\":\"jsondecode\"}", http.StatusInternalServerError)
		return
	}
	if _, err := api.Store.GetList(l.Id); err == nil {
		log.Print("POST ", r.URL.Path, " exists: ", l.Id)
		http.Error(w, "{\"error\":\"exists\"}", http.StatusInternalServerError)
		return
	}
	api.Store.StoreList(&l)
	log.Print("POST ", r.URL.Path, " listadded:", l.Id)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(list{Id: l.Id})
}

func ViewList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	l, err := api.Store.GetList(vars["id"])
	log.Print("GET ", r.URL.Path, " found:", l != nil)
	if err != nil {
		http.Error(w, "{\"error\":\"no list\"}", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(l)
}

func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)
	api = *listAPIinit(r, "/list/") //TODO feels really dirty to use global state
	api.Store.StoreList(&list{"abc", "some name", []item{}})
	web.Run(r, "/web/", "./web/", "http://localhost:8000/list/")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/web/", http.StatusMovedPermanently)
	})

	server := &http.Server{
		Addr:         ":8000",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
