package main

import (
	"encoding/json"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"
)

type Category struct {
	Title string
	Image string
}

type Categories []Category

type IndexContent struct {
	Categories Categories
}

// Initialize some global variables
var categories = Categories{
	{Title: "Univers jouets enfants", Image: "/static/images/coin_jouets_poupees.png"},
	{Title: "Mes projets en cours", Image: "/static/images/projets_en_cours.png"},
	{Title: "Aux pinceaux!", Image: "/static/images/aux_pinceaux.png"},
	{Title: "Mes plaids tout doux", Image: "/static/images/patc_quilt_plaid.png"},
	{Title: "Mon coin couture", Image: "static/images/couture1.png"},
	{Title: "Mes petits objets en carton", Image: "static/images/carton_categ.png"},
}

type Post struct {
	Title   string `json:"title"`
	Image   string `json:"image"`
	Content string `json:"content"`
	Date    string `json:"date"`
}

type BlogContent struct {
	Posts []Post
}

// Compile all templates and cache them. Add special pipelines.
var templates = template.Must(template.New("main").Funcs(template.FuncMap{"markDown": markDowner, "time": userFriendlyTimer}).ParseGlob("templates/*"))

// Transform content in markdown into html.
func markDowner(content []byte) template.HTML {
	s := blackfriday.MarkdownCommon(content)
	return template.HTML(s)
}

// Transform date in from a specific layout into another one more friendly for users
func userFriendlyTimer(date string) string {
	parsed, err := time.Parse("Mon, 02 Jan 2006 15:04:05", date)
	if err == nil {
		return parsed.Format("02/01/2006") // see for internationalization later
	} else {
		return date
	}
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	err := templates.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", IndexContent{Categories: categories})
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "contact", nil)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	markdownContent, err := ioutil.ReadFile("content/about.md")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "about", markdownContent)
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	postsAsJson, err := ioutil.ReadFile("content/posts.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	posts := []Post{}
	if err := json.Unmarshal(postsAsJson, &posts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "blog", posts)
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
