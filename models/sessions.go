package models

import (
	"database/sql"
	"time"

	"github.com/comebacknader/wink/config"
)

type Session struct {
	ID       int    `json:"id,omitempty"`
	UUID     string `json:"uuid,omitempty"`
	UserID   int    `json:"user_id,omitempty"`
	UserType string `json:"userType,omitempty"`
}

type UserSession struct {
	Credential string `json:"credential,omitempty"`
	Hash       string `json:"password,omitempty"`
}

// CreateSession creates a session.
func CreateSession(usrID int, userType string, sID string, activeTime time.Time) {
	_, err := config.DB.
		Exec("INSERT INTO sessions (uuid, user_id, usertype, activity) VALUES ($1, $2, $3, $4)",
			sID, usrID, userType, activeTime)
	if err != nil {
		panic(err)
	}
}

// GetUserIDByCookie gets the user ID from a supplied Cookie.
func GetUserIDByCookie(cookie string) int {
	sesh := Session{}
	err := config.DB.QueryRow("SELECT user_id FROM sessions WHERE uuid = $1", cookie).
		Scan(&sesh.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		}
		panic(err)
	}
	return sesh.UserID
}

// DelSessionByUUID deletes a session by the supplied uuid.
func DelSessionByUUID(uuid string) {
	_, err := config.DB.
		Exec("DELETE FROM sessions WHERE uuid = $1", uuid)
	if err != nil {
		panic(err)
	}
}

// Deletes Session by Username
func DelSessionByUsername(username string) {
	user := User{}
	err := config.DB.QueryRow("SELECT ID FROM users WHERE username = $1", username).
		Scan(&user.ID)
	_, err2 := config.DB.
		Exec("DELETE FROM sessions WHERE user_id = $1",
			user.ID)
	if err != nil {
		if err == sql.ErrNoRows {

		} else {
			panic(err)
		}
	}
	if err2 != nil {
		panic(err)
	}
}

// UpdateSessionActivity updates a session's activity.
func UpdateSessionActivity(uuid string) {
	newTime := time.Now().UTC()
	_, err := config.DB.
		Exec("UPDATE sessions SET activity = $1 WHERE uuid = $2", newTime, uuid)
	if err != nil {
		panic(err)
	}
}

// DeleteOldSessions deletes sessions that are 4+ hours old.
func DeleteOldSessions() {
	_, err := config.DB.
		Exec("DELETE FROM sessions WHERE activity < now() - interval '4 hours'")
	if err != nil {
		panic(err)
	}
}
