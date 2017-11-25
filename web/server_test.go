package web

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestTemplates(t *testing.T) {
	_, err := template.ParseGlob("templates/*")
	if err != nil {
		t.Error(err)
	}
}

func TestStatic(t *testing.T) {
	r := mux.NewRouter()
	Web(r, "", "static")
	ts := httptest.NewServer(r)

	http.Get(ts.URL + "/static/main.css")
}
