package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func DotEnv(key string) string {
	if os.Getenv(key) == "" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal(err)
		}
	}
	return os.Getenv(key)
}

type Label struct {
	Id   int
	Name string
}

type Link struct {
	Name         string
	Url          string
	Labels       []Label
	Created_date string
}

func getAllLinks() []Link {
	return linksStub
}

var indexHandler = func(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("./templates/index.html"))
	temp.Execute(w, getAllLinks())
}

func createLink(Link) Link {
	log.Println("Creating new link")
	return linkStub
}

var createLinkHandler = func(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("./templates/new.html"))
	temp.Execute(w, nil)
}

func getLink(linkId int) Link {
	log.Println("Getting link with ID", linkId)
	return linkStub
}

var getLinkHandler = func(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("./templates/link.html"))
	temp.Execute(w, nil)
}

func updateLink(linkId int, updatedLink Link) Link {
	log.Println("Updating link with ID:", linkId)
	return linkStub
}

var updateLinkHandler = func(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("./template/link.html"))
	temp.Execute(w, nil)
}

func deleteLink(linkId int) bool {
	log.Println("Deleting link with ID:", linkId)
	return true
}

var deleteLinkHandler = func(w http.ResponseWriter, r *http.Request) {
	return
}

func main() {
	log.Println("Sup")

	// API ROUTES
	// /linkz
	// - GET list_all_linkz
	// - POST create_link
	// /links/:linkId
	// - GET read_link
	// - PUT update_link
	// - DELETE delete_link

	/* 	list_all_linkz
	uses Find() to get all links
	returns all Links
	*/
	http.HandleFunc("/", indexHandler)

	/*  create_link
	takes Link
	uses save(Link) to create a new link
	returns the new Link
	*/
	http.HandleFunc("/link/new", createLinkHandler)

	/*	read_link
		takes Link.Id
		uses findOne(Link.Id) to get a single link
		returns a Link
	*/

	http.HandleFunc("/link/", getLinkHandler)

	/*	update_link
		takes Link.Id and Link
		uses findOneAndUpdate to modify an existing link
		returns the updated link
	*/
	http.HandleFunc("/link/update", updateLinkHandler)

	/*	delete_link
		takes Link.Id
		uses Remove(Link.Id)
		returns the deleted link?
	*/
	http.HandleFunc("/link/delete", deleteLinkHandler)

	http.ListenAndServe(":3000", nil)

}
