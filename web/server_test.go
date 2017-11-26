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

func TestBasicFiles(t *testing.T) {
	testcases := []struct {
		name string
		path string
	}{
		{"static file", "/static/main.css"},
		{"startpage", "/"},
		{"basic listview", "/show/testid"},
	}

	apis := getTestAPI()

	r := mux.NewRouter()
	Run(r, "/", "./", apis.URL+"/")
	ts := httptest.NewServer(r)

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := http.Get(ts.URL + tc.path)
			if err != nil || resp.StatusCode != http.StatusOK {
				t.Errorf("failed GET of file %s err=%s, resp.StatusCode=%d", tc.path, err, resp.StatusCode)
			}
		})
	}
}

func getTestAPI() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "{\"id\":\"testid\",\"name\":\"test\",\"items\":[]}", 200)
		}))
}
