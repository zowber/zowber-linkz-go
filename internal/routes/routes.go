package routes

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/zowber/zowber-linkz-go/internal/data"
	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
)

var db, err = data.NewPortfolioDbClient()

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/link/new", createHandler)
	mux.HandleFunc("/link/label/new", labelHandler)
	mux.HandleFunc("/link/edit", editHandler)
	mux.HandleFunc("/link/delete", deleteHandler)

	return mux
}

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

	switch r.Method {
	case "GET":
		{
			temp := template.Must(template.ParseFiles("./templates/create.html"))
			temp.Execute(w, nil)
		}
	case "POST":
		{
			name := r.PostFormValue("name")
			url := r.PostFormValue("url")

			// this seems a bit nasty
			formRaw := r.Form
			var labels []linkzapp.Label
			for key, value := range formRaw {
				if strings.Contains(key, "label") {
					labels = append(labels, linkzapp.Label{Id: len(labels), Name: value[0]})
				}
			}

			newLink, err := db.Insert(&linkzapp.Link{Name: name, Url: url, Labels: labels})
			if err != nil {
				errorHandler(w, r, http.StatusInternalServerError, err)
				return
			}

			temp := template.Must(template.ParseFiles("./templates/link.html"))
			temp.Execute(w, newLink)
		}
	}
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

		link := &linkzapp.Link{
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

//var labelCount = 0

var labelHandler = func(w http.ResponseWriter, r *http.Request) {
	name := r.PostFormValue("new-label")
	rawId := strings.Split(name, " ")

	var id = "label"
	for i := 0; i < len(rawId); i++ {
		id = id + "_" + rawId[i]
	}

	data := map[string]string{"name": name, "id": id}

	temp := template.Must(template.ParseFiles("./templates/label.html"))
	temp.Execute(w, data)
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
