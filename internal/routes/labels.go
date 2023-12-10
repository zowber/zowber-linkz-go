package routes

import (
	"html/template"
	"net/http"

	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
)

func labelsHandler(appProps linkzapp.AppProps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type PageProps struct {
			Settings linkzapp.Settings
			Labels   []linkzapp.Label
		}

		labels, err := db.AllLabels()
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err)
		}

		pageProps := PageProps{
			Settings: appProps.Settings,
			Labels:   labels,
		}

		switch r.Method {
		case "GET":
			tmpl := template.Must(template.New("labels").ParseFiles("./templates/head.html", "./templates/header.html", "./templates/labels.tmpl.html", "./templates/footer.html"))
			tmpl.ExecuteTemplate(w, "labels", pageProps)
		default:
			errorHandler(w, r, http.StatusNotImplemented, err)
		}

	}
}
