package routes

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/zowber/zowber-linkz-go/internal/data/sqlite"
)

// TODO: Sort on import so ascending dated links have ascending ids
// TODO: Figure out the one true CSV format for import/export
// TODO: Pagiation
// TODO: Implement /label/:id/links
// TODO: Stats/analytics?
// TODO: First run/setup i.e., create tables, store some kind of config in the db, etc.

var db, err = sqlite.NewDbClient()

func idToStr(id int) string {
	idStr := strconv.Itoa(id)
	return idStr
}

func formatDate(unixTime int) string {
	timeVal := time.Unix(int64(unixTime), 0)
	formattedTime := timeVal.Format("02 Jan 2006")
	return formattedTime
}

var funcMap = template.FuncMap{
	"idToStr":    idToStr,
	"formatDate": formatDate,
}

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/links", linksHandler)
	mux.HandleFunc("/link/edit", editHandler)
	mux.HandleFunc("/link", linkHandler)
	mux.HandleFunc("/labels", labelsHandler)
	mux.HandleFunc("/label", labelHandler)
    // mux.HandleFunc("/label/:id/links")

	mux.HandleFunc("/scripts/links.js", staticHandler)

	mux.HandleFunc("/import", importHandler)
	mux.HandleFunc("/export", exportHandler)

	return mux
}
