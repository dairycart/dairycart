package main

import "net/http"

// PageResponse represents an HTML page response to a request
type PageResponse struct {
	Title   string
	Content string
}

func renderTemplates(w http.ResponseWriter, title, content string) {
	err := templates.ExecuteTemplate(w, "indexPage", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
