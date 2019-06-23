package models

import (
	"time"

	"github.com/comebacknader/wink/config"
)

type Purchase struct {
	ID       int       `json:"id,omitempty"`
	UserId   int       `json:"user_id,omitempty"`
	Typ      string    `json:"type,omitempty"`
	Price    string    `json:"price,omitempty"`
	Username string    `json:"username,omitempty"`
	Bought   time.Time `json:"bought,omitempty"`
}

// AddPurchase adds a new purchase.
func AddPurchase(username string, price string, typ string, bought time.Time) {
	// Get User_ID by username
	usr, _ := GetUserByName(username)
	// Exec purchase
	_, err := config.DB.
		Exec("INSERT INTO purchases (user_id, type, price, username, bought) VALUES ($1, $2, $3, $4, $5)",
			usr.ID, typ, price, username, bought)
	if err != nil {
		panic(err)
	}
}

// GetAllPurchases gets all the purchases.
func GetAllPurchases() ([]Purchase, error) {
	prchs := []Purchase{}
	rows, err := config.DB.Query("SELECT * FROM purchases")
	if err != nil {
		return prchs, err
	}
	defer rows.Close()
	for rows.Next() {
		purch := Purchase{}
		err := rows.Scan(&purch.ID, &purch.UserId, &purch.Typ, &purch.Price,
			&purch.Username, &purch.Bought)
		if err != nil {
			panic(err)
		}
		prchs = append(prchs, purch)
	}
	return prchs, nil
}
