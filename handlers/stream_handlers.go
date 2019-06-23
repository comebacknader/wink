package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/comebacknader/wink/config"
	"github.com/comebacknader/wink/models"
	"github.com/julienschmidt/httprouter"
)

// StreamerData holds information about Streamer
type StreamerData struct {
	UserCred  UserStatus
	CurrStrmr string
	CurrUser  string
	TotCoins  int
	Stream    models.Stream
	Error     string
	Success   string
}

type StreamUpdate struct {
	Title   string `json:"title,omitempty"`
	Game    string `json:"game, omitempty"`
	Twit    string `json:"twitter, omitempty"`
	SiteOne string `json:"siteone, omitempty"`
}

// GetDashboard gets the streamer's dashboard page.
func GetDashboard(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if req.Method != http.MethodGet {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	usr, loggedIn := GetCurrentUser(req)
	if loggedIn == false {
		http.Redirect(w, req, "/", http.StatusUnauthorized)
		return
	}
	if usr.UserType != "streamer" {
		http.Redirect(w, req, "/", http.StatusUnauthorized)
		return
	}

	usrname := ps.ByName("username")
	if usr.Username != usrname {
		http.Redirect(w, req, "/", http.StatusUnauthorized)
		return
	}

	// Get the Stream of the user
	strm, strmExist := models.GetStreamByUID(usr.ID)
	if strmExist == false {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	// Store necessary data in StreamerData
	strmData := StreamerData{}
	strmData.CurrStrmr = usr.Username
	if strm.Twit == "Enter Twitter" {
		strm.Twit = ""
	}
	if strm.SiteOne == "Enter Personal Website" {
		strm.SiteOne = ""
	}
	strmData.Stream = strm
	strmData.UserCred.IsLogIn = true
	strmData.UserCred.IsStreamer = true
	strmData.CurrUser = usr.Username
	strmData.TotCoins = usr.Coins
	// Execute Template
	err := tpl.ExecuteTemplate(w, "dashboard.gohtml", strmData)
	if err != nil {
		panic(err)
	}
	return
}

// UpdateStream updates the streamer's stream.
func UpdateStream(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	c, err := req.Cookie("session")
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	usrId := models.GetUserIDByCookie(c.Value)
	usr, _ := models.GetUserById(usrId)
	if usr.UserType != "streamer" {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	models.UpdateSessionActivity(c.Value)
	title := req.FormValue("title")
	game := req.FormValue("game")
	twit := req.FormValue("twitter")
	siteone := req.FormValue("siteone")

	if title == "" {
		http.Error(w, "Title cannot be blank.", 400)
		return
	}
	if len(title) > 100 {
		http.Error(w, "Title too long. 100 character max.", 400)
		return
	}
	if game == "" {
		http.Error(w, "Game title cannot be blank.", 400)
		return
	}
	if len(game) > 50 {
		http.Error(w, "Game title too long. 50 character max.", 400)
		return
	}
	if twit == "" {
		twit = "Enter Twitter"
	}
	if siteone == "" {
		siteone = "Enter Personal Website"
	}
	// Update the stream information
	models.UpdateStream(title, game, twit, siteone, usrId)
	stream := StreamUpdate{title, game, twit, siteone}
	jsonStrm, err := json.Marshal(stream)
	if err != nil {
		http.Error(w, "Streams cannot be updated at this time.", 500)
		return
	}

	_, writeErr := fmt.Fprint(w, string(jsonStrm))
	if writeErr != nil {
		http.Error(w, "There was a problem writing json to writer.", 500)
		return
	}
	return
}

// UpdateOnline updates the status of the stream being offline/online.
func UpdateOnline(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	c, err := req.Cookie("session")
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	usrId := models.GetUserIDByCookie(c.Value)
	usr, _ := models.GetUserById(usrId)
	if usr.UserType != "streamer" {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	models.UpdateSessionActivity(c.Value)
	online := req.FormValue("online")
	if online == "" {
		http.Error(w, "Online cannot be blank.", 400)
		return
	}
	if online == "online" {
		models.StreamOnline(usrId)
		fmt.Fprint(w, `{"status":"online"}`)
		return
	} else {
		models.StreamOffline(usrId)
		fmt.Fprint(w, `{"status":"offline"}`)
		return
	}
}

// GetStreaming gets the streamer's public-facing streaming page.
func GetStreaming(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if req.Method != http.MethodGet {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	strmData := StreamerData{}
	currUser, loggedIn := GetCurrentUser(req)
	if loggedIn == true {
		strmData.UserCred.IsLogIn = true
		if currUser.UserType == "streamer" {
			strmData.UserCred.IsStreamer = true
		}
		if currUser.UserType == "admin" {
			strmData.UserCred.IsAdmin = true
		}
	} else {
		strmData.UserCred.IsLogIn = false
	}

	username := ps.ByName("username")
	usr, exist := models.GetUserByName(username)
	if exist == false {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	// Get the Stream of the user
	strm, strmExist := models.GetStreamByUID(usr.ID)
	if strmExist == false {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	if strm.Twit == "Enter Twitter" {
		strm.Twit = ""
	}
	if strm.SiteOne == "Enter Personal Website" {
		strm.SiteOne = ""
	}
	strmData.Stream = strm
	strmData.CurrStrmr = username
	strmData.CurrUser = currUser.Username
	strmData.TotCoins = currUser.Coins
	err := tpl.ExecuteTemplate(w, "streaming.gohtml", strmData)
	if err != nil {
		panic(err)
	}
	return
}

// GetFrontpage gets the page where admins determine frontpage streamer.
func GetFrontpage(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if req.Method != http.MethodGet {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	usr, loggedIn := GetCurrentUser(req)
	RedirectNotAdmin(w, req)
	adminData := UserData{}
	if loggedIn == true {
		adminData.UserCred.IsLogIn = true
	}
	if usr.UserType == "admin" {
		adminData.UserCred.IsAdmin = true
	}

	// Get a list of frontpage streamers.
	frontStrmrs, error := models.GetFrontpagers()

	switch {
	case error == sql.ErrNoRows:
		frontStrmrs = []models.User{}
	case error != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	adminData.Users = frontStrmrs

	err := tpl.ExecuteTemplate(w, "frontpage.gohtml", adminData)
	if err != nil {
		config.HandleError(w, err)
	}
	return
}

// PostFrontpage adds/removes a streamer from the Frontpage list.
func PostFrontpage(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	loggedIn := AlreadyLoggedIn(req)
	RedirectNotAdmin(w, req)
	name := req.FormValue("streamer")
	action := req.FormValue("addremove")
	usr, _ := models.GetUserByName(name)

	adminData := UserData{}
	if loggedIn == true {
		adminData.UserCred.IsLogIn = true
	}

	if usr.UserType != "streamer" {
		w.WriteHeader(400)
		adminData.Error = "User is not a streamer."

		err := tpl.ExecuteTemplate(w, "frontpage.gohtml", adminData)
		if err != nil {
			config.HandleError(w, err)
		}
		return
	}
	if action == "add" {
		models.AddFrontpager(usr.ID)
	}
	if action == "remove" {
		models.RemoveFrontpager(usr.ID)
	}

	frontStrmrs, error := models.GetFrontpagers()

	switch {
	case error == sql.ErrNoRows:
		frontStrmrs = []models.User{}
	case error != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	adminData.Users = frontStrmrs
	w.WriteHeader(200)
	adminData.Success = "Successfully updated."

	err := tpl.ExecuteTemplate(w, "frontpage.gohtml", adminData)
	if err != nil {
		config.HandleError(w, err)
	}
	return
}
