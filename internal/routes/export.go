package routes

import (
	"encoding/csv"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
)

func exportHandler(appProps linkzapp.AppProps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        type PageProps struct {
            Settings linkzapp.Settings 
        }

        pageProps := PageProps{
            Settings: appProps.Settings,
        }

		switch r.Method {
		case "GET":
			tmpl := template.Must(template.New("export.html").ParseFiles("./templates/head.html", "./templates/header.html", "./templates/export.html", "./templates/footer.html"))
			tmpl.ExecuteTemplate(w, "export", pageProps)
		case "POST":
			links, err := db.All()
			if err != nil {
				log.Println("Err getting all links", err)
				errorHandler(w, r, http.StatusInternalServerError, err)
			}

			action := r.PostFormValue("action")
			if action == "csv" {
				w.Header().Set("Content-Disposition", "attachment; filename=export.csv")
				w.Header().Set("Content-Type", "text/csv")
				writer := csv.NewWriter(w)

				for _, link := range links {
					var labels []string
					for _, label := range link.Labels {
						labels = append(labels, label.Name)
					}
					record := []string{strconv.Itoa(*link.Id), link.Name, link.Url, strconv.Itoa(link.CreatedAt)}
					record = append(record, labels...)
					writer.Write(record)
				}
				writer.Flush()
			}

			if action == "json" {
				w.Header().Set("Content-Disposition", "attachment; filename=export.json")
				w.Header().Set("Content-Type", "application/json")
				writer := json.NewEncoder(w)
				writer.Encode(links)
			}

		default:
			errorHandler(w, r, http.StatusMethodNotAllowed, err)
			return
		}
	}
}
