package routes

import (
	"log"
	"net/http"
	"text/template"

	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
)

func setupHandler(appProps AppProps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("In the settings handler...")

		type PageProps struct {
			Settings linkzapp.Settings
		}

		pageProps := PageProps{
			Settings: appProps.Settings,
		}

		if r.Method == "GET" {
			tmpl := template.Must(template.ParseFiles("./templates/head.html", "./templates/header.html", "./templates/setup.tmpl.html", "./templates/footer.html"))
			tmpl.ExecuteTemplate(w, "setup", pageProps)
		}

        if r.Method == "POST" {
            var (
                name = r.PostFormValue("name")
                theme = r.PostFormValue("theme")
            )

            newUser := linkzapp.User{
                Name: name,
            }

            user, err := db.InsertUser(&newUser) 
            if err != nil {
                log.Println("Err creating new user")
                errorHandler(w, r, http.StatusInternalServerError, err)
            }

            newSettings := linkzapp.Settings{
                ColorScheme: theme,
            }
            
            if err = db.InsertSettings(user, &newSettings); err != nil {
                log.Println("Err creating settings for new user", err)
                errorHandler(w, r, http.StatusInternalServerError, err)
            }

            http.Redirect(w, r, "/links", http.StatusSeeOther)
        }
	}
}
