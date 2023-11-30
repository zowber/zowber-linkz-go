package routes

import (
	"html/template"
	"net/http"

	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
)

func settingsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type PageProps struct {
			Settings linkzapp.Settings
		}

		// get settings from DB
		settings, _ := db.GetSettings()

		props := PageProps{Settings: *settings}

		if r.Method == "GET" {
			tmpl := template.Must(template.New("settings").Funcs(funcMap).ParseFiles("./templates/settings.tmpl.html", "./templates/head.html", "./templates/header.html", "./templates/footer.html"))
			tmpl.ExecuteTemplate(w, "settings", props)
		}

        if r.Method == "POST" {

           //var (
           //    name = r.PostFormValue("name")
           //    colourScheme = r.PostFormValue("colour-scheme") 
           //)
           //
           //newUser := linkzapp.User{
           //    Name: name,
           //}

           //newSettings := linkzapp.Settings{
           //    ColourScheme: colourScheme,
           //}

        }
	}
}
