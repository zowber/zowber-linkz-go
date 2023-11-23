package routes

import (
	"net/http"
)

var staticHandler = func(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/scripts/links.js")
}
