package routes

import (
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/zowber/zowber-linkz-go/internal/data"
	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var db, err = data.NewPortfolioDbClient()

func oidToStr(oid primitive.ObjectID) string {
	return oid.Hex()
}

var funcMap = template.FuncMap{
	"oidToStr": oidToStr,
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

			// this seems a bit nasty
			formRaw := r.Form
			var labels []linkzapp.Label
			for key, value := range formRaw {
				if strings.Contains(key, "label_") {
					labels = append(labels, linkzapp.Label{Id: key, Name: value[0]})
				}
			}

			newLink, err := db.Insert(&linkzapp.Link{Name: name, Url: url, Labels: labels, CreatedAt: time.Now().Unix()})
			if err != nil {
				errorHandler(w, r, http.StatusInternalServerError, err)
				return
			}

			tmpl := template.Must(template.New("link.html").Funcs(funcMap).ParseFiles("./templates/link.html"))
			tmpl.ExecuteTemplate(w, "link", newLink)
		}
	}
}

var linkHandler = func(w http.ResponseWriter, r *http.Request) {
	// oidStr := r.URL.Path[len("/link/"):]
	oidStr := r.URL.Query().Get("id")
	oid, err := primitive.ObjectIDFromHex(oidStr)
	if err != nil {
		errorHandler(w, r, http.StatusBadRequest, err)
		return
	}

	link, err := db.One(oid)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err)
		return
	}

	tmpl := template.Must(template.New("link.html").Funcs(funcMap).ParseFiles("./templates/link.html"))
	tmpl.ExecuteTemplate(w, "link", link)
}

var editHandler = func(w http.ResponseWriter, r *http.Request) {
	oidStr := r.URL.Query().Get("id")
	oid, err := primitive.ObjectIDFromHex(oidStr)
	if err != nil {
		errorHandler(w, r, http.StatusBadRequest, err)
		return
	}

	switch r.Method {
	case "GET":
		linkToEdit, err := db.One(oid)
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err)
			return
		}

		tmpl := template.Must(template.New("edit.html").Funcs(funcMap).ParseFiles("./templates/edit.html", "./templates/label.html"))
		tmpl.Execute(w, linkToEdit)
	case "PUT":
		name := r.PostFormValue("name")
		url := r.PostFormValue("url")

		// this seems a bit nasty
		formRaw := r.Form
		var labels []linkzapp.Label
		for key, value := range formRaw {
			if strings.Contains(key, "label_") {
				labels = append(labels, linkzapp.Label{Id: key, Name: value[0]})
			}
		}

		link := &linkzapp.Link{
			Id:     oid,
			Name:   name,
			Url:    url,
			Labels: labels,
		}
		updatedLink, err := db.Update(link)
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

		data := map[string]string{"Name": name, "Id": id}
		tmpl := template.Must(template.New("label.html").ParseFiles("./templates/label.html"))
		tmpl.ExecuteTemplate(w, "label", data)
	}
}

var deleteHandler = func(w http.ResponseWriter, r *http.Request) {
	oidStr := r.URL.Query().Get("id")
	oid, err := primitive.ObjectIDFromHex(oidStr)
	if err != nil {
		errorHandler(w, r, http.StatusBadRequest, err)
		return
	}

	err = db.Delete(oid)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err)
		return
	}
}
