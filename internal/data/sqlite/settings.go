package sqlite

import (
	"log"
	"github.com/zowber/zowber-linkz-go/pkg/linkzapp"
)

func (d *SQLiteClient) GetSettings() (*linkzapp.Settings, error) {
	db := d.client

	var userId int
	var colourScheme string
	err := db.QueryRow(`
        SELECT * FROM settings;
    `).Scan(&userId, &colourScheme)

	settings := linkzapp.Settings{ColorScheme: colourScheme}

	return &settings, err
}

func (d *SQLiteClient) InsertSettings(user *linkzapp.User, settings *linkzapp.Settings) error {
	db := d.client

	_, err := db.Exec(`
        INSERT INTO settings (user_id, color_scheme)
        VALUES (?, ?)
    `, user.Id, settings.ColorScheme)
	if err != nil {
		log.Println("Err inserting settings for new user", err)
	}

	return err
}

func (d *SQLiteClient) UpdateSettings(user *linkzapp.User, settings *linkzapp.Settings) error {
	db := d.client

	_, err := db.Exec(`
        UPDATE settings
        SET colour_scheme = ?
        WHERE user_id = ?;
        `, settings.ColorScheme, user.Name)

	return err
}
