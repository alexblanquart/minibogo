package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/alexblanquart/minibo/dater"
	"github.com/russross/blackfriday"
)

type Category struct {
	Title string
	Image string
}

type IndexContent struct {
	Categories []Category
}

type PostContent struct {
	Post Post
	MetaBlog
	CurrentURL string
}

type MetaBlog struct {
	Recent []Post
	Tags   []string
}

type BlogContent struct {
	Posts []Post
	MetaBlog
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

type Product struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Post        string `json:"post"`
}

type ProductsContent struct {
	Products []Product
}

type Tutorial struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Video string `json:"video"`
}

type TutorialsContent struct {
	Tutorials []Tutorial
}

var (
	categories    []Category
	posts, recent []Post
	tags          []string
	metaBlog      MetaBlog
	host          string = "http://localhost:8080" // TODO: to change when live!
	products      []Product
	tutorials     []Tutorial

	baseTempl, indexTempl, blogTempl, postTempl, aboutTempl, contactTempl, productsTempl, tutorialsTempl *template.Template
)

// Return the complete list of current funcs used through all the templates
func getAllFuncs() template.FuncMap {
	return template.FuncMap{"markDown": markDowner, "date": dater.FriendlyDater, "holder": holder}
}

// Return the base template presently used to compute all templates being executed
func getBaseTemplate() *template.Template {
	return template.Must(template.New("base").Funcs(getAllFuncs()).ParseFiles("templates/base.html",
		"templates/header.html", "templates/navigation.html", "templates/footer.html"))
}

// Add specified templates to the base template to create the final template to be executed later
func getTemplate(filenames ...string) *template.Template {
	return template.Must(template.Must(getBaseTemplate().Clone()).ParseFiles(filenames...))
}

func holder(path string) string {
	if _, err := os.Stat(path[1:]); os.IsNotExist(err) {
		return "holder.js/340x340"
	} else {
		return path
	}
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

// Read markdown file inside content folder
func getContent(name string) ([]byte, error) {
	return ioutil.ReadFile("content/" + name)
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	// filtered posts only corresponding to tag
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
	renderTemplate(w, blogTempl, BlogContent{Posts: filtered, MetaBlog: metaBlog})
}

// Return all posts sorted by date, by reading the temporay json file
func getPosts() ([]Post, error) {
	postsAsJson, err := ioutil.ReadFile("posts.json")
	if err != nil {
		return nil, err
	}
	posts := []Post{}
	if err := json.Unmarshal(postsAsJson, &posts); err != nil {
		return nil, err
	}
	sort.Sort(ByDate(posts))
	return posts, nil
}

type ByDate []Post

func (d ByDate) Len() int           { return len(d) }
func (d ByDate) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d ByDate) Less(i, j int) bool { return dater.Parse(d[i].Date).After(dater.Parse(d[j].Date)) }

func postHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[6:] // from pattern "/post/{{ID}}"
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
	currentURL := host + r.URL.Path
	renderTemplate(w, postTempl, PostContent{Post: post, MetaBlog: metaBlog, CurrentURL: currentURL})
}

// Return all products, by reading the temporay json file
func getProducts() ([]Product, error) {
	productsAsJson, err := ioutil.ReadFile("products.json")
	if err != nil {
		return nil, err
	}
	products := []Product{}
	if err := json.Unmarshal(productsAsJson, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, productsTempl, ProductsContent{Products: products})
}

func tutorialsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, tutorialsTempl, TutorialsContent{Tutorials: tutorials})
}

// Return all tutorials, by reading the temporay json file
func getTutorials() ([]Tutorial, error) {
	tutorialsAsJson, err := ioutil.ReadFile("tutorials.json")
	if err != nil {
		return nil, err
	}
	tutorials := []Tutorial{}
	if err := json.Unmarshal(tutorialsAsJson, &tutorials); err != nil {
		return nil, err
	}
	return tutorials, nil
}

// For now there is no database!
func init() {
	// main categories
	categories = []Category{
		{Title: "Univers jouets enfants", Image: "coin_jouets_poupees.png"},
		{Title: "Mes projets en cours", Image: "projets_en_cours.png"},
		{Title: "Aux pinceaux!", Image: "aux_pinceaux.png"},
		{Title: "Mes plaids tout doux", Image: "patc_quilt_plaid.png"},
		{Title: "Mon coin couture", Image: "couture1.png"},
		{Title: "Mes petits objets en carton", Image: "carton_categ.png"},
	}
	// all posts sorted by date!
	posts, _ = getPosts()
	// get 5 mosts recent posts
	recent = posts[:5]
	// gather tags
	tags = []string{"tous"}
	var setOfTags = map[string]bool{}
	for _, p := range posts {
		for _, t := range p.Tags {
			if !setOfTags[t] {
				tags = append(tags, t)
				setOfTags[t] = true
			}
		}
	}
	// recent + tags
	metaBlog = MetaBlog{Recent: recent, Tags: tags}
	// all products
	products, _ = getProducts()
	// all tutorials
	tutorials, _ = getTutorials()

	// templates
	indexTempl = getTemplate("templates/index.html", "templates/categories.html", "templates/news.html")
	blogTempl = getTemplate("templates/blog.html", "templates/sidebar.html")
	postTempl = getTemplate("templates/post.html", "templates/sidebar.html")
	aboutTempl = getTemplate("templates/about.html")
	contactTempl = getTemplate("templates/contact.html")
	productsTempl = getTemplate("templates/products.html")
	tutorialsTempl = getTemplate("templates/tutorials.html")
}

func main() {
	// Serve static content
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routing
	http.HandleFunc("/tutorials", tutorialsHandler)
	http.HandleFunc("/products", productsHandler)
	http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/blog/", blogHandler)
	http.HandleFunc("/post/", postHandler)
	http.HandleFunc("/", indexHandler)

	// Start server
	http.ListenAndServe(":8080", nil)
}
