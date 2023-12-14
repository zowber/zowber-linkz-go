package routes

import (
	"html/template"
	"net/http"

	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
)

func labelLinksHandler(appProps linkzapp.AppProps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type PageProps struct {
			Settings linkzapp.Settings
			Links    []linkzapp.Link
		}

		var links []linkzapp.Link
		pageProps := PageProps{
			Settings: appProps.Settings,
			Links:    links,
		}

		switch r.Method {
		case "GET":
			tmpl := template.Must(template.New("links-list").ParseFiles("./templates/links-list.html"))
			tmpl.ExecuteTemplate(w, "links-list", pageProps.Links)
		default:
			errorHandler(w, r, http.StatusNotImplemented, err)
		}
	}
}
