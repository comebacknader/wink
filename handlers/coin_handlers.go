package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/comebacknader/wink/config"
	"github.com/comebacknader/wink/models"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
)

type Tip struct {
	Amt  int    `json:"amt, omitempty"`
	Sndr string `json:"sender, omitempty"`
}

type CoinData struct {
	UserCred  UserStatus
	CurrUser  string
	CurrStrmr string
	Purchases []models.Purchase
	PrchToken string
	Coins     int
	Error     string
	Success   string
}

// GetBuyCoins gets the page to buy coins.
func GetBuyCoins(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if req.Method != http.MethodGet {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	cnData := CoinData{}
	usr, loggedIn := GetCurrentUser(req)

	if loggedIn == false {
		http.Redirect(w, req, "/", http.StatusUnauthorized)
		return
	}

	cnData.UserCred.IsLogIn = true
	cnData.CurrUser = usr.Username

	if usr.UserType == "admin" {
		cnData.UserCred.IsAdmin = true
	}

	if usr.UserType == "streamer" {
		cnData.UserCred.IsStreamer = true
	} else {
		cnData.UserCred.IsStreamer = false
	}

	// Need to get the PurchToken of the User
	newTkn := uuid.Must(uuid.NewV4()).String()
	models.ResetPurchToken(usr.Username, newTkn)

	cnData.PrchToken = newTkn
	cnData.Coins = usr.Coins

	// Execute Template
	err := tpl.ExecuteTemplate(w, "buycoins.gohtml", cnData)
	if err != nil {
		config.HandleError(w, err)
	}
	return
}

// PostTip sends a tip from a user to a streamer.
func PostTip(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	// Check if the sender is already logged in
	usr, loggedIn := GetCurrentUser(req)
	if loggedIn == false {
		http.Redirect(w, req, "/", http.StatusUnauthorized)
		return
	}

	// Get the values of the submission
	sndr := usr
	rcvr := req.FormValue("receiver")
	strAmount := req.FormValue("amount")
	amt, err := strconv.Atoi(strAmount)
	if err != nil {
		http.Error(w, "Error with amount sent.", http.StatusMethodNotAllowed)
		return
	}

	// Check if sender has enough coins to send
	if sndr.Coins < amt {
		http.Error(w, "You do not have enough coins to give.", 400)
		return
	}

	// Subtract coins from sender
	models.SubCoinsFromUser(sndr.Username, amt)
	// Add coins to receiver
	models.AddCoinsToUser(rcvr, amt)

	// Return success message in json
	enc := json.NewEncoder(w)
	encErr := enc.Encode(Tip{Amt: amt, Sndr: sndr.Username})
	if encErr != nil {
		panic(encErr)
	}
	return
}

// PostPurchase allows a user to purchase points.
// The request is sent from ccBill's servers.
func PostPurchase(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	hookType := req.FormValue("eventType")
	if hookType == "NewSaleSuccess" {
		splitAddr := strings.Split(req.RemoteAddr, ":")
		if len(splitAddr) > 2 {
			fmt.Println("IP: Loopback Address")
			return
		}
		ip := net.ParseIP(splitAddr[0])
		_, range1, _ := net.ParseCIDR("64.38.240.0/24")
		_, range2, _ := net.ParseCIDR("64.38.241.0/24")
		_, range3, _ := net.ParseCIDR("64.38.212.0/24")
		_, range4, _ := net.ParseCIDR("64.38.215.0/24")

		if !range1.Contains(ip) && !range2.Contains(ip) &&
			!range3.Contains(ip) && !range4.Contains(ip) {
			//Return error
			http.Error(w, "You are not authorized to do that.", 401)
			fmt.Println("IP not in proper range:", ip)
			return
		}
		// Make sure it has form values
		username := req.FormValue("X-uname")
		price := req.FormValue("billedInitialPrice")
		token := req.FormValue("X-token")

		purchToken, _ := models.GetPurchToken(username)
		// Make sure purchaseToken === token submitted via form.
		if purchToken != token {
			http.Error(w, "You are not authorized to do that.", 401)
			return
		}
		if price == "10.99" {
			models.AddCoinsToUser(username, 100)
		} else if price == "20.99" {
			models.AddCoinsToUser(username, 200)
		} else if price == "44.99" {
			models.AddCoinsToUser(username, 500)
		} else if price == "62.99" {
			models.AddCoinsToUser(username, 750)
		} else if price == "79.99" {
			models.AddCoinsToUser(username, 1000)
		} else {
			http.Error(w, "Price is invalid.", 403)
			return
		}
		newTkn := uuid.Must(uuid.NewV4()).String()
		models.ResetPurchToken(username, newTkn)
		// Add Purchase to Purchases Table
		bought := time.Now().UTC()
		models.AddPurchase(username, price, "coins", bought)
		w.WriteHeader(200)
		return
	}
	fmt.Println("EventType didn't match.")
	return
}

// GetPurchases lists all the purchases.
func GetPurchases(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if req.Method != http.MethodGet {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	RedirectNotAdmin(w, req)
	usr, loggedIn := GetCurrentUser(req)
	if loggedIn == false {
		http.Redirect(w, req, "/", http.StatusUnauthorized)
		return
	}
	tplData := CoinData{}
	tplData.UserCred.IsLogIn = true
	tplData.UserCred.IsAdmin = true
	tplData.CurrUser = usr.Username

	// Get list of purchases.
	prchs, err := models.GetAllPurchases()

	switch {
	case err == sql.ErrNoRows:
		tplData.Purchases = prchs
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	tplData.Purchases = prchs

	err = tpl.ExecuteTemplate(w, "purchases.gohtml", tplData)
	if err != nil {
		panic(err)
	}
	return
}

// DenyPurch gets called when ccBill purchase fails.
func DenyPurch(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	ip := req.FormValue("ipAddress")
	fmt.Println("Purchase Denied by ", ip)
}
