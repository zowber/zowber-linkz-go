package routes

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
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
	mux.HandleFunc("/link", linkHandler)
	mux.HandleFunc("/link/create-placeholder", createPlaceholderHandler)
	mux.HandleFunc("/link/new", createHandler)
	mux.HandleFunc("/link/label/new", labelHandler)
	mux.HandleFunc("/link/edit", editHandler)
	mux.HandleFunc("/link/delete", deleteHandler)

	return mux
}

var errorHandler = func(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	w.WriteHeader(statusCode)
	tmpl := template.Must(template.ParseFiles("./templates/error.html"))
	tmpl.Execute(w, err)
}

var createPlaceholderHandler = func(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/create-placeholder.html"))
	tmpl.ExecuteTemplate(w, "create-placeholder", nil)
}

var indexHandler = func(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound, err)
		return
	}

	links, err := db.All()
	if err != nil {
		log.Print(err.Error())
	}

	tmpl := template.Must(template.New("index.html").Funcs(funcMap).ParseFiles("./templates/index.html", "./templates/header.html", "./templates/create-placeholder.html", "./templates/links.html", "./templates/link.html", "./templates/footer.html"))
	tmpl.Execute(w, links)
}

var createHandler = func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		{
			tmpl := template.Must(template.ParseFiles("./templates/create.html"))
			tmpl.Execute(w, nil)
		}
	case "POST":
		{
			name := r.PostFormValue("name")
			url := r.PostFormValue("url")

			//this seems a bit nasty
			formRaw := r.Form
			var labels []linkzapp.Label
			for key, value := range formRaw {
				if strings.Contains(key, "label_") {
					keyToI, _ := strconv.Atoi(strings.SplitAfter(key, "_")[0])
					labels = append(labels, linkzapp.Label{Id: keyToI, Name: value[0]})
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
            log.Println("inserted. id", newLinkId)

            newLink, err := db.One(newLinkId)
            if err != nil {
                log.Println("Err getting inserted link", err)
            }

            tmpl := template.Must(template.New("link.html").Funcs(funcMap).ParseFiles("./templates/link.html"))
			tmpl.ExecuteTemplate(w, "link", newLink)
		}
	}
}

var linkHandler = func(w http.ResponseWriter, r *http.Request) {
	// oidStr := r.URL.Path[len("/link/"):]
	// oidStr := r.URL.Query().Get("id")
	// oid, err := primitive.ObjectIDFromHex(oidStr)
	// if err != nil {
	// 	errorHandler(w, r, http.StatusBadRequest, err)
	// 	return
	// }

	idStr := r.URL.Query().Get("id")
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

		//this seems a bit nasty
		formRaw := r.Form
		var labels []linkzapp.Label
		for key, value := range formRaw {
			if strings.Contains(key, "label_") {
				keyToI, _ := strconv.Atoi(strings.SplitAfter(key, "_")[0])
				labels = append(labels, linkzapp.Label{Id: keyToI, Name: value[0]})
			}
		}

 		link := &linkzapp.Link{
 			Name:   name,
 			Url:    url,
 			Labels: labels,
 		}

 		updatedLink, err := db.Update(id, link)
 		if err != nil {
 			errorHandler(w, r, http.StatusInternalServerError, err)
 			return
 		}

 		tmpl := template.Must(template.New("link.html").Funcs(funcMap).ParseFiles("./templates/link.html"))
 		tmpl.ExecuteTemplate(w, "link", updatedLink)
    }
}

var labelHandler = func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		name := r.PostFormValue("new-label")
		rawId := strings.Split(name, " ")

		id := "label_"
		for i := 0; i < len(rawId); i++ {
			id = id + rawId[i]
		}

		data := map[string]string{"Id": id, "Name": name}
		tmpl := template.Must(template.New("label.html").ParseFiles("./templates/label.html"))
		tmpl.ExecuteTemplate(w, "label", data)
	}
}

var deleteHandler = func(w http.ResponseWriter, r *http.Request) {
	//oidStr := r.URL.Query().Get("id")
	//oid, err := primitive.ObjectIDFromHex(oidStr)
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
}
