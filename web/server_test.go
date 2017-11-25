package web

import (
	"html/template"
	"testing"
)

func TestTemplates(t *testing.T) {
	_, err := template.ParseGlob("templates/*")
	if err != nil {
		t.Error(err)
	}
}
