package web

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Gellardo/giftlist/api"
)

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

		var l api.List
		err = json.NewDecoder(resp.Body).Decode(&l)
		if err != nil {
			log.Printf("Failed json decode: err=%s", err)
			return
		}
		log.Println(l)

		renderPage(w, basedir, "show.html", struct {
			ID   string
			List api.List
		}{ID: vars["id"], List: l})
	}
}

func getUpdateHandler(basedir, listapiurl string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		listid := vars["id"]
		itemid := vars["item"]

		log.Printf("%s %s: update page", r.Method, r.URL.Path, listid, itemid)

		resp, err := http.Get(listapiurl + listid)
		if err != nil {
			log.Printf("%s: Failed REST call err=%s", r.URL.Path, err)
			return
		} else if resp.StatusCode != http.StatusOK {
			log.Printf("%s: REST call status=%d", r.URL.Path, resp.StatusCode)
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		var l api.List
		err = json.NewDecoder(resp.Body).Decode(&l)
		if err != nil {
			log.Printf("%s: Failed json decode err=%s", r.URL.Path, err)
			return
		}
		log.Println(l)

		var item api.Item
		for _, it := range l.Items {
			if it.ID == itemid {
				item = it
				break
			}
		}
		if item.ID == "" {
			log.Printf("%s: item not found", r.URL.Path)
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		renderPage(w, basedir, "edit.html", struct {
			ListID string
			Item   api.Item
		}{ListID: listid, Item: item})
	}
}

// Run a webserver for the giftlist using the router p.
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
	s.HandleFunc("/edit/{id}/{item}", getUpdateHandler(basedir, listapiurl))
}
