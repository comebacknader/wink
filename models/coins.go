package models

import (
	"github.com/comebacknader/wink/config"
)

// SubCoinsFromUser subtracts the coins from a User.
func SubCoinsFromUser(username string, amt int) {
	_, err := config.DB.Exec("UPDATE users SET coins = coins - $1 WHERE username = $2",
		amt, username)
	if err != nil {
		panic(err)
	}
}

// AddCoinsToUser subtracts the coins from a User.
func AddCoinsToUser(username string, amt int) {
	_, err := config.DB.Exec("UPDATE users SET coins = coins + $1 WHERE username = $2",
		amt, username)
	if err != nil {
		panic(err)
	}
}
