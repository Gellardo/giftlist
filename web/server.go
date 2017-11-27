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
	Id    string `json:"id"`
	Name  string `json:"name,omitempty"`
	Items []item `json:"items,omitempty"`
}

//TODO duplicate; import if list-API is refactored into an unique package
type item struct {
	Id       string `json:"id"`
	Name     string `json:"name,omitempty"`
	Link     string `json:"link,omitempty"`
	Assigned bool   `json:"assigned,omitempty"`
}

func renderPage(w http.ResponseWriter, basedir, name string, data interface{}) error {
	t := template.Must(template.ParseGlob(basedir + "templates/*"))
	err := t.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Printf("page %s: err=%s", name, err)
	}
	return err
}

func getListHandler(basedir, listapiurl string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		log.Println(r.Method + " " + r.URL.Path + ": show " + vars["id"])
		resp, err := http.Get(listapiurl + vars["id"])
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

		renderPage(w, basedir, "show.html", struct {
			Id   string
			List list
		}{Id: vars["id"], List: l})
	}
}

// Adds a webserver for the giftlist to the router p.
// All paths must include the trailing '/'.
//
// The webserver is added to *p* under the path *prefix*.
// *basedir* points to the directory containing the static and template directory.
// *listapiurl* is the full URL of the listapi (e.g. "http://localhost:8000/api/list/")
func Run(p *mux.Router, prefix, basedir, listapiurl string) {
	if prefix[len(prefix)-1] != '/' ||
		basedir[len(basedir)-1] != '/' ||
		listapiurl[len(listapiurl)-1] != '/' {
		panic("web.Run(): prefix, basedir and listapiurl must end in '/'")
	}

	s := p.PathPrefix(prefix).Subrouter()
	s.PathPrefix("/static/").Handler(
		http.StripPrefix(prefix+"static", http.FileServer(http.Dir(basedir+"static"))),
	)
	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print(r.Method + " " + r.URL.Path)
		renderPage(w, basedir, "index.html", nil)
	})
	s.HandleFunc("/show/{id}", getListHandler(basedir, listapiurl))
}
