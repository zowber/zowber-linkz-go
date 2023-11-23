package routes

import (
    "html/template"
    "net/http"
)

var indexHandler = func(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound, err)
		return
	}
	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	tmpl.Execute(w, err)
}
