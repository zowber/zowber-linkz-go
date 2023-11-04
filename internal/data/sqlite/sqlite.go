package sqlite

import (
	"database/sql"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
)

// CREATE TABLE "links" (
// 	"id"	INTEGER NOT NULL,
// 	"name"	TEXT NOT NULL,
// 	"url"	TEXT NOT NULL,
// 	"createdat"	INTEGER NOT NULL,
// 	PRIMARY KEY("id" AUTOINCREMENT)
// );

// CREATE TABLE "labels" (
// 	"id"	INTEGER NOT NULL,
// 	"name"	TEXT NOT NULL,
// 	PRIMARY KEY("id" AUTOINCREMENT)
// );

// CREATE TABLE "link_labels" (
// 	"link_id"	INTEGER NOT NULL,
// 	"label_id"	INTEGER NOT NULL,
// 	FOREIGN KEY("link_id") REFERENCES "links"("id"),
// 	FOREIGN KEY("label_id") REFERENCES "labels"("id"),
// 	PRIMARY KEY("link_id","label_id")
// );

type SQLiteClient struct {
	client *sql.DB
}

func NewDbClient() (*SQLiteClient, error) {
	db, err := sql.Open("sqlite3", "links2.sqlite")
	if err != nil {
		log.Println(err)
	}
	return &SQLiteClient{db}, nil
}

func (d *SQLiteClient) All() ([]*linkzapp.Link, error) {
	db := d.client

	rows, err := db.Query(`
		SELECT * FROM links
	`)
	if err != nil {
		log.Println(err)
	}

	var links []*linkzapp.Link
	for rows.Next() {
		var id, createdAt int
		var name, url string
		if err := rows.Scan(&id, &name, &url, &createdAt); err != nil {
			log.Println(err)
		}

		//get the labels
		labels := []linkzapp.Label{}
		labelRows, err := db.Query(`
			SELECT labels.id, labels.name AS name
			FROM labels
			INNER JOIN link_labels ON labels.id = link_labels.label_id
			WHERE link_labels.link_id = ?;
		`, id)
		if err != nil {
			log.Println(err)
		}
		for labelRows.Next() {
			var labelId int
			var labelName string
			if err := labelRows.Scan(&labelId, &labelName); err != nil {
				log.Println("here", err)
			}
			label := &linkzapp.Label{Id: labelId, Name: labelName}
			labels = append(labels, *label)
		}
		link := &linkzapp.Link{Id: &id, Name: name, Url: url, Labels: labels, CreatedAt: createdAt}
		links = append(links, link)
	}

	return links, nil
}

func (d *SQLiteClient) One(id int) (*linkzapp.Link, error) {
	db := d.client

	row := db.QueryRow("SELECT * from links WHERE id = " + strconv.Itoa(id) + " LIMIT 1;")

	var Id, CreatedAt int
	var Name, Url string
	var Labels []linkzapp.Label
	if err := row.Scan(&Id, &Name, &Url, &Labels, &CreatedAt); err != nil {
		log.Println(err)
	}
	link := &linkzapp.Link{Id: &Id, Name: Name, Url: Url, Labels: Labels, CreatedAt: CreatedAt}

	return link, nil
}

func (d *SQLiteClient) Insert(link *linkzapp.Link) (*linkzapp.Link, error) {
	db := d.client

    tx, err := db.Begin()
    if err != nil {
        log.Println(err)
    }

	// insert the link
    var LinkId int 
    err = tx.QueryRow(`
        INSERT INTO links (name, url, createdat)
        VALUES ( ?, ?, ? );
        SELECT last_insert_rowid();
        `, link.Name, link.Url, link.CreatedAt).Scan(&LinkId)
	if err != nil {
        log.Println("Err insering link:", err)
	}

    // insert the label(s)
	labelIds := make([]int, len(link.Labels))
	for i, label := range link.Labels {
		err = tx.QueryRow(`
            INSERT INTO labels (name)
           	VALUES (?)
			SELECT last_insert_rowid();
		`, label.Name).Scan(&labelIds[i])
		if err != nil {
            log.Println("Err inserting link labels:", err)
		}
	}

	// insert the relations
    for _, labelId := range labelIds {
        _, err = tx.Exec(`
            INSERT INTO links_labels (link_id, label_id)
            VALUES (?, ?)
        `, link.Id, labelId)
        if err != nil {
            log.Println("Err inserting link-label relation:", err)}
    }


	newId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	newIdInt := int(newId)

	insertedLink := &linkzapp.Link{
		Id:        &newIdInt,
		Name:      link.Name,
		Url:       link.Url,
		Labels:    link.Labels,
		CreatedAt: link.CreatedAt,
	}

	return insertedLink, err
}

func (d *SQLiteClient) Update(id int, link *linkzapp.Link) (*linkzapp.Link, error) {
	db := d.client

	_, err := db.Exec("UPDATE links SET Name = ?, Url = ?, Labels = ?, CreatedAt = ? WHERE Id = ?",
		link.Name, link.Url, link.Labels, link.CreatedAt, link.Id)
	if err != nil {
		log.Println(err)
	}

	return nil, nil
}

func (d *SQLiteClient) Delete(id int) error {
	db := d.client

	_, err := db.Exec("DELETE from links WHERE id = " + strconv.Itoa(id))
	if err != nil {
		log.Println(err)
	}

	return err
}
