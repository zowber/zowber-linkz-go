package sqlite

import (
	"log"

	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
)

func (d *SQLiteClient) CountUsers() (int, error) {
	db := d.client

	var count int
	err := db.QueryRow(`
        SELECT COUNT(*) from users;
    `).Scan(&count)
	if err != nil {
		log.Println("Err counting users", err)
	}

	return count, err
}

func (d *SQLiteClient) GetUser() (*linkzapp.User, error) {
    db := d.client

    var user linkzapp.User

    err := db.QueryRow(`
        SELECT *
        FROM users
        WHERE id = ?
    `, 1).Scan(&user.Id, &user.Name)
    if err != nil {
        log.Println("Err getting user", err)
    }

    return &user, err

}

func (d *SQLiteClient) InsertUser(user *linkzapp.User) (*linkzapp.User, error) {
	db := d.client

	newUser := linkzapp.User{}

	err := db.QueryRow(`
        INSERT INTO users (name) 
        VALUES (?)
        RETURNING *
    `, user.Name).Scan(&newUser.Id, &newUser.Name)
	if err != nil {
		log.Println("Err inserting new user", err)
	}

	return &newUser, err
}
