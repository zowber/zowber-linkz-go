package routes

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
)

func statsHandler(appProps linkzapp.AppProps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type TotalLinksByMonth struct {
			Month int
			Total int
		}

		type TotalLinksByYear struct {
			Year   int
			Total  int
			Months []TotalLinksByMonth
		}

		type TotalLinksByYearAndMonth struct {
			Years []TotalLinksByYear
		}

		type PageProps struct {
			Settings                 linkzapp.Settings
			TotalLinks               int
			TotalLinksByYearAndMonth TotalLinksByYearAndMonth
		}

		// get settings from DB
		settings, err := db.GetSettings()
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err)
		}

		totalLinks, err := db.TotalLinksCount()
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err)
		}

		// links by month/year
		// 2023 (145 links added)
		// Jan (5), Feb (25)...
		// 2022...
		totalLinksByYearAndMonth := TotalLinksByYearAndMonth{
			Years: []TotalLinksByYear{
				{
					Year:  2023,
					Total: 500,
					Months: []TotalLinksByMonth{
						{
							Month: 1,
							Total: 100,
						},
					},
				},
			},
		}
		fmt.Println(totalLinksByYearAndMonth)

		// popular tags by month/year
		// 2023 (ai, golang, dubai)
		// Jan (ai), Feb (dubai)...
		// 2022...

		props := PageProps{
			Settings:                 *settings,
			TotalLinks:               totalLinks,
			TotalLinksByYearAndMonth: totalLinksByYearAndMonth,
		}

		tmpl := template.Must(template.New("stats").ParseFiles("./templates/head.html", "./templates/header.html", "./templates/stats.tmpl.html", "./templates/footer.html"))

		tmpl.ExecuteTemplate(w, "stats", props)

	}
}
