package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/comebacknader/wink/config"
)

type User struct {
	ID              int       `json:"id,omitempty"`
	Username        string    `json:"username,omitempty"`
	Email           string    `json:"email,omitempty"`
	Hash            string    `json:"password,omitempty"`
	UserType        string    `json:"userType,omitempty"`
	ResetPassToken  string    `json:"resetPassToken,omitempty"`
	ResetPassExpiry time.Time `json:"-"`
	Coins           int       `json:"coins,omitempty"`
	PurchToken      string    `json:"purchToken,omitempty"`
}

// GetAllUsers gets all the users in the database.
func GetAllUsers() ([]User, error) {
	rows, err := config.DB.Query("SELECT ID, username, email, usertype FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	usrs := make([]User, 0)

	for rows.Next() {
		usr := User{}
		err := rows.Scan(&usr.ID, &usr.Username, &usr.Email, &usr.UserType)
		if err != nil {
			panic(err)
		}
		usrs = append(usrs, usr)
	}

	return usrs, nil
}

// Posts a User to the Database
func PostUser(usr User) error {
	_, err := config.DB.
		Exec("INSERT INTO users (username, email, hash, usertype, coins) VALUES ($1, $2, $3, $4, $5)",
			usr.Username, usr.Email, usr.Hash, usr.UserType, 0)
	if err != nil {
		panic(err)
	}
	return nil
}

// Deletes a User from the Database
func DeleteUser(username string) error {
	_, err := config.DB.Exec("DELETE FROM users WHERE username = $1", username)
	if err != nil {
		panic(err)
	}
	return nil
}

// Check if a User exists with supplied Username
// Return true if exists, false otherwise
func CheckUserName(usr string) bool {
	user := User{}
	err := config.DB.QueryRow("SELECT email FROM users WHERE username = $1", usr).Scan(&user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err)
	}
	fmt.Println("Got here.")
	return true
}

// Check if a User exists with supplied Email
// Return true if exists, false otherwise
func CheckUserEmail(email string) bool {
	user := User{}
	err := config.DB.QueryRow("SELECT username FROM users WHERE email = $1", email).Scan(&user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err)
	}
	return true
}

// GetUserByEmail gets the user by the supplied email. Return bool if exists.
func GetUserByEmail(email string) (User, bool) {
	usr := User{}
	err := config.DB.
		QueryRow("SELECT ID, username, email, hash, usertype FROM users WHERE email = $1", email).
		Scan(&usr.ID, &usr.Username, &usr.Email, &usr.Hash, &usr.UserType)
	if err != nil {
		if err == sql.ErrNoRows {
			return usr, false
		}
		panic(err)
	}
	return usr, true
}

// GetUserByName gets the user by their username. Return bool if exists.
func GetUserByName(username string) (User, bool) {
	usr := User{}
	err := config.DB.
		QueryRow("SELECT ID, username, email, hash, usertype FROM users WHERE username = $1", username).
		Scan(&usr.ID, &usr.Username, &usr.Email, &usr.Hash, &usr.UserType)
	if err != nil {
		if err == sql.ErrNoRows {
			return usr, false
		}
		panic(err)
	}
	return usr, true
}

// GetUserById gets the user by their user ID. Return bool if exists.
func GetUserById(uid int) (User, bool) {
	usr := User{}
	err := config.DB.
		QueryRow("SELECT ID, username, email, hash, usertype, coins FROM users WHERE ID = $1", uid).
		Scan(&usr.ID, &usr.Username, &usr.Email, &usr.Hash, &usr.UserType, &usr.Coins)
	if err != nil {
		if err == sql.ErrNoRows {
			return usr, false
		}
		panic(err)
	}
	return usr, true
}

// UserExistById returns whether user exists by supplied user ID.
func UserExistById(uid int) bool {
	usr := User{}
	err := config.DB.
		QueryRow("SELECT ID FROM users WHERE ID = $1", uid).
		Scan(&usr.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		panic(err)
	}
	return true
}

// StoreTokenAndExpiry updates the reset password token and expiration date.
func StoreTokenAndExpiry(usr User, token string, tme time.Time) {
	_, err := config.DB.
		Exec("UPDATE users SET resetpasstoken = $1, resetpassexpiry = $2 WHERE username = $3",
			token, tme, usr.Username)
	if err != nil {
		panic(err)
	}
}

// GetUserByToken gets the user by the supplied password reset token.
func GetUserByToken(token string) (User, bool) {
	usr := User{}
	err := config.DB.
		QueryRow("SELECT ID, username, usertype, email, resetpassexpiry FROM users WHERE resetpasstoken = $1", token).
		Scan(&usr.ID, &usr.Username, &usr.UserType, &usr.Email, &usr.ResetPassExpiry)
	if err != nil {
		if err == sql.ErrNoRows {
			return usr, false
		}
		panic(err)
	}
	return usr, true
}

// UpdateUserPassword updates the password of a user.
func UpdateUserPassword(email string, password string) {
	_, err := config.DB.
		Exec("UPDATE users SET hash = $1 WHERE email = $2",
			password, email)
	if err != nil {
		panic(err)
	}
}

// UpdateUserToStreamer updates a user's type to streamer.
func UpdateUserToStreamer(usrid int) {
	_, err := config.DB.
		Exec("UPDATE users SET usertype = $1 WHERE id = $2",
			"streamer", usrid)
	if err != nil {
		panic(err)
	}
	_, err = config.DB.
		Exec(`INSERT INTO streams (user_id, title, game, withgame, 
			online, twitter, siteone) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			usrid, "Enter Title", "Enter Game", true, false,
			"Enter Twitter", "Enter Personal Website")
	if err != nil {
		panic(err)
	}
}

// UpdateUserToAdmin updates a user's type to admin.
func UpdateUserToAdmin(usrid int) {
	_, err := config.DB.
		Exec("UPDATE users SET usertype = $1 WHERE id = $2",
			"admin", usrid)
	if err != nil {
		panic(err)
	}
}

// UpdateAdminToUser updates an admin's type to user.
func UpdateAdminToUser(usrid int) {
	_, err := config.DB.
		Exec("UPDATE users SET usertype = $1 WHERE id = $2",
			"user", usrid)
	if err != nil {
		panic(err)
	}
}

// UpdateStreamerToUser updates a streamer's type to user.
func UpdateStreamerToUser(usrid int) {
	_, err := config.DB.
		Exec("UPDATE users SET usertype = $1 WHERE id = $2",
			"user", usrid)
	if err != nil {
		panic(err)
	}
	_, err = config.DB.
		Exec("DELETE FROM streams WHERE user_id = $1", usrid)
	if err != nil {
		panic(err)
	}
}

// GetTokenByEmail gets the reset password token by email.
func GetTokenByEmail(email string) (string, error) {
	usr := User{}
	err := config.DB.
		QueryRow("SELECT resetpasstoken FROM users WHERE email = $1", email).
		Scan(&usr.ResetPassToken)
	if err != nil {
		return "", err
	}
	return usr.ResetPassToken, nil
}

// GetPurchToken gets the purchase token needed to buy coins.
func GetPurchToken(username string) (string, bool) {
	usr := User{}
	err := config.DB.
		QueryRow("SELECT purchtoken FROM users WHERE username = $1", username).
		Scan(&usr.PurchToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false
		}
		panic(err)
	}
	return usr.PurchToken, true
}

// ResetPurchToken resets the user's purchase token.
func ResetPurchToken(username string, token string) {
	_, err := config.DB.
		Exec("UPDATE users SET purchtoken = $1 WHERE username = $2",
			token, username)
	if err != nil {
		panic(err)
	}
}
