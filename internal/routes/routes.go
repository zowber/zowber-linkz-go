package routes

import (
	"encoding/csv"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/zowber/zowber-linkz-go/internal/data/sqlite"
	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
)

var db, err = sqlite.NewDbClient()

func idToStr(id int) string {
	idStr := strconv.Itoa(id)
	return idStr
}

var funcMap = template.FuncMap{
	"idToStr": idToStr,
}

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/links", linksHandler)
	mux.HandleFunc("/link/edit", editHandler)
	mux.HandleFunc("/link", linkHandler)
	mux.HandleFunc("/labels", labelsHandler)
	mux.HandleFunc("/label", labelHandler)
	// /label/:id/links
	mux.HandleFunc("/scripts/links.js", staticHandler)

	mux.HandleFunc("/import", importHandler)

	return mux
}

var staticHandler = func(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/scripts/links.js")
}

var errorHandler = func(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	w.WriteHeader(statusCode)
	tmpl := template.Must(template.ParseFiles("./templates/error.html"))
	tmpl.Execute(w, err)
}

var indexHandler = func(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound, err)
		return
	}
	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	tmpl.Execute(w, err)
}

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

			tmpl := template.Must(template.New("links.html").Funcs(funcMap).ParseFiles("./templates/links.html", "./templates/header.html", "./templates/links-list.html", "./templates/link.html", "./templates/footer.html"))
			tmpl.Execute(w, links)
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

var labelsHandler = func(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement /labels
	errorHandler(w, r, http.StatusNotImplemented, nil)
}

var labelHandler = func(w http.ResponseWriter, r *http.Request) {
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

var importHandler = func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl := template.Must(template.New("import.html").ParseFiles("./templates/header.html", "./templates/import.html", "./templates/footer.html"))
		tmpl.ExecuteTemplate(w, "import", nil)
	case "POST":
		log.Println("post to importHandler")

		file, _, err := r.FormFile("file")
		if err != nil {
			log.Println("Err reading file", err)
			return
		}

		// parse the csv
		file.Seek(0, 0)
		reader := csv.NewReader(file)
		records, err := reader.ReadAll()

		var links []*linkzapp.Link

		for i, record := range records {
			counter := i
			var link linkzapp.Link
			for j, val := range record {
				if j == 0 {
					link.Name = val
				}
				if j == 1 {
					link.Url = val
				}
				if j >= 2 && j < len(record) {
					link.Labels = append(link.Labels, linkzapp.Label{Name: val})
				}
			}
			link.Id = &counter
			link.CreatedAt = int(time.Now().Unix())
			links = append(links, &link)
		}

		log.Println("links:", links)

		action := r.MultipartForm.Value["action"][0]

		if action == "preview" {
			tmpl := template.Must(template.New("links-list.html").Funcs(funcMap).ParseFiles("./templates/links-list.html", "./templates/link.html"))
			tmpl.ExecuteTemplate(w, "links-list", links)
		}

		if action == "import" {
			log.Println("import")
            for _, link := range links {
                log.Println("inserting:", link)
                _, err := db.Insert(link)
                if err != nil {
                    log.Println("Err inserting imported link", err)
                }
            }

            tmpl := template.Must(template.New("import-result").Parse("<p>Imported {{ . }} links.</p>"))
            tmpl.Execute(w, len(links))
		}

	default:
		errorHandler(w, r, http.StatusMethodNotAllowed, err)
		return
	}
}
