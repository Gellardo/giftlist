package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestListAPI(t *testing.T) {

	testcases := []struct {
		name    string
		method  string
		path    string
		body    string
		expCode int
		expBody string
	}{
		{"view list", "GET", "/testid/", "", http.StatusOK, "{\"id\":\"testid\",\"name\":\"test\"}"},
		{"create list", "POST", "/", "{\"id\":\"someid\", \"name\":\"test123\"}", http.StatusCreated, "{\"id\":\"someid\"}"},
		{"view list", "GET", "/someid/", "", http.StatusOK, "{\"id\":\"someid\",\"name\":\"test123\"}"},
		{"add item", "POST", "/testid/", "{\"name\":\"testitem\"}", http.StatusCreated, ""},
		{"view list+item", "GET", "/testid/", "", http.StatusOK, "{\"id\":\"testid\",\"name\":\"test\",\"items\":[{\"id\":\"0\",\"name\":\"testitem\"}]}"},
		{"error in json", "POST", "/", "{\"id\":\"testid}", http.StatusBadRequest, ""},
		{"nonexistent list", "GET", "/testid123/", "", http.StatusNotFound, ""},
	}

	r := mux.NewRouter()
	Setup(r, "/")

	store.StoreList(&List{ID: "testid", Name: "test"})

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			if err != nil {
				t.Error("failed to construct Request")
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			var l List
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

func TestItemUpdate(t *testing.T) {
	testcases := []struct {
		name       string
		path       string
		body       string
		expStatus  int
		listbefore *List
		listafter  *List
	}{
		// successful cases
		{"update status", "/lid/items/0", "{'assigned':true}", 200,
			&List{ID: "lid", Items: []Item{{ID: "0", Assigned: false}}},
			&List{ID: "lid", Items: []Item{{ID: "0", Assigned: true}}}},
		{"update name", "/lid/items/0", "{'name':'test'}", 200,
			&List{ID: "lid", Items: []Item{{ID: "0", Name: "name"}}},
			&List{ID: "lid", Items: []Item{{ID: "0", Name: "test"}}}},
		{"update link", "/lid/items/0", "{'link':'google.de'}", 200,
			&List{ID: "lid", Items: []Item{{ID: "0", Link: "amazon.com"}}},
			&List{ID: "lid", Items: []Item{{ID: "0", Link: "google.de"}}}},
		{"update all", "/lid/items/0", "{'name':'test','link':'google.de','assigned':true}", 200,
			&List{ID: "lid", Items: []Item{{ID: "0", Name: "name", Link: "amazon.com", Assigned: false}}},
			&List{ID: "lid", Items: []Item{{ID: "0", Name: "test", Link: "google.de", Assigned: true}}}},

		// failure cases
		{"update itemid", "/lid/items/0", "{'id':'test'}", http.StatusBadRequest,
			&List{ID: "lid", Items: []Item{{ID: "0"}}}, nil},
		{"complex update itemid", "/lid/items/0", "{'id':'1','name':'test','link':'google.de','assigned':true}", http.StatusBadRequest,
			&List{ID: "lid", Items: []Item{{ID: "0"}}}, nil},
		{"jsonerr", "/lid/items/0", "{'id':'1}", http.StatusBadRequest,
			&List{ID: "lid", Items: []Item{{ID: "0"}}}, nil},
		{"nonexistant item", "/lid/items/1", "{'name':'test'}", http.StatusNotFound,
			&List{ID: "lid", Items: []Item{{ID: "0"}}}, nil},
	}

	r := mux.NewRouter()
	Setup(r, "/")
	estore, _ := store.(*easyStore)

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, tc.path, strings.NewReader(tc.body))
			if err != nil {
				t.Error("failed to construct Request")
			}

			estore.empty()
			estore.StoreList(tc.listbefore)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			//test for correct result
			if w.Result().StatusCode != tc.expStatus {
				t.Errorf("status not as expected:\n%d\n%d", w.Result().StatusCode, tc.expStatus)
			}
			if storedl, err := estore.GetList(tc.listbefore.ID); err != nil {
			} else if tc.listafter == nil {
				if !reflect.DeepEqual(tc.listbefore, storedl) {
					t.Errorf("unexpected change:\n%v\n%v", tc.listbefore, storedl)
				}
			} else if !reflect.DeepEqual(tc.listafter, storedl) {
				t.Errorf("result not as expected:\n%v\n%v", tc.listafter, storedl)
			}
		})
	}
}
