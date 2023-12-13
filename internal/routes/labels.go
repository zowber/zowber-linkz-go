package routes

import (
	"html/template"
	"net/http"

	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
)

func labelsHandler(appProps linkzapp.AppProps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type LabelWithLinkCount struct {
			Id        int
			Name      string
			LinkCount int
		}

		type PageProps struct {
			Settings            linkzapp.Settings
			LabelsWithLinkCount []LabelWithLinkCount
		}

		labels, err := db.AllLabels()
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err)
		}

		var labelsWithLinkCount []LabelWithLinkCount
		for _, label := range labels {
			linkCount, _ := db.TotalLinksCountForLabel(*label.Id)
            if label.Name == "" {
                label.Name = "Unlabeled"
            }
			labelsWithLinkCount = append(labelsWithLinkCount, LabelWithLinkCount{Id: *label.Id, Name: label.Name, LinkCount: linkCount})
		}

		pageProps := PageProps{
			Settings:            appProps.Settings,
			LabelsWithLinkCount: labelsWithLinkCount,
		}

		switch r.Method {
		case "GET":
			tmpl := template.Must(template.New("labels").ParseFiles("./templates/head.html", "./templates/header.html", "./templates/labels.tmpl.html", "./templates/footer.html"))
			tmpl.ExecuteTemplate(w, "labels", pageProps)
		default:
			errorHandler(w, r, http.StatusNotImplemented, err)
		}

	}
}
