package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListAPI(t *testing.T) {
	api := listAPIinit()

	req, _ := http.NewRequest("POST", "/", strings.NewReader("{\"id\":\"testid\", \"name\":\"test\"}"))
	w := httptest.NewRecorder()
	api.Router.ServeHTTP(w, req)
	var l list
	if err := json.NewDecoder(w.Body).Decode(&l); err != nil || w.Code != http.StatusCreated || l.Id != "testid" {
		t.Errorf("create: jerr='%s' Status=%d Body='%s'", err, w.Code, w.Body)
	}

	req, _ = http.NewRequest("GET", "/"+l.Id+"/", nil)
	w = httptest.NewRecorder()
	api.Router.ServeHTTP(w, req)
	if w.Code == http.StatusOK {
		if err := json.NewDecoder(w.Body).Decode(&l); err != nil || l.Id != "testid" || l.Name != "test" {
			fmt.Println(w.Code, http.StatusCreated)
			t.Errorf("create: jerr='%s' Status=%d Body='%s'", err, w.Code, w.Body)
		}
	} else {
		t.Errorf("create: Status=%d Body='%s'", w.Code, w.Body)
	}

	req, _ = http.NewRequest("POST", "/"+l.Id+"/", strings.NewReader("{\"name\":\"testitem\"}"))
	w = httptest.NewRecorder()
	api.Router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("create: Status=%d Body='%s'", w.Code, w.Body)
	}

	req, _ = http.NewRequest("GET", "/"+l.Id+"/", nil)
	w = httptest.NewRecorder()
	api.Router.ServeHTTP(w, req)
	if w.Code == http.StatusOK {
		if err := json.NewDecoder(w.Body).Decode(&l); err != nil || l.Id != "testid" || l.Name != "test" ||
			len(l.Items) != 1 || l.Items[0].Name != "testitem" {
			fmt.Println(w.Code, http.StatusCreated)
			t.Errorf("create: jerr='%s' Status=%d Body='%s'", err, w.Code, w.Body)
		}
	} else {
		t.Errorf("create: Status=%d Body='%s'", w.Code, w.Body)
	}
}
