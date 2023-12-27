package sqlite

import (
	"log"
	"sort"
	"time"

	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
)

func (d *SQLiteClient) GetPopularLabelsByYearAndMonth() (linkzapp.PopularLabelsByYearAndMonth, error) {
	var res linkzapp.PopularLabelsByYearAndMonth

	// get all the links
	links, err := d.All()
	if err != nil {
		log.Println("Err", err)
	}

	labelCounts := make(map[string]map[string]int)

	for _, link := range links {
		month := time.Unix(int64(link.CreatedAt), 0).Format("2006-01")
		if labelCounts[month] == nil {
			labelCounts[month] = make(map[string]int)
		}
		for _, label := range link.Labels {
			labelCounts[month][label.Name]++
		}
	}

	// Extract top 3 labels for each month
	topLabels := make(map[string][]string)
	for month, counts := range labelCounts {
		labels := make([]string, 0, len(counts))
		for label := range counts {
			if label != "" {
				labels = append(labels, label)
			}
		}

		sort.Slice(labels, func(i, j int) bool {
			return counts[labels[i]] > counts[labels[j]]
		})

		if len(labels) > 3 {
			labels = labels[:2]
		}

		topLabels[month] = labels
	}

	log.Println(topLabels)

	return res, nil
}
