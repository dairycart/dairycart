package main

import (
	"net/http"
)

func renderTemplates(w http.ResponseWriter, title, content string) {
	err := templates.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
