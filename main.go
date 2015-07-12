package main

import (
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"net/http"
)

type Category struct {
	Label string
	Image string
}

type Categories []Category

type IndexContent struct {
	Categories Categories
}

// Initialize some global variables
var categories = Categories{
	{Label: "Univers jouets enfants", Image: "/static/images/coin_jouets_poupees.png"},
	{Label: "Mes projets en cours", Image: "/static/images/projets_en_cours.png"},
	{Label: "Aux pinceaux!", Image: "/static/images/aux_pinceaux.png"},
	{Label: "Mes plaids tout doux", Image: "/static/images/patc_quilt_plaid.png"},
	{Label: "Mon coin couture", Image: "static/images/couture1.png"},
	{Label: "Mes petits objets en carton", Image: "static/images/carton_categ.png"},
}

// Compile all templates and cache them. Add special funcs at the end.
var templates = template.Must(template.New("main").Funcs(template.FuncMap{"markDown": markDowner}).ParseGlob("templates/*"))

func markDowner(content []byte) template.HTML {
	s := blackfriday.MarkdownCommon(content)
	return template.HTML(s)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index", IndexContent{Categories: categories})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "contact", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	markdownContent, err := ioutil.ReadFile("markdown/about.md")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := templates.ExecuteTemplate(w, "about", markdownContent); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "blog", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	// Serve static content
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routing
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/blog", blogHandler)

	// Start server
	http.ListenAndServe(":8080", nil)
}
