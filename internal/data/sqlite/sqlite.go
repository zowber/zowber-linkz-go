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
		var id, createdat int
		var name, url string
		if err := rows.Scan(&id, &name, &url, &createdat); err != nil {
			log.Println(err)
		}

		// get the labels
		labelRows, err := db.Query(`
			SELECT labels.id, labels.name AS name
			FROM labels
			INNER JOIN link_labels ON labels.id = link_labels.label_id
			WHERE link_labels.link_id = ?;
		`, id)
		if err != nil {
			log.Println(err)
		}

		labels := []linkzapp.Label{}
		for labelRows.Next() {
			var labelId int
			var labelName string
			if err := labelRows.Scan(&labelId, &labelName); err != nil {
				log.Println("here", err)
			}
			label := &linkzapp.Label{Id: labelId, Name: labelName}
			labels = append(labels, *label)
		}
		link := &linkzapp.Link{Id: &id, Name: name, Url: url, Labels: labels, CreatedAt: createdat}
		links = append(links, link)
	}

	return links, err
}

func (d *SQLiteClient) One(id int) (*linkzapp.Link, error) {
	db := d.client

	log.Println("getting one with id:", id)
	row := db.QueryRow("SELECT * from links WHERE id = " + strconv.Itoa(id) + " LIMIT 1;")

	var Id, Createdat int
	var Name, Url string
	if err := row.Scan(&Id, &Name, &Url, &Createdat); err != nil {
		log.Println(err)
	}

	// get the label(s)
	var Labels []linkzapp.Label
	labelRows, err := db.Query(`
		SELECT labels.id, labels.name AS name
		FROM labels
		INNER JOIN link_labels ON labels.id = link_labels.label_id
		WHERE link_labels.link_id = ?;
	`, Id)
	if err != nil {
		log.Println(err)
	}

	for labelRows.Next() {
		var labelId int
		var labelName string
		if err := labelRows.Scan(&labelId, &labelName); err != nil {
			log.Println("Err scanning label rows", err)
		}
		label := &linkzapp.Label{Id: labelId, Name: labelName}
		Labels = append(Labels, *label)
	}
	link := &linkzapp.Link{Id: &id, Name: Name, Url: Url, Labels: Labels, CreatedAt: Createdat}

	return link, nil
}

func (d *SQLiteClient) Insert(link *linkzapp.Link) (int, error) {
	db := d.client

	// start a new txn
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
	}

	// insert the link
	var linkId int
	err = tx.QueryRow(`
        INSERT INTO links (name, url, createdat)
        VALUES ( ?, ?, ? )
        RETURNING id;
        `, link.Name, link.Url, link.CreatedAt).Scan(&linkId)
	if err != nil {
		log.Println("Err insering link:", err)
	}

	// insert the label(s)
	labelIds := make([]int, len(link.Labels))
	for i, label := range link.Labels {
		err = tx.QueryRow(`
            INSERT INTO labels (name)
           	VALUES (?)
            RETURNING id;
		`, label.Name).Scan(&labelIds[i])
		if err != nil {
			log.Println("Err inserting link labels:", err)
		}
	}

	// insert the relations
	for _, labelId := range labelIds {
		_, err = tx.Exec(`
            INSERT INTO link_labels (link_id, label_id)
            VALUES (?, ?)
        `, linkId, labelId)
		if err != nil {
			log.Println("Err inserting link-label relation:", err)
		}
	}

	// commit the txn
	if err := tx.Commit(); err != nil {
		log.Println("Err commiting insert:", err)
	}

	return linkId, err
}

func (d *SQLiteClient) Update(id int, link *linkzapp.Link) (*linkzapp.Link, error) {
	db := d.client

    // 1. have labels changed?
    // 2. then delete or add labels
    // 3. has link changed?
    // 4. then update link
    // 5. update associations?
	
    // 1. get current labels
    // get the labels
	labelRows, err := db.Query(`
		SELECT labels.id, labels.name AS name
		FROM labels
		INNER JOIN link_labels ON labels.id = link_labels.label_id
		WHERE link_labels.link_id = ?;
		`, link.Id)
	if err != nil {
		log.Println(err)
	}

	labels := []linkzapp.Label{}
	for labelRows.Next() {
		var labelId int
		var labelName string
		if err := labelRows.Scan(&labelId, &labelName); err != nil {
			log.Println("here", err)
		}
		label := &linkzapp.Label{Id: labelId, Name: labelName}
		labels = append(labels, *label)
	}


    _, err := db.Exec("UPDATE links SET Name = ?, Url = ?, Labels = ?, CreatedAt = ? WHERE Id = ?",
		link.Name, link.Url, link.Labels, link.CreatedAt, link.Id)
	if err != nil {
		log.Println(err)
	}

	return nil, nil
}

func (d *SQLiteClient) Delete(id int) error {
	db := d.client
    // delete associations
    // leave labels if not associated with other links
    // delete link
	_, err := db.Exec("DELETE from links WHERE id = " + strconv.Itoa(id))
	if err != nil {
		log.Println(err)
	}

	return err
}
