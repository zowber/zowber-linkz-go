package sqlite

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteClient struct {
	client *sql.DB
}

func NewDbClient() (*SQLiteClient, error) {
	db, err := sql.Open("sqlite3", "newdb.sqlite")
	if err != nil {
		log.Println("Error opening db", err)
	}

	// TODO: Check if tables exist, if not create them

	return &SQLiteClient{db}, err
}

func (d *SQLiteClient) CreateTables() error {
    db := d.client
	
    stmt := `
        CREATE TABLE "links" (
    	"id"	INTEGER NOT NULL,
    	"name"	TEXT NOT NULL,
    	"url"	TEXT NOT NULL,
    	"createdat"	INTEGER NOT NULL,
    	PRIMARY KEY("id" AUTOINCREMENT)
        );

        CREATE TABLE "labels" (
    	"id"	INTEGER NOT NULL,
    	"name"	TEXT NOT NULL,
    	PRIMARY KEY("id" AUTOINCREMENT)
        );

        CREATE TABLE "link_labels" (
	    "link_id"	INTEGER NOT NULL,
	    "label_id"	INTEGER NOT NULL,
	    FOREIGN KEY("link_id") REFERENCES "links"("id"),
	    FOREIGN KEY("label_id") REFERENCES "labels"("id"),
	    PRIMARY KEY("link_id","label_id")
        );

        CREATE TABLE users (
    	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    	name TEXT NOT NULL
        );

        CREATE TABLE settings (
	    user_id INTEGER NOT NULL,
	    color_scheme TEXT,
	    CONSTRAINT settings_FK FOREIGN KEY (user_id) REFERENCES users(id)
        );
    `
    _, err := db.Exec(stmt)
	if err != nil {
		log.Println("Err creatng db tables", err)
	}

	return err
}
