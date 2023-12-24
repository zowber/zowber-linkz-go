package sqlite

import (
	"log"
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

    latestUnixTime := links[0].CreatedAt
    oldestUnixTime := links[len(links)-1].CreatedAt

    latestLinkTime := time.Unix(int64(latestUnixTime), 0)
    oldestLinkTime := time.Unix(int64(oldestUnixTime), 0)

    latestYear := latestLinkTime.Year()
    oldestYear := oldestLinkTime.Year()

    currYear := oldestYear
    for i := oldestYear; i <= latestYear ; i++ {
        currYear++
    }

	return res, nil
}
