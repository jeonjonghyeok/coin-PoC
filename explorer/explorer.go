package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/jeonjonghyeok/coin/blockchain"
)

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

const (
	templateDir string = "explorer/templates/"
)

func home(w http.ResponseWriter, r *http.Request) {
	data := homeData{"Home", nil}
	templates.ExecuteTemplate(w, "home", data)
}
func add(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(w, "add", nil)
	case "POST":
		r.ParseForm()
		blockchain.Blockchain().AddBlock()
		http.Redirect(w, r, "/", http.StatusMovedPermanently)

	}
}

func Start(port int) {
	handler := http.NewServeMux()
	blockchain.Blockchain()

	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	handler.HandleFunc("/", home)
	handler.HandleFunc("/add", add)

	fmt.Printf("Listening on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}
