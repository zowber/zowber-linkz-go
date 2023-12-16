package routes

import (
	"encoding/csv"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
)

// TODO: Sort on import so ascending dated links have ascending ids

func importHandler(appProps linkzapp.AppProps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        type PageProps struct {
            Settings linkzapp.Settings
        }
        pageProps := PageProps{
            Settings: appProps.Settings,
        }

		switch r.Method {
		case "GET":
			tmpl := template.Must(template.New("import.html").ParseFiles("./templates/head.html", "./templates/header.html", "./templates/import.html", "./templates/footer.html"))
			tmpl.ExecuteTemplate(w, "import", pageProps)
		case "POST":
			file, fileHeader, err := r.FormFile("file")
			if err != nil {
				log.Println("Err reading file", err)
				errorHandler(w, r, http.StatusInternalServerError, err)
			}

			fileMime := fileHeader.Header["Content-Type"][0]
			var links []*linkzapp.Link

			file.Seek(0, 0)

			if fileMime == "text/csv" {
				reader := csv.NewReader(file)
				records, err := reader.ReadAll()
				if err != nil {
					log.Println("Err reading csv file", err)
					errorHandler(w, r, http.StatusInternalServerError, err)
				}

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
		    }
            
            // care, sorts in place
            sort.Slice(links, func(i, j int) bool { return links[i].CreatedAt < links[j].CreatedAt })

			if fileMime == "application/json" {
				decoder := json.NewDecoder(file)
				err := decoder.Decode(&links)
				if err != nil {
					log.Println("Err decoding json", err)
				}

				// TODO: template breaks with nil ids
				for i, link := range links {
					counter := i
					link.Id = &counter
				}
			}

			action := r.MultipartForm.Value["action"][0]

			if action == "preview" {

				for _, link := range links {
					log.Println(link)
				}

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
}
