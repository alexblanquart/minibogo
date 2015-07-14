package main

import (
	"encoding/json"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
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
	{Title: "Univers jouets enfants", Image: "coin_jouets_poupees.png"},
	{Title: "Mes projets en cours", Image: "projets_en_cours.png"},
	{Title: "Aux pinceaux!", Image: "aux_pinceaux.png"},
	{Title: "Mes plaids tout doux", Image: "patc_quilt_plaid.png"},
	{Title: "Mon coin couture", Image: "couture1.png"},
	{Title: "Mes petits objets en carton", Image: "carton_categ.png"},
}

type Post struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Image   string `json:"image"`
	Text    string `json:"text"`
	Date    string `json:"date"`
	Content []byte
}

type BlogContent struct {
	Posts []Post
}

// Compile all templates and cache them. Add special pipelines.
var templates = template.Must(template.New("main").Funcs(template.FuncMap{"markDown": markDowner,
	"time": userFriendlyTimer, "thumbnail": thumbnailer}).ParseGlob("templates/*"))

// From a path, try to find the thumbnail associated image in the special directory
func thumbnailer(path string) string {
	name := filepath.Base(path)
	ext := filepath.Ext(path)
	nameWithoutExt := name[:len(name)-len(ext)]
	newPath := "static/images/thumbs/" + nameWithoutExt + ".png"
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		newPath = "holder.js/400x340"
	}
	return newPath
}

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
	content, err := getContent("about.md")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "about", content)
}

func getContent(name string) ([]byte, error) {
	return ioutil.ReadFile("content/" + name)
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := getPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "blog", posts)
}

func getPosts() ([]Post, error) {
	postsAsJson, err := ioutil.ReadFile("content/posts.json")
	if err != nil {
		return nil, err
	}
	posts := []Post{}
	if err := json.Unmarshal(postsAsJson, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[6:] // from pattern "/post/{{ID}}"
	posts, err := getPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var post = Post{}
	for _, p := range posts {
		if p.ID == id {
			post = p
		}
	}
	content, err := getContent(post.Text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	post.Content = content
	renderTemplate(w, "post", post)
}

func main() {
	// Serve static content
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routing
	http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/blog", blogHandler)
	http.HandleFunc("/post/", postHandler)
	http.HandleFunc("/", indexHandler)

	// Start server
	http.ListenAndServe(":8080", nil)
}
