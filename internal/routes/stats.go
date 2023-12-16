package routes

import (
	"html/template"
	"net/http"

	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
)

func statsHandler(appProps linkzapp.AppProps) http.HandlerFunc {
    return func(w http.ResponseWriter, r * http.Request) {

        tmpl := template.Must(template.New("stats").ParseFiles("./templates/stats.tmpl.html"))

        tmpl.ExecuteTemplate(w, "stats", nil)
        
    }
}
