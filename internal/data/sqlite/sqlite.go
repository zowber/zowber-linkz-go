package sqlite

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
)

type SQLiteClient struct {
	client *sql.DB
}

func NewDbClient() (*SQLiteClient, error) {
    db, err := sql.Open("sqlite3", "linkz.sqlite")
	if err != nil {
		log.Println("Error opening db", err)
	}

	return &SQLiteClient{db}, err
}

func CreateTables() error {
	db, err := NewDbClient()
    if err != nil {
        log.Println("Err getting new db client", err)
    }
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
    `
    _, err = db.client.Exec(stmt)
    if err != nil {
        log.Println("Err creatng db tables", err)
    }

    return err
}

func (d *SQLiteClient) All() ([]*linkzapp.Link, error) {
	db := d.client

	log.Println("Get All")
	rows, err := db.Query(`
		SELECT * FROM links
        ORDER BY createdat DESC;
	`)
	if err != nil {
		log.Println(err)
	}

	var links []*linkzapp.Link
	for rows.Next() {
		var id, createdat int
		var name, url string
		if err := rows.Scan(&id, &name, &url, &createdat); err != nil {
			log.Println("Err scanning link row", err)
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
				log.Println("Err scanning label row", err)
			}
			label := &linkzapp.Label{Id: &labelId, Name: labelName}
			labels = append(labels, *label)
		}
		link := &linkzapp.Link{Id: &id, Name: name, Url: url, Labels: labels, CreatedAt: createdat}
		links = append(links, link)
	}

	return links, err
}

func (d *SQLiteClient) One(id int) (*linkzapp.Link, error) {
	db := d.client
	// TODO: is it necessary do do defer db.Close() ??

	log.Println("Get One with id", id)
	var Id, Createdat int
	var Name, Url string

	err := db.QueryRow(`
        SELECT * from links
        WHERE id = ?
        LIMIT 1;
    `, id).Scan(&Id, &Name, &Url, &Createdat)
	if err != nil {
		log.Printf("Err getting link %d: %s", id, err)
		return nil, err
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
		label := &linkzapp.Label{Id: &labelId, Name: labelName}
		Labels = append(Labels, *label)
	}
	link := &linkzapp.Link{Id: &Id, Name: Name, Url: Url, Labels: Labels, CreatedAt: Createdat}

	return link, nil
}

func (d *SQLiteClient) Insert(link *linkzapp.Link) (int, error) {
	db := d.client

	// start a new txn
	tx, err := db.Begin()
	if err != nil {
		log.Println("Err begining Insert txn", err)
	}

	// insert link and get it's id
	var linkId int
	err = tx.QueryRow(`
        INSERT INTO links (name, url, createdat)
        VALUES ( ?, ?, ? )
        RETURNING id;
        `, link.Name, link.Url, link.CreatedAt).Scan(&linkId)
	if err != nil {
		log.Println("Err Inserting link:", err)
	}

	// get id of any existing labels
	// insert any new labels and get their id
	labelIds := make([]int, len(link.Labels))
	for i, label := range link.Labels {
		var existsId int
		err = tx.QueryRow(`
            SELECT id FROM labels WHERE name = ?;
            `, label.Name).Scan(&existsId)
		if err != nil {
			log.Println("Err checking if label exists", err)
		}

		if existsId != 0 {
			labelIds[i] = existsId
		} else {
			err = tx.QueryRow(`
                INSERT INTO labels (name)
                VALUES (?)
                RETURNING id;
            `, label.Name).Scan(&labelIds[i])
			if err != nil {
				log.Println("Err Inserting new label:", err)
			}
		}
	}

	// insert the relations
	for _, labelId := range labelIds {
		_, err = tx.Exec(`
            INSERT INTO link_labels (link_id, label_id)
            VALUES (?, ?)
        `, linkId, labelId)
		if err != nil {
			log.Println("Err Inserting link-label association:", err)
		}
	}

	// commit the txn
	if err := tx.Commit(); err != nil {
		log.Println("Err commiting insert:", err)
	}

	return linkId, err
}

func (d *SQLiteClient) Update(id int, link *linkzapp.Link) error {
	db := d.client

	// get prev labels
	prevLabelRows, err := db.Query(`
		SELECT labels.id, labels.name AS name
		FROM labels
		INNER JOIN link_labels ON labels.id = link_labels.label_id
		WHERE link_labels.link_id = ?;
        `, id)
	if err != nil {
		log.Println(err)
	}

	prevLabels := []linkzapp.Label{}

	for prevLabelRows.Next() {
		var labelId int
		var labelName string
		if err := prevLabelRows.Scan(&labelId, &labelName); err != nil {
			log.Println(err)
		}
		label := linkzapp.Label{Id: &labelId, Name: labelName}
		prevLabels = append(prevLabels, label)
	}

	// determine which prevLabels are no longer wanted
	currLabels := link.Labels
	var unwantedLabels []linkzapp.Label

	newMap := make(map[string]bool)
	for _, label := range currLabels {
		newMap[label.Name] = true
	}
	for _, prevLabel := range prevLabels {
		if _, exists := newMap[prevLabel.Name]; !exists {
			unwantedLabels = append(unwantedLabels, prevLabel)
		}
	}

	// get the ids of unwanted labels
	for i, label := range unwantedLabels {
		var uwlId int
		if err := db.QueryRow(`
            SELECT id FROM labels
            WHERE Name = ?
            `, label.Name).Scan(&uwlId); err != nil {
			log.Println(err)
		}
		unwantedLabels[i] = linkzapp.Label{Id: &uwlId, Name: label.Name}
	}

	// remove the associations between the link and unwanted labels
	for _, label := range unwantedLabels {
		_, err := db.Exec(`
	        DELETE FROM link_labels
	        WHERE link_id = ? AND label_id = ?
	    `, id, *label.Id)
		if err != nil {
			log.Println("Err removing associations", err)
		}
	}

	// TODO: remove orphan labels

	// insert new labels
	for _, label := range currLabels {
		_, err := db.Exec(`
            INSERT OR IGNORE INTO labels (name)
            VALUES (?)
        `, label.Name)
		if err != nil {
			log.Println("Err inserting new labels", err)
		}
	}

	// insert associations to new labels (some of which may exist)
	// get ids for the labels
	for i, label := range currLabels {
		var labelId int
		err := db.QueryRow(`
            SELECT id FROM labels
            WHERE (name) = ?
        `, label.Name).Scan(&labelId)
		if err != nil {
			log.Println("Err getting id for currLabels", err)
		}
		currLabels[i] = linkzapp.Label{Id: &labelId, Name: label.Name}
	}

	// insert associations between link and target labels
	for _, label := range currLabels {
		_, err = db.Exec(`
            INSERT OR IGNORE INTO link_labels (link_id, label_id)
            VALUES (?, ?)
        `, id, *label.Id)
		if err != nil {
			log.Println("Err Inserting link-label association:", err)
		}
	}

	// update the link
	_, err = db.Exec(`
        UPDATE links SET Name = ?, Url = ? 
        WHERE Id = ?
    `, link.Name, link.Url, id)
	if err != nil {
		log.Println("Err updating link", err)
	}

	// should we return early in some cases?
	return err
}

func (d *SQLiteClient) Delete(id int) error {
	db := d.client

	// delete associations
	_, err := db.Exec(`
        DELETE from link_labels
        WHERE link_id = ?
    `, id)
	if err != nil {
		log.Println("Err deleting associations", err)
	}

	// delete link
	_, err = db.Exec(`
        DELETE from links
        WHERE id = ?
    `, id)
	if err != nil {
		log.Println("Err deleting link", err)
	}

	return err
}
