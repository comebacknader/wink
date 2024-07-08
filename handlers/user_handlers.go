package handlers

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/comebacknader/wink/config"
	"github.com/comebacknader/wink/models"
	"github.com/julienschmidt/httprouter"
)

var tpl *template.Template

func init() {
	tpl = config.Tpl
}

type UserStatus struct {
	IsLogIn    bool
	IsStreamer bool
	IsAdmin    bool
}

type UserData struct {
	UserCred  UserStatus
	CurrUser  string
	CurrStrmr string
	Users     []models.User
	User      models.User
	Error     string
	Success   string
}

// ShowAllUsers lists all the users.
func ShowAllUsers(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if req.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	RedirectNotAdmin(w, req)
	usr, loggedIn := GetCurrentUser(req)

	data := UserData{}
	if loggedIn == true {
		data.UserCred.IsLogIn = true
	} else {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	usrs, err := models.GetAllUsers()

	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, req)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	data.Users = usrs
	data.CurrUser = usr.Username
	data.UserCred.IsStreamer = false
	data.UserCred.IsAdmin = true

	tpl.ExecuteTemplate(w, "users_all.gohtml", data)
	return
}

// UpdateUserType updates the UserType of the user.
func UpdateUserType(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	RedirectNotAdmin(w, req)
	username := req.FormValue("username")
	newtype := req.FormValue("usertype")

	tplData := UserData{}
	tplData.UserCred.IsLogIn = true
	// Need to populate users when re-rendering template
	usrs, err := models.GetAllUsers()

	tplData.Users = usrs

	// Get information on user (userType)
	usr, exists := models.GetUserByName(username)
	if exists != true {
		w.WriteHeader(400)
		tplData.Error = "No username by that name."
		err = tpl.ExecuteTemplate(w, "users_all.gohtml", tplData)
		if err != nil {
			panic(err)
		}
		return
	}

	if usr.UserType == "user" && newtype == "streamer" {
		models.UpdateUserToStreamer(usr.ID)
	} else if usr.UserType == "user" && newtype == "admin" {
		models.UpdateUserToAdmin(usr.ID)
	} else if usr.UserType == "admin" && newtype == "user" {
		models.UpdateAdminToUser(usr.ID)
	} else if usr.UserType == "streamer" && newtype == "user" {
		models.UpdateStreamerToUser(usr.ID)
	} else if usr.UserType == "admin" && newtype == "streamer" {
		// Same thing as going from user --> streamer
		models.UpdateUserToStreamer(usr.ID)
	} else {
		w.WriteHeader(400)
		tplData.Error = "Can't update user with parameters."
		err := tpl.ExecuteTemplate(w, "users_all.gohtml", tplData)
		if err != nil {
			panic(err)
		}
		return
	}

	usrs, err = models.GetAllUsers()
	if err != nil {
		panic(err)
		w.WriteHeader(500)
		tplData.Error = "There was a problem updating."
		err := tpl.ExecuteTemplate(w, "users_all.gohtml", tplData)
		if err != nil {
			panic(err)
		}
		return
	}

	tplData.Users = usrs
	w.WriteHeader(200)
	tplData.Success = "Successfully updated."
	err = tpl.ExecuteTemplate(w, "users_all.gohtml", tplData)
	if err != nil {
		panic(err)
	}
	return

}
