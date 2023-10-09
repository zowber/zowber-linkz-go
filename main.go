package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

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

var errorHandler = func(w http.ResponseWriter, r *http.Request, statusCode int) {
	w.WriteHeader(statusCode)
	temp := template.Must(template.ParseFiles("./templates/error.html"))
	temp.Execute(w, statusCode)
}

// func getAllLinks() []Link {
// 	return linksStub
// }

var indexHandler = func(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	temp := template.Must(template.ParseFiles("./templates/index.html"))
	temp.Execute(w, linksStub)
}

// func createLink(Link) Link {
// 	log.Println("Creating new link")
// 	return linkStub
// }

var createLinkHandler = func(w http.ResponseWriter, r *http.Request) {

	name := r.PostFormValue("name")
	url := r.PostFormValue("url")

	log.Printf("Got name %s and url %s from form values.", name, url)

	newLink := Link{
		Name: name,
		Url:  url,
	}

	temp := template.Must(template.ParseFiles("./templates/link.html"))
	temp.Execute(w, newLink)
}

// func getLink(linkId int) Link {
// 	log.Println("Getting link with ID", linkId)
// 	return linkStub
// }

// var getLinkHandler = func(w http.ResponseWriter, r *http.Request) {
// 	temp := template.Must(template.ParseFiles("./templates/link.html"))
// 	temp.Execute(w, nil)
// }

// func editLink(linkId int, updatedLink Link) Link {
// 	log.Println("Updating link with ID:", linkId)

// 	return linkStub
// }

var editLinkHandler = func(w http.ResponseWriter, r *http.Request) {

	log.Println("Req to edit with verb", r.Method)

	linkIdStr := r.URL.Query().Get("id")
	linkId, err := strconv.Atoi(linkIdStr)
	if err != nil {
		errorHandler(w, r, http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		log.Println("Req to edit link with Id", linkId)

		linkToEdit := Link{
			Id:   linkId,
			Name: "Link to edit",
			Url:  "http://aintitcool.com",
		}

		temp := template.Must(template.ParseFiles("./templates/edit.html"))
		temp.Execute(w, linkToEdit)
	case "PUT":
		name := r.PostFormValue("name")
		url := r.PostFormValue("url")

		log.Printf("Got name %s and url %s from form values.", name, url)

		newLink := Link{
			Id:   linkId,
			Name: name,
			Url:  url,
		}

		temp := template.Must(template.ParseFiles("./templates/link.html"))
		temp.Execute(w, newLink)
	}
}

// func deleteLink(linkId int) bool {
// 	// Delete the link
// 	return true
// }

var deleteLinkHandler = func(w http.ResponseWriter, r *http.Request) {
	linkId := r.URL.Query().Get("id")
	log.Println("Deleting the link with Id:", linkId)
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

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/link/new", createLinkHandler)
	http.HandleFunc("/link/edit", editLinkHandler)
	http.HandleFunc("/link/delete", deleteLinkHandler)

	/*

		list_all_linkz
			uses Find() to get all links
			returns all Links

		create_link
			takes Link
			uses save(Link) to create a new link
			returns the new Link

		read_link
			takes Link.Id
			uses findOne(Link.Id) to get a single link
			returns a Link

		?? http.HandleFunc("/link/", getLinkHandler)

		update_link
			takes Link.Id and Link
			uses findOneAndUpdate to modify an existing link
			returns the updated link

		delete_link
			takes Link.Id
			uses Remove(Link.Id)
			returns the deleted link?
	*/

	http.ListenAndServe(":3000", nil)

}
