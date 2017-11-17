package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListAPI(t *testing.T) {

	testcases := []struct {
		name    string
		method  string
		url     string
		body    string
		expCode int
		expBody string
	}{
		{"view list", "GET", "/testid/", "", http.StatusOK, "{\"id\":\"testid\",\"name\":\"test\"}"},
		{"create list", "POST", "/", "{\"id\":\"someid\", \"name\":\"test123\"}", http.StatusCreated, "{\"id\":\"someid\"}"},
		{"view list", "GET", "/someid/", "", http.StatusOK, "{\"id\":\"someid\",\"name\":\"test123\"}"},
		{"error in json", "POST", "/", "{\"id\":\"testid}", http.StatusInternalServerError, ""},
		{"nonexistent list", "GET", "/testid123/", "", http.StatusNotFound, ""},
		{"add item", "POST", "/testid/", "{\"name\":\"testitem\"}", http.StatusCreated, ""},
		{"view list+item", "GET", "/testid/", "", http.StatusOK, "{\"id\":\"testid\",\"name\":\"test\",\"items\":[{\"name\":\"testitem\"}]}"},
	}

	api = *listAPIinit() //TODO feels really dirty to use the state
	api.Store.StoreList(&list{Id: "testid", Name: "test"})
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.url, strings.NewReader(tc.body))
			if err != nil {
				t.Error("failed to construct Request")
			}

			w := httptest.NewRecorder()
			api.Router.ServeHTTP(w, req)

			var l list
			if tc.expCode == w.Code && (tc.expBody == "" || tc.expBody+"\n" == w.Body.String()) {
				if err := json.NewDecoder(w.Body).Decode(&l); tc.expBody != "" && err != nil {
					t.Errorf("failed to decode answer: jerr='%s' body='%s'", err, w.Body)
				}
			} else {
				t.Errorf("did not respond as expected\n%d '%s'\n%d '%s'", tc.expCode, tc.expBody, w.Code, w.Body.String())
			}
		})
	}

}
