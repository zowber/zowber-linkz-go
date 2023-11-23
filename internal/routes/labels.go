package routes

import (
	"net/http"
)

var labelsHandler = func(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement /labels
	errorHandler(w, r, http.StatusNotImplemented, nil)
}

