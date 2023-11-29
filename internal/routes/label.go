package routes

import (
	"html/template"
	"net/http"
	"time"

	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
)

func labelHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			idStr := r.URL.Query().Get("id")

			if idStr == "" {
				errorHandler(w, r, http.StatusNotImplemented, nil)
			}

			if idStr != "" {

			}
		case "POST":
			idStr := r.URL.Query().Get("id")

			if idStr == "" {
				tempId := int(time.Now().Unix())
				name := r.PostFormValue("new-label")

				label := linkzapp.Label{Id: &tempId, Name: name}

				tmpl := template.Must(template.New("label.html").ParseFiles("./templates/label.html"))
				tmpl.ExecuteTemplate(w, "label", label)
			}
			if idStr != "" {
				errorHandler(w, r, http.StatusNotImplemented, err)
			}
		}
	}
}
