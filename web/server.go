package web

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//TODO duplicate; import if list-API is refactored into an unique package
type list struct {
	Id    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Items []item `json:"items,omitempty"`
}

//TODO duplicate; import if list-API is refactored into an unique package
type item struct {
	Name string `json:"name"`
}

func renderPage(w http.ResponseWriter, name string, data interface{}) {
	t := template.Must(template.ParseFiles("web/templates/basics.template", "web/templates/index.html", "web/templates/show.html"))
	err := t.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Printf("page %s: err=%s", name, err)
	}
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println(r.Method + " " + r.URL.Path + ": show " + vars["id"])
	resp, err := http.Get("http://localhost:8000/list/" + vars["id"])
	if err != nil {
		log.Printf("Failed REST call: err=%s", err)
		return
	} else if resp.StatusCode != http.StatusOK {
		log.Printf("REST call status=%d", resp.StatusCode)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	var l list
	err = json.NewDecoder(resp.Body).Decode(&l)
	if err != nil {
		log.Printf("Failed json decode: err=%s", err)
		return
	}
	log.Println(l)

	renderPage(w, "show.html", struct {
		Id   string
		List list
	}{Id: vars["id"], List: l})
}

func Run(r *mux.Router, prefix string) {
	s := r.PathPrefix(prefix).Subrouter()
	s.PathPrefix("/static/").Handler(
		http.StripPrefix(prefix+"static", http.FileServer(http.Dir("./web/static"))),
	)

	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print(r.Method + " " + r.URL.Path)
		renderPage(w, "index.html", nil)
	})
	s.HandleFunc("/show/{id}", ListHandler)
}
