package routes

import (
	"net/http"
)

func labelsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement /labels
		errorHandler(w, r, http.StatusNotImplemented, nil)
	}
}
