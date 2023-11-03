package sqlite

import (
	"database/sql"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
)

// CREATE TABLE links (
//     id INTEGER PRIMARY KEY AUTOINCREMENT,
//     name TEXT,
//     url TEXT,
//     created_at INTEGER
// );

// CREATE TABLE labels (
//     id INTEGER PRIMARY KEY,
//     name TEXT
// );

// CREATE TABLE link_labels (
//     link_id INTEGER,
//     label_id INTEGER,
//     FOREIGN KEY (link_id) REFERENCES links (id),
//     FOREIGN KEY (label_id) REFERENCES labels (id),
//     PRIMARY KEY (link_id, label_id)
// );

// get
// SELECT links.id, links.name AS link_name, links.url, links.created_at, GROUP_CONCAT(labels.name) AS label_names
// FROM links
// LEFT JOIN link_labels ON links.id = link_labels.link_id
// LEFT JOIN labels ON link_labels.label_id = labels.id
// WHERE links.id = :link_id
// GROUP BY links.id;

// insert one
// INSERT INTO links (name, url, created_at) VALUES ('LinkName', 'LinkURL', 'CreatedAtTimestamp');
// SELECT last_insert_rowid();
// INSERT INTO link_labels (link_id, label_id) VALUES (X, A);
// INSERT INTO link_labels (link_id, label_id) VALUES (X, B);
// INSERT INTO link_labels (link_id, label_id) VALUES (X, C);

type SQLiteClient struct {
	client *sql.DB
}

func NewDbClient() (*SQLiteClient, error) {
	db, err := sql.Open("sqlite3", "links.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	return &SQLiteClient{db}, nil
}

func (d *SQLiteClient) All() ([]*linkzapp.Link, error) {
	db := d.client

	rows, err := db.Query("SELECT Id, Name, Url, Labels, CreatedAt FROM links")
	if err != nil {
		log.Println(err)
	}

	var links []*linkzapp.Link
	for rows.Next() {
		var Id, CreatedAt int
		var Name, Url string
		var Labels []linkzapp.Label
		if err := rows.Scan(&Id, &Name, &Url, &Labels, &CreatedAt); err != nil {
			log.Println(err)
		}
		link := &linkzapp.Link{Id: &Id, Name: Name, Url: Url, Labels: Labels, CreatedAt: CreatedAt}
		links = append(links, link)
	}

	return links, nil
}

func (d *SQLiteClient) One(id int) (*linkzapp.Link, error) {
	db := d.client

	row := db.QueryRow("SELECT * from links WHERE id = " + strconv.Itoa(id) + " LIMIT 1")

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

	res, err := db.Exec("INSERT INTO links (Name, Url, Labels, CreatedAt) VALUES ( ?, ?, ?, ? )",
		link.Name, link.Url, link.Labels, link.CreatedAt)
	if err != nil {
		log.Println(err)
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
