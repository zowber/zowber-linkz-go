package routes

import (
	"log"
	"net/http"
)

//

func indexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			errorHandler(w, r, http.StatusNotFound, err)
			return
		}

		// check to see if a user already exists
		userCount, err := db.CountUsers()
		if err != nil {
			log.Println("Err getting count of users", err)
		}

		log.Println("users:", userCount)

		if userCount > 0 {
            http.Redirect(w, r, "/links", http.StatusSeeOther)
		} else {
            http.Redirect(w, r, "/setup", http.StatusSeeOther)
        }

	}
}
