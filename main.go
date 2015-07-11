package main

import (
	"html/template"
	"net/http"
)

// compile all templates and cache them
var templates = template.Must(template.ParseGlob("templates/*"))

func handler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	// Serve static content
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// For now only index
	http.HandleFunc("/", handler)

	// Start server
	http.ListenAndServe(":8080", nil)
}
