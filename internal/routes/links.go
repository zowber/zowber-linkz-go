package routes

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var linksHandler = func(w http.ResponseWriter, r *http.Request) {

	accepts := make(map[string]bool)
	for _, el := range strings.Split(r.Header["Accept"][0], ",") {
		accepts[el] = true
	}

	if accepts["text/html"] || accepts["*/*"] {
		switch r.Method {
		case "GET":
			links, err := db.All()
			if err != nil {
				log.Print(err.Error())
			}

			tmpl := template.Must(template.New("links.html").Funcs(funcMap).ParseFiles("./templates/header.html", "./templates/links.html", "./templates/links-list.html", "./templates/link.html", "./templates/footer.html"))
			tmpl.ExecuteTemplate(w, "links.html", links)
		default:
			errorHandler(w, r, http.StatusMethodNotAllowed, err)
		}
	}

	if accepts["application/json"] {
		switch r.Method {
		case "GET":
			links, err := db.All()
			if err != nil {
				log.Println(err)
			}

			jsonData, err := json.Marshal(links)
			if err != nil {
				log.Println(err)
			}

			w.Header().Set("Content-Type:", "application/json")
			w.Write(jsonData)

		default:
			errorHandler(w, r, http.StatusMethodNotAllowed, err)
		}
	}
}
