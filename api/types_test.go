package api

import (
	"reflect"
	"testing"
)

func TestListGetItem(t *testing.T) {
	i := Item{ID: "id", Name: "name", Link: "link", Assigned: true}
	l := List{Items: []Item{i}}
	geti, err := l.getItem("id")
	if err != nil || !reflect.DeepEqual(i, *geti) {
		t.Error("should return the same object")
		t.Logf("expect: %v\nresult: %v", i, *geti)
	}
	if _, err = l.getItem("nonexistant"); err == nil {
		t.Error("should error if item does not exists")
	}
}
func TestListAddItem(t *testing.T) {
	var l List
	i := Item{ID: "id", Name: "name"}
	l.addItem(i)
	if geti, err := l.getItem("id"); err != nil || !reflect.DeepEqual(i, *geti) {
		t.Error("insert+get should not change the object")
		t.Logf("expect: %v\nresult: %v", i, *geti)
	}

	if err := l.addItem(i); err == nil {
		t.Error("double insert of item with the same id should not succeed")
	}

}
