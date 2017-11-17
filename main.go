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
	Name  string `json:"name,omitempty"`
	Items []item `json:"items,omitempty"`
}

type item struct {
	Name string `json:"name"`
}

type listApi struct {
	Router *mux.Router
}

func listAPIinit() *listApi {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/", CreateList).Methods(http.MethodPost)
	r.HandleFunc("/{id}/", ViewList).Methods(http.MethodGet)
	r.HandleFunc("/{id}/", CreateItem).Methods(http.MethodPost)

	return &listApi{r}
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	l := lists[vars["id"]]

	var i item
	if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
		log.Print("POST ", r.RequestURI, " jsonerr: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	l.Items = append(l.Items, i)
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
	if lists[l.Id] != nil {
		log.Print("POST / alreadyexsisting: ", l.Id)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	lists[l.Id] = &l
	log.Print("POST / listadded:", l.Id)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(list{Id: l.Id})
}
func ViewList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	l := lists[vars["id"]]
	log.Print("GET /", vars["id"], "/ found:", l != nil)
	if l == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	json.NewEncoder(w).Encode(l)
}

func main() {
	api := listAPIinit()
	lists["abc"] = &list{"abc", "some name", []item{}}

	log.Fatal(http.ListenAndServe(":8000", api.Router))
}
