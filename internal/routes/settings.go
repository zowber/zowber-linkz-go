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
            User linkzapp.User
		}

		// get settings from DB
		settings, err := db.GetSettings()
        if err != nil {
            errorHandler(w, r, http.StatusInternalServerError, err)
        }
        // TODO: feels like this should get passed in from routes.go or whatever
        user, err := db.GetUser()
        if err != nil {
            errorHandler(w, r, http.StatusInternalServerError, err)
        }

        props := PageProps{Settings: *settings, User: *user }

		if r.Method == "GET" {
			tmpl := template.Must(template.New("settings").Funcs(funcMap).ParseFiles("./templates/settings.tmpl.html", "./templates/head.html", "./templates/header.html", "./templates/footer.html"))
			tmpl.ExecuteTemplate(w, "settings", props)
		}

		if r.Method == "POST" {
            //TODO: implement update settings
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
