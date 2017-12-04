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
		{"add item", "POST", "/testid/items/", "{\"name\":\"testitem\"}", http.StatusCreated, ""},
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
		{"update status", "/lid/items/0/", "{\"assigned\":true}", 200,
			&List{ID: "lid", Items: []Item{{ID: "0", Assigned: false}}},
			&List{ID: "lid", Items: []Item{{ID: "0", Assigned: true}}}},
		{"update name", "/lid/items/0/", "{\"name\":\"test\"}", 200,
			&List{ID: "lid", Items: []Item{{ID: "0", Name: "name"}}},
			&List{ID: "lid", Items: []Item{{ID: "0", Name: "test"}}}},
		{"update link", "/lid/items/0/", "{\"link\":\"google.de\"}", 200,
			&List{ID: "lid", Items: []Item{{ID: "0", Link: "amazon.com"}}},
			&List{ID: "lid", Items: []Item{{ID: "0", Link: "google.de"}}}},
		{"update all", "/lid/items/0/", "{\"name\":\"test\",\"link\":\"google.de\",\"assigned\":true}", 200,
			&List{ID: "lid", Items: []Item{{ID: "0", Name: "name", Link: "amazon.com", Assigned: false}}},
			&List{ID: "lid", Items: []Item{{ID: "0", Name: "test", Link: "google.de", Assigned: true}}}},

		// failure cases
		{"update itemid", "/lid/items/0/", "{\"id\":\"test\"}", http.StatusBadRequest,
			&List{ID: "lid", Items: []Item{{ID: "0"}}}, nil},
		{"complex update itemid", "/lid/items/0/", "{\"id\":\"1\",\"name\":\"test\",\"link\":\"google.de\",\"assigned\":true}", http.StatusBadRequest,
			&List{ID: "lid", Items: []Item{{ID: "0"}}}, nil},
		{"jsonerr", "/lid/items/0/", "{\"id\":\"1}", http.StatusBadRequest,
			&List{ID: "lid", Items: []Item{{ID: "0"}}}, nil},
		{"nonexistant item", "/lid/items/1/", "{\"name\":\"test\"}", http.StatusNotFound,
			&List{ID: "lid", Items: []Item{{ID: "0"}}}, nil},
		{"nonexistant list", "/listid/items/0/", "{\"name\":\"test\"}", http.StatusNotFound,
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
				t.Errorf("status not as expected:\nexpect: %d\nresult: %d", tc.expStatus, w.Result().StatusCode)
			}
			if storedl, err := estore.GetList(tc.listbefore.ID); err != nil {
				t.Errorf("list disappeared !?!: %v", err)
			} else if tc.listafter == nil {
				if !reflect.DeepEqual(tc.listbefore, storedl) {
					t.Errorf("unexpected change:\nexpect: %v\nresult: %v", tc.listbefore, storedl)
				}
			} else if !reflect.DeepEqual(tc.listafter, storedl) {
				t.Errorf("result not as expected:\nexpect: %v\nresult: %v", tc.listafter, storedl)
			}
		})
	}
}

func TestItemDelete(t *testing.T) {
	testcases := []struct {
		name       string
		path       string
		expStatus  int
		listbefore *List
		listafter  *List
	}{
		// successful cases
		{"delete from single element list", "/lid/items/0/", http.StatusOK,
			&List{ID: "lid", Items: []Item{{ID: "0"}}},
			&List{ID: "lid", Items: []Item{}}},
		{"delete 0 from multi element list", "/lid/items/0/", http.StatusOK,
			&List{ID: "lid", Items: []Item{{ID: "0"}, {ID: "1"}, {ID: "2"}}},
			&List{ID: "lid", Items: []Item{{ID: "1"}, {ID: "2"}}}},
		{"delete 1 from multi element list", "/lid/items/1/", http.StatusOK,
			&List{ID: "lid", Items: []Item{{ID: "0"}, {ID: "1"}, {ID: "2"}}},
			&List{ID: "lid", Items: []Item{{ID: "0"}, {ID: "2"}}}},
		{"delete 2 from multi element list", "/lid/items/2/", http.StatusOK,
			&List{ID: "lid", Items: []Item{{ID: "0"}, {ID: "1"}, {ID: "2"}}},
			&List{ID: "lid", Items: []Item{{ID: "0"}, {ID: "1"}}}},

		// failure cases
		{"no items", "/lid/items/1/", http.StatusNotFound,
			&List{ID: "lid", Items: []Item{{}}}, &List{ID: "lid", Items: []Item{{}}}},
		{"nonexistant item", "/lid/items/1/", http.StatusNotFound,
			&List{ID: "lid", Items: []Item{{ID: "0"}}}, &List{ID: "lid", Items: []Item{{ID: "0"}}}},
		{"nonexistant list", "/listid/items/0/", http.StatusNotFound,
			&List{ID: "lid", Items: []Item{{ID: "0"}}}, &List{ID: "lid", Items: []Item{{ID: "0"}}}},
	}

	r := mux.NewRouter()
	Setup(r, "/")
	estore, _ := store.(*easyStore)

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, tc.path, nil) //strings.NewReader(tc.body))
			if err != nil {
				t.Error("failed to construct Request")
			}

			estore.empty()
			estore.StoreList(tc.listbefore)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			//test for correct result
			if w.Result().StatusCode != tc.expStatus {
				t.Errorf("status not as expected:\nexpect: %d\nresult: %d", tc.expStatus, w.Result().StatusCode)
			}
			if storedl, err := estore.GetList(tc.listbefore.ID); err != nil {
				t.Errorf("list disappeared !?!: %v", err)
			} else if !reflect.DeepEqual(tc.listafter, storedl) {
				t.Errorf("result not as expected:\nexpect: %v\nresult: %v", tc.listafter, storedl)
			}
		})
	}
}
