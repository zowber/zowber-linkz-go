package sqlite

import (
	"log"

	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
)

func (d *SQLiteClient) AllLabels() ([]linkzapp.Label, error) {
    
    db := d.client

    var labels []linkzapp.Label

    rows, err := db.Query(`
        SELECT *
        FROM labels
        ORDER BY name COLLATE NOCASE ASC 
    `)
    if err != nil {
        log.Println("Err getting labels from db")
    }
    var label linkzapp.Label
    for rows.Next() {
        var (
            id int
            name string
        )
        if err := rows.Scan(&id, &name); err != nil {
            log.Println("Err scanning label row")
        }
        label = linkzapp.Label{Id: &id, Name: name}
        labels = append(labels, label)
    }

    return labels, err
}
