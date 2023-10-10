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

var db, err = newPortfolioDbClient()

var errorHandler = func(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	w.WriteHeader(statusCode)
	temp := template.Must(template.ParseFiles("./templates/error.html"))
	temp.Execute(w, err)
}

var indexHandler = func(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound, err)
		return
	}

	links, err := db.All()
	if err != nil {
		log.Print(err.Error())
	}

	temp := template.Must(template.ParseFiles("./templates/index.html"))
	temp.Execute(w, links)
}

var createHandler = func(w http.ResponseWriter, r *http.Request) {

	name := r.PostFormValue("name")
	url := r.PostFormValue("url")

	log.Printf("Got name %s and url %s from form.", name, url)

	newLink, err := db.Insert(Link{Name: name, Url: url})
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err)
		return
	}

	temp := template.Must(template.ParseFiles("./templates/link.html"))
	temp.Execute(w, newLink)
}

var editHandler = func(w http.ResponseWriter, r *http.Request) {

	linkIdStr := r.URL.Query().Get("id")
	linkId, err := strconv.Atoi(linkIdStr)
	if err != nil {
		errorHandler(w, r, http.StatusBadRequest, err)
		return
	}

	switch r.Method {
	case "GET":

		linkToEdit, err := db.One(linkId)
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err)
			return
		}

		temp := template.Must(template.ParseFiles("./templates/edit.html"))
		temp.Execute(w, linkToEdit)

	case "PUT":
		name := r.PostFormValue("name")
		url := r.PostFormValue("url")

		link := Link{
			Id:   linkId,
			Name: name,
			Url:  url,
		}

		updatedLink, err := db.Update(link)
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err)
			return
		}

		temp := template.Must(template.ParseFiles("./templates/link.html"))
		temp.Execute(w, updatedLink)
	}
}

var deleteHandler = func(w http.ResponseWriter, r *http.Request) {
	linkIdStr := r.URL.Query().Get("id")
	linkId, err := strconv.Atoi(linkIdStr)
	if err != nil {
		errorHandler(w, r, http.StatusBadRequest, err)
		return
	}

	err = db.Delete(linkId)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err)
		return
	}

}

func main() {
	log.Println("Sup")

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/link/new", createHandler)
	http.HandleFunc("/link/edit", editHandler)
	http.HandleFunc("/link/delete", deleteHandler)

	http.ListenAndServe(":3000", nil)
}
