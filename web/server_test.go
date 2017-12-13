package web

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
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
		name           string
		path           string
		expectedStatus int
		bodyContains   string
	}{
		{"static file", "/static/main.css", http.StatusOK, ""},
		{"startpage", "/", http.StatusOK, ""},
		{"basic listview", "/show/testid", http.StatusOK, ""},
		{"basic editview first item", "/edit/testid/1", http.StatusOK, "some fancy present"},
		{"basic editview third item", "/edit/testid/2", http.StatusOK, "a lame book"},

		{"listview no list", "/show/test", http.StatusNotFound, ""},
		{"editview no item", "/edit/testid/100", http.StatusNotFound, ""},
		{"editview no list", "/edit/test/1", http.StatusNotFound, ""},
	}

	apis := getTestAPI()

	r := mux.NewRouter()
	Run(r, "/", "./", apis.URL+"/")
	ts := httptest.NewServer(r)

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := http.Get(ts.URL + tc.path)
			if err != nil || resp.StatusCode != tc.expectedStatus {
				t.Errorf("failed GET of file %s err=%v, resp.StatusCode=%d", tc.path, err, resp.StatusCode)
			}
			body, _ := ioutil.ReadAll(resp.Body)
			if tc.bodyContains != "" && !strings.Contains(string(body), tc.bodyContains) {
				t.Errorf("body does not contain '%s'\n%s", tc.bodyContains, body)
			}
		})
	}
}

func getTestAPI() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if strings.Contains(r.URL.Path, "/testid") {
				http.Error(w,
					`{"id":"testid","name":"asdf","items":[
					{"id":"0","name":"a","link":"b","assigned":true},
					{"id":"1","name":"some fancy present","link":"aaaaaaa","assigned":true},
					{"id":"2","name":"a lame book","link":"amazon"}
				]}`,
					http.StatusOK)
			} else {
				http.Error(w, "{\"error\":\"Not found\"}", http.StatusNotFound)
			}
		}))
}
