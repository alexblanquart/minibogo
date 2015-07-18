package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/alexblanquart/minibo/dater"
	"github.com/russross/blackfriday"
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
var categories = []Category{
	{Title: "Univers jouets enfants", Image: "coin_jouets_poupees.png"},
	{Title: "Mes projets en cours", Image: "projets_en_cours.png"},
	{Title: "Aux pinceaux!", Image: "aux_pinceaux.png"},
	{Title: "Mes plaids tout doux", Image: "patc_quilt_plaid.png"},
	{Title: "Mon coin couture", Image: "couture1.png"},
	{Title: "Mes petits objets en carton", Image: "carton_categ.png"},
}

type Post struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Image   string   `json:"image"`
	Text    string   `json:"text"`
	Date    string   `json:"date"`
	Tags    []string `json:"tags"`
	Content []byte
}

var baseTemplate = getBaseTemplate()
var indexTempl = getTemplate("templates/index.html", "templates/categories.html", "templates/news.html")
var blogTempl = getTemplate("templates/blog.html")
var postTempl = getTemplate("templates/post.html")
var aboutTempl = getTemplate("templates/about.html")
var contactTempl = getTemplate("templates/contact.html")

// Return the complete list of current funcs used through all the templates
func getAllFuncs() template.FuncMap {
	return template.FuncMap{"markDown": markDowner, "friendlyDater": dater.FriendlyDater, "thumbnail": thumbnailer}
}

// Return the base template presently used to compute all templates being executed
func getBaseTemplate() *template.Template {
	return template.Must(template.New("base").Funcs(getAllFuncs()).ParseFiles("templates/base.html",
		"templates/header.html", "templates/navigation.html", "templates/footer.html"))
}

// Add specified templates to the base template to create the final template to be executed later
func getTemplate(filenames ...string) *template.Template {
	return template.Must(template.Must(baseTemplate.Clone()).ParseFiles(filenames...))
}

// From a path, try to find the thumbnail associated image in the special directory
func thumbnailer(path string) string {
	name := filepath.Base(path)
	ext := filepath.Ext(path)
	nameWithoutExt := name[:len(name)-len(ext)]
	newPath := "static/images/thumbs/" + nameWithoutExt + ".png"
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		return "holder.js/340x340"
	} else {
		return "/" + newPath
	}
}

// Transform content in markdown into html.
func markDowner(content []byte) template.HTML {
	s := blackfriday.MarkdownCommon(content)
	return template.HTML(s)
}

func renderTemplate(w http.ResponseWriter, tmpl *template.Template, data interface{}) {
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, indexTempl, IndexContent{Categories: categories})
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, contactTempl, nil)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	content, err := getContent("about.md")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, aboutTempl, content)
}

func getContent(name string) ([]byte, error) {
	return ioutil.ReadFile("content/" + name)
}

type BlogContent struct {
	Posts []Post
	Tags  []string
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := getPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var tags = []string{"tous"}
	var setOfTags = map[string]bool{}
	for _, p := range posts {
		for _, t := range p.Tags {
			if !setOfTags[t] {
				tags = append(tags, t)
				setOfTags[t] = true
			}
		}

	}
	tag := r.URL.Path[6:] // from pattern "/blog/{{tag}}"
	var filtered = []Post{}
	if tag != "" && tag != "tous" {
		for _, p := range posts {
			for _, t := range p.Tags {
				if tag == t {
					filtered = append(filtered, p)
					continue
				}
			}
		}
	} else {
		filtered = posts
	}

	renderTemplate(w, blogTempl, BlogContent{Posts: filtered, Tags: tags})
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
	renderTemplate(w, postTempl, post)
}

func main() {
	// Serve static content
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routing
	http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/blog/", blogHandler)
	http.HandleFunc("/post/", postHandler)
	http.HandleFunc("/", indexHandler)

	// Start server
	http.ListenAndServe(":8080", nil)
}
