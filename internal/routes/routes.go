package routes

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/zowber/zowber-linkz-go/internal/data/sqlite"
	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
)

// TODO: Sort on import so ascending dated links have ascending ids
// TODO: Figure out the one true CSV format for import/export
// TODO: Implement /label/:id/links
// TODO: Stats/analytics?
// TODO: First run/setup i.e., create tables, store some kind of config in the db, etc.

type AppProps struct {
	Settings linkzapp.Settings
}

var db, err = sqlite.NewDbClient()

func idToStr(id int) string {
	idStr := strconv.Itoa(id)
	return idStr
}

func idStrToId(idStr string) int {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// errorHandler(w, r, http.StatusBadRequest, err)
		return -1
	}
	return id
}

// TODO: Refactor Link struct to use type time.Time
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
	settings, _ := db.GetSettings()
	appProps := AppProps{Settings: *settings}

	staticDir := "./static"
	staticServer := http.FileServer(http.Dir(staticDir))

	mux := http.NewServeMux()
	mux.Handle("/", indexHandler())
	mux.Handle("/links", linksHandler())
	mux.Handle("/link/edit", editHandler())
	mux.Handle("/link", linkHandler())
	mux.Handle("/labels", labelsHandler())
	mux.Handle("/label", labelHandler())
	// mux.HandleFunc("/label/:id/links")
	mux.Handle("/static", http.StripPrefix("/", staticServer))
	mux.Handle("/settings", settingsHandler())
	mux.Handle("/import", importHandler(appProps))
	mux.Handle("/export", exportHandler())

	return mux
}
