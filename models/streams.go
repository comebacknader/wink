package models

import (
	"database/sql"
	_ "time"

	"github.com/comebacknader/wink/config"
)

type Stream struct {
	ID       int
	UserID   int
	Title    string
	Game     string
	Withgame bool
	Online   bool
	Twit     string
	SiteOne  string
}

type Frontpagers struct {
	ID     int
	UserID int
}

// GetStreamByID gets the stream by a User's ID.
func GetStreamByUID(UID int) (Stream, bool) {
	strm := Stream{}
	err := config.DB.
		QueryRow("SELECT ID, user_id, title, game, online, twitter, siteone FROM streams WHERE user_id = $1", UID).
		Scan(&strm.ID, &strm.UserID, &strm.Title, &strm.Game, &strm.Online, &strm.Twit, &strm.SiteOne)
	if err != nil {
		if err == sql.ErrNoRows {
			return strm, false
		}
		panic(err)
	}
	return strm, true
}

// UpdateStream updates the stream by User's ID
func UpdateStream(title string, game string, twit string, siteone string, UID int) {
	_, err := config.DB.
		Exec("UPDATE streams SET title = $1, game = $2, twitter = $3, siteone = $4 WHERE user_id = $5",
			title, game, twit, siteone, UID)
	if err != nil {
		panic(err)
	}
}

// StreamOnline changes the status of a Stream to Online.
func StreamOnline(UID int) {
	_, err := config.DB.
		Exec("UPDATE streams SET online = $1 WHERE user_id = $2",
			true, UID)
	if err != nil {
		panic(err)
	}
}

// StreamOffline changes the status of a Stream to Offline.
func StreamOffline(UID int) {
	_, err := config.DB.
		Exec("UPDATE streams SET online = $1 WHERE user_id = $2",
			false, UID)
	if err != nil {
		panic(err)
	}
}

// GetFrontpagers gets ALL the streamers that are on the frontpage
func GetFrontpagers() ([]User, error) {
	fpagers := []User{}
	// SELECT ALL USERS WHERE id = FRONTPAGE(USER_ID)
	rows, err := config.DB.Query(`SELECT users.username FROM users
	 INNER JOIN frontpage ON users.id = frontpage.user_id`)
	if err != nil {
		return fpagers, err
	}
	defer rows.Close()
	usrs := make([]User, 0)

	for rows.Next() {
		usr := User{}
		err := rows.Scan(&usr.Username)
		if err != nil {
			return usrs, err
		}
		usrs = append(usrs, usr)
	}
	return usrs, nil
}

// GetFrontpager gets ONE Frontpage streamer.
func GetFrontpager() (User, error) {
	usr := User{}
	err := config.DB.QueryRow(`SELECT users.id FROM users
	 INNER JOIN frontpage ON users.id = frontpage.user_id ORDER BY id DESC LIMIT 1`).
		Scan(&usr.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return usr, nil
		}
		panic(err)
	}
	return usr, nil
}

// AddFrontpager adds a streamer to the Frontpage.
func AddFrontpager(UID int) {
	_, err := config.DB.Exec("INSERT INTO frontpage (user_id) VALUES ($1)", UID)
	if err != nil {
		panic(err)
	}
}

// RemoveFrontpager adds a streamer to the Frontpage.
func RemoveFrontpager(UID int) {
	_, err := config.DB.Exec("DELETE FROM frontpage WHERE user_id = $1", UID)
	if err != nil {
		panic(err)
	}
}

// DeleteStreamByUID deletes a stream by username.
func DeleteStreamByUID(UID int) {
	_, err := config.DB.Exec("DELETE FROM streams WHERE user_id = $1", UID)
	if err != nil {
		panic(err)
	}
}
