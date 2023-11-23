package routes

import (
    "html/template"
    "net/http"
)

var errorHandler = func(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	w.WriteHeader(statusCode)
	tmpl := template.Must(template.ParseFiles("./templates/error.html"))
	tmpl.Execute(w, err)
}
