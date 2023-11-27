package routes

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
)

type PageData struct {
	Links      []*linkzapp.Link
	Page       int
	PerPage    int
	TotalLinks int
	HasPrev    bool
	PrevPage   int
	HasNext    bool
	NextPage   int
}

var linksHandler = func(w http.ResponseWriter, r *http.Request) {

	totalLinks, _ := db.TotalLinksCount()

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(r.URL.Query().Get("perpage"))
	if err != nil || perPage < 25 {
		perPage = 25
	}

	hasPrev := page > 1
	prevPage := page - 1

	hasNext := (page * perPage) < totalLinks
	nextPage := page + 1

	offset := (page - 1) * perPage

	links, err := db.All(perPage, offset)
	if err != nil {
		log.Print(err.Error())
	}

	pageData := PageData{
		Links:      links,
		Page:       page,
		PerPage:    perPage,
		TotalLinks: totalLinks,
		HasPrev:    hasPrev,
		PrevPage:   prevPage,
		HasNext:    hasNext,
		NextPage:   nextPage,
	}

	accepts := make(map[string]bool)
	for _, el := range strings.Split(r.Header["Accept"][0], ",") {
		accepts[el] = true
	}

	if accepts["text/html"] || accepts["*/*"] {
		switch r.Method {
		case "GET":
			tmpl := template.Must(template.New("links.html").Funcs(funcMap).ParseFiles("./templates/header.html", "./templates/links.tmpl.html", "./templates/links-list.html", "./templates/link.html", "./templates/footer.html"))
			tmpl.ExecuteTemplate(w, "links", pageData)
		default:
			errorHandler(w, r, http.StatusMethodNotAllowed, err)
		}
	}

	if accepts["application/json"] {
		switch r.Method {
		case "GET":
			jsonData, err := json.Marshal(pageData.Links)
			if err != nil {
				log.Println(err)
			}

			w.Header().Set("Content-Type:", "application/json")
			w.Write(jsonData)

		default:
			errorHandler(w, r, http.StatusMethodNotAllowed, err)
		}
	}
}
