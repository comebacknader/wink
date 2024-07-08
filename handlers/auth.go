package handlers

import (
	"net/http"

	"github.com/comebacknader/wink/models"
)

// AlreadyLoggedIn determines whether the User is already logged in.
func AlreadyLoggedIn(req *http.Request) bool {
	c, err := req.Cookie("session")
	if err != nil {
		return false
	}
	// Get a User ID from the cookie's value
	// Checks session database and returns User ID
	uid := models.GetUserIDByCookie(c.Value)
	if uid == 0 {
		return false
	}
	// Make sure User exists with User ID
	exist := models.UserExistById(uid)

	// Update session's activity
	models.UpdateSessionActivity(c.Value)

	return exist
}

// RedirectNotAdmin is a helper function that redirects user if not admin.
func RedirectNotAdmin(w http.ResponseWriter, req *http.Request) {
	// Redirect if not admin
	c, err := req.Cookie("session")
	if err != nil {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	usrId := models.GetUserIDByCookie(c.Value)
	// Should be GetUserTypeById(usrId)
	usr, _ := models.GetUserById(usrId)
	if usr.UserType != "admin" {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
}

// GetCurrentUser gets the current user logged in.
func GetCurrentUser(req *http.Request) (models.User, bool) {
	nilUser := models.User{}
	c, err := req.Cookie("session")
	if err != nil {
		return nilUser, false
	}
	uid := models.GetUserIDByCookie(c.Value)
	if uid == 0 {
		return nilUser, false
	}
	user, exist := models.GetUserById(uid)
	if exist == false {
		return nilUser, false
	}
	// Update session's activity
	models.UpdateSessionActivity(c.Value)
	return user, true
}
