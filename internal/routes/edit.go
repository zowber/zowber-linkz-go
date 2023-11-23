package routes

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
)

var editHandler = func(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorHandler(w, r, http.StatusBadRequest, err)
		return
	}

	switch r.Method {
	case "GET":
		linkToEdit, err := db.One(id)
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err)
			return
		}

		tmpl := template.Must(template.New("edit.html").Funcs(funcMap).ParseFiles("./templates/edit.html", "./templates/label.html"))
		tmpl.Execute(w, linkToEdit)
	case "PUT":
		name := r.PostFormValue("name")
		url := r.PostFormValue("url")

		//this seems a bit less nasty than it was
		formRaw := r.Form
		var labels []linkzapp.Label
		for key, value := range formRaw {
			if strings.Contains(key, "label-") {
				labels = append(labels, linkzapp.Label{Name: value[0]})
			}
		}

		link := &linkzapp.Link{
			Name:   name,
			Url:    url,
			Labels: labels,
		}

		err := db.Update(id, link)
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err)
			return
		}

		updatedLink, err := db.One(id)
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err)
		}

		tmpl := template.Must(template.New("link.html").Funcs(funcMap).ParseFiles("./templates/link.html"))
		tmpl.ExecuteTemplate(w, "link.html", updatedLink)
	}
}


