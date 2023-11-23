package routes

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
)

var linkHandler = func(w http.ResponseWriter, r *http.Request) {

	accepts := make(map[string]bool)
	for _, el := range strings.Split(r.Header["Accept"][0], ",") {
		accepts[el] = true
	}

	if accepts["text/html"] || accepts["*/*"] {
		log.Println("Client accepts text/html")

		switch r.Method {
		case "GET":
			log.Println("GET")

			idStr := r.URL.Query().Get("id")

			if idStr == "" {
				tmpl := template.Must(template.ParseFiles("./templates/create.html"))
				tmpl.Execute(w, nil)
			}

			if idStr != "" {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					errorHandler(w, r, http.StatusBadRequest, err)
					return
				}

				link, err := db.One(id)
				if err != nil {
					errorHandler(w, r, http.StatusInternalServerError, err)
					return
				}

				tmpl := template.Must(template.New("link.html").Funcs(funcMap).ParseFiles("./templates/link.html"))
				tmpl.ExecuteTemplate(w, "link", link)
			}
		case "POST":
			log.Println("POST")

			name := r.PostFormValue("name")
			url := r.PostFormValue("url")

			//this seems a bit nasty
			formRaw := r.Form
			var labels []linkzapp.Label
			for key, value := range formRaw {
				if strings.Contains(key, "label-") {
					log.Println("Building label with key, value[0]:", key, value[0])
					labels = append(labels, linkzapp.Label{Name: value[0]})
				}
			}

			link := &linkzapp.Link{
				Name:      name,
				Url:       url,
				Labels:    labels,
				CreatedAt: int(time.Now().Unix()),
			}

			log.Println("inserting", link)
			newLinkId, err := db.Insert(link)
			if err != nil {
				errorHandler(w, r, http.StatusInternalServerError, err)
				return
			}

			newLink, err := db.One(newLinkId)
			if err != nil {
				log.Println("Err getting inserted link", err)
			}

			tmpl := template.Must(template.New("link.html").Funcs(funcMap).ParseFiles("./templates/link.html"))
			tmpl.ExecuteTemplate(w, "link", newLink)

		case "PUT":

			log.Println("PUT")
			log.Println("using linkHandler/PUT")

			idStr := r.URL.Query().Get("id")

			if idStr == "" {

				id, err := strconv.Atoi(idStr)
				if err != nil {
					errorHandler(w, r, http.StatusBadRequest, err)
					return
				}

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

				err = db.Update(id, link)
				if err != nil {
					errorHandler(w, r, http.StatusInternalServerError, err)
					return
				}

				updatedLink, err := db.One(id)
				if err != nil {
					errorHandler(w, r, http.StatusInternalServerError, err)
				}

				tmpl := template.Must(template.New("link.html").Funcs(funcMap).ParseFiles("./templates/link.html"))
				tmpl.ExecuteTemplate(w, "link", updatedLink)
			}

			if idStr != "" {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					errorHandler(w, r, http.StatusBadRequest, err)
					return
				}

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

				err = db.Update(id, link)
				if err != nil {
					errorHandler(w, r, http.StatusInternalServerError, err)
					return
				}

				updatedLink, err := db.One(id)
				if err != nil {
					errorHandler(w, r, http.StatusInternalServerError, err)
				}

				tmpl := template.Must(template.New("link.html").Funcs(funcMap).ParseFiles("./templates/link.html"))
				tmpl.ExecuteTemplate(w, "link", updatedLink)
			}

		case "DELETE":

			log.Println("DELETE")
			log.Println("using linkHandler/DELETE")

			idStr := r.URL.Query().Get("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				errorHandler(w, r, http.StatusBadRequest, err)
				return
			}

			err = db.Delete(id)
			if err != nil {
				errorHandler(w, r, http.StatusInternalServerError, err)
				return
			}

		default:
			errorHandler(w, r, http.StatusMethodNotAllowed, err)
		}
	}
	if accepts["application/json"] {
		log.Println("Client accepts application/json")

		switch r.Method {
		case "GET":

			log.Println("GET")
			log.Println("using linkHandler/GET")
			idStr := r.URL.Query().Get("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				errorHandler(w, r, http.StatusBadRequest, err)
				return
			}

			log.Println("Id:", id)

			link, err := db.One(id)
			if err != nil {
				errorHandler(w, r, http.StatusInternalServerError, err)
				return
			}

			log.Println("link:", link)

			jsonData, err := json.Marshal(link)
			if err != nil {
				log.Println("Err masrshaling JSON", err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)
		}
	}
}
