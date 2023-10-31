package sqlite

import (
	"database/sql"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
)

type SQLiteClient struct {
	db *sql.DB
}

func NewDbClient() (*SQLiteClient, error) {
	db, err := sql.Open("sqlite3", "links.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	return &SQLiteClient{db}, nil
}

func (d *SQLiteClient) All() ([]*linkzapp.Link, error) {
	db, err := sql.Open("sqlite3", "links.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT Id, Name, Url, Labels, CreatedAt FROM links")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var links []*linkzapp.Link
	for rows.Next() {
		var Id, CreatedAt int
		var Name, Url, Labels string
		if err := rows.Scan(&Id, &Name, &Url, &Labels, &CreatedAt); err != nil {
			log.Fatal(err)
		}
		link := &linkzapp.Link{&Id, Name, Url, Labels, CreatedAt}
		links = append(links, link)
	}

	return links, nil
}

// One
func (d *SQLiteClient) One(id int) (*linkzapp.Link, error) {
	db, err := sql.Open("sqlite3", "links.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("DELETE from links WHERE id = " + strconv.Itoa(id))
	if err != nil {
		log.Fatal(err)
	}

	return nil, nil
}

// Insert
func (d *SQLiteClient) Insert(link *linkzapp.Link) (*linkzapp.Link, error) {
	db, err := sql.Open("sqlite3", "links.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	res, err := db.Exec("INSERT INTO links (Name, Url, Labels, CreatedAt) VALUES ( ?, ?, ?, ? )",
		link.Name, link.Url, link.Labels, link.CreatedAt)
	if err != nil {
		log.Fatal(err)
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

// Update
func (d *SQLiteClient) Update(id int, link *linkzapp.Link) (*linkzapp.Link, error) {
	db, err := sql.Open("sqlite3", "links.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("UPDATE links SET Name = ?, Url = ?, Labels = ?, CreatedAt = ? WHERE Id = ?",
		link.Name, link.Url, link.Labels, link.CreatedAt, link.Id)
	if err != nil {
		log.Fatal(err)
	}

	return nil, nil
}

// Delete
func (d *SQLiteClient) Delete(id int) error {
	db, err := sql.Open("sqlite3", "links.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("DELETE from links WHERE id = " + strconv.Itoa(id))
	if err != nil {
		log.Fatal(err)
	}

	return err
}
