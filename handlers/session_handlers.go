package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/badoux/checkmail"
	"github.com/comebacknader/wink/config"
	"github.com/comebacknader/wink/models"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// Data passed to all pages dealing with Sessions
type SessionData struct {
	UserCred UserStatus
	CurrUser string
	Error    []string
	Success  string
	Token    string
}

func GetLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	loggedIn := AlreadyLoggedIn(r)
	if loggedIn == true {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		seshData := SessionData{}
		seshData.UserCred.IsLogIn = false
		err := tpl.ExecuteTemplate(w, "login.gohtml", seshData)
		config.HandleError(w, err)
	}
	return
}

// GetSignup gets the signup page.
func GetSignup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	loggedIn := AlreadyLoggedIn(r)
	if loggedIn == true {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		seshData := SessionData{}
		seshData.UserCred.IsLogIn = false
		err := tpl.ExecuteTemplate(w, "signup.gohtml", seshData)
		config.HandleError(w, err)
	}
}

// PostLogin logs a user in.
func PostLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	cred := r.FormValue("credential")
	password := r.FormValue("password")

	// Check if User submitted Username or Email
	var credential string
	var user models.User
	var exist bool
	seshData := SessionData{}
	err := checkmail.ValidateFormat(cred)
	if err == nil {
		user, exist = models.GetUserByEmail(cred)
		if exist == false {
			seshData.Error = append(seshData.Error, "Sorry, that email doesn't exist.")
			w.WriteHeader(400)
			err := tpl.ExecuteTemplate(w, "login.gohtml", seshData)
			config.HandleError(w, err)
			return
		}
		credential = "Email"

	} else {
		user, exist = models.GetUserByName(cred)
		if exist == false {
			seshData.Error = append(seshData.Error, "Sorry, that username doesn't exist.")
			w.WriteHeader(400)
			err := tpl.ExecuteTemplate(w, "login.gohtml", seshData)
			config.HandleError(w, err)
			return
		}
		credential = "Username"
	}

	// Compare user submitted password to password in DB
	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
	if err != nil {
		seshData.Error = append(seshData.Error, credential+" and/or password do not match.")
		w.WriteHeader(400)
		err := tpl.ExecuteTemplate(w, "login.gohtml", seshData)
		config.HandleError(w, err)
		return
	}

	// Create Session and store it in Cookie and DB
	sID := uuid.Must(uuid.NewV4())
	activeTime := time.Now().UTC()

	models.CreateSession(user.ID, user.UserType, sID.String(), activeTime)

	cookie := &http.Cookie{
		Name:     "session",
		Value:    sID.String(),
		Expires:  time.Now().UTC().Add(time.Hour * 24),
		MaxAge:   86400,
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return

}

// PostSignup signs up a user.
func PostSignup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	usr := models.User{}
	usr.Username = r.FormValue("username")
	usr.Email = r.FormValue("email")
	usr.Hash = r.FormValue("password")

	seshData := SessionData{}

	valErr := ValidateUserFields(w, usr, seshData)
	if valErr == 0 {
		return
	}

	doesNameExist := models.CheckUserName(usr.Username)
	if doesNameExist == true {
		seshData.Error = append(seshData.Error, "Username already taken.")
		w.WriteHeader(400)
		err := tpl.ExecuteTemplate(w, "signup.gohtml", seshData)
		config.HandleError(w, err)
		return
	}

	fmt.Println("Got here!!")

	doesEmailExist := models.CheckUserEmail(usr.Email)
	if doesEmailExist == true {
		seshData.Error = append(seshData.Error, "Email already taken.")
		w.WriteHeader(400)
		err := tpl.ExecuteTemplate(w, "signup.gohtml", seshData)
		config.HandleError(w, err)
		return
	}

	usr.UserType = "user"

	// Encrypt the password using Bcrypt
	hashPass, err := bcrypt.GenerateFromPassword([]byte(usr.Hash), bcrypt.MinCost)

	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	usr.Hash = string(hashPass[:])

	err = models.PostUser(usr)

	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	user, _ := models.GetUserByName(usr.Username)

	sID := uuid.Must(uuid.NewV4())
	cookie := &http.Cookie{
		Name:     "session",
		Value:    sID.String(),
		Path:     "/",
		Expires:  time.Now().UTC().Add(time.Hour * 24),
		MaxAge:   86400,
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	activeTime := time.Now().UTC()

	models.CreateSession(user.ID, user.UserType, sID.String(), activeTime)

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

// User Logout
func PostLogout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	expCook := &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().UTC(),
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, expCook)
	models.DelSessionByUUID(cookie.Value)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// DelOldSessions deletes all of the old sessions.
func DelOldSessions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	loggedIn := AlreadyLoggedIn(r)
	if loggedIn == false {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	c, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	usrId := models.GetUserIDByCookie(c.Value)
	usr, _ := models.GetUserById(usrId)
	if usr.UserType != "admin" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	models.DeleteOldSessions()
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

// ForgotPass is the GET handler for forgotten passwords.
func ForgotPass(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if req.Method != http.MethodGet {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	seshData := SessionData{}
	err := tpl.ExecuteTemplate(w, "forgot.gohtml", seshData)
	if err != nil {
		panic(err)
		return
	}
	return
}

// PostForgotPass is the POST handler for forgotten passwords.
func PostForgotPass(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// Store now.UTC + 4 hours into the database
	if req.Method != http.MethodPost {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	email := req.FormValue("email")
	usr, exist := models.GetUserByEmail(email)
	seshData := SessionData{}
	if exist == false {
		w.WriteHeader(400)
		seshData.Error = append(seshData.Error, "No account with that email.")
		err := tpl.ExecuteTemplate(w, "forgot.gohtml", seshData)
		if err != nil {
			panic(err)
		}
		return
	}

	tkn := uuid.Must(uuid.NewV4()).String()
	tme := time.Now().UTC().Add(time.Hour * time.Duration(2))
	models.StoreTokenAndExpiry(usr, tkn, tme)

	host := os.Getenv("WINK_MAIL_HOST")
	usrname := os.Getenv("WINK_MAIL_U")
	pass := os.Getenv("WINK_MAIL_P")
	auth := smtp.PlainAuth("", usrname, pass, host)
	to := []string{usr.Email}
	msg := []byte("From:" + "Wink.GG<admin@wink.gg>" + "\r\n" +
		"To:" + usr.Email + "\r\n" +
		"Subject: Password Reset Link \r\n" +
		"\r\n" +
		"Click on this link to reset your password: " +
		"https://www.wink.gg/reset?token=" + tkn + "\r\n")

	err := smtp.SendMail(host+":587", auth, "admin@wink.gg", to, msg)

	w.WriteHeader(200)
	seshData.Success = "An email was sent with instructions on how to reset your password."
	err = tpl.ExecuteTemplate(w, "forgot.gohtml", seshData)
	if err != nil {
		panic(err)
	}
	return
}

// GetResetPass is the GET handler for resetting password.
func GetResetPass(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if req.Method != http.MethodGet {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	seshData := SessionData{}
	seshData.Token = req.URL.Query().Get("token")
	err := tpl.ExecuteTemplate(w, "reset.gohtml", seshData)
	if err != nil {
		panic(err)
		return
	}
	return
}

// PostResetPass is the POST handler for resetting password.
func PostResetPass(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if req.Method != http.MethodPost {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	pass := req.FormValue("password")
	conf := req.FormValue("confirmation")
	tkn := req.FormValue("token")

	seshData := SessionData{}

	if pass == "" || conf == "" {
		w.WriteHeader(400)
		seshData.Error = append(seshData.Error,
			"Password fields can't be empty.")
		seshData.Token = tkn
		err := tpl.ExecuteTemplate(w, "reset.gohtml", seshData)
		if err != nil {
			panic(err)
		}
		return
	}

	if pass != conf {
		w.WriteHeader(400)
		seshData.Error = append(seshData.Error,
			"Password and password confirmation don't match.")
		seshData.Token = tkn
		err := tpl.ExecuteTemplate(w, "reset.gohtml", seshData)
		if err != nil {
			panic(err)
		}
		return
	}

	if len(pass) < 6 || len(pass) > 50 {
		seshData.Error = append(seshData.Error, "Password must be between 6 and 50 characters.")
		w.WriteHeader(400)
		err := tpl.ExecuteTemplate(w, "reset.gohtml", seshData)
		config.HandleError(w, err)
		return
	}

	// Check that a user exists with token

	usr, exist := models.GetUserByToken(tkn)
	if exist != true {
		w.WriteHeader(400)
		seshData.Error = append(seshData.Error,
			"Your confirmation link is invalid.")
		err := tpl.ExecuteTemplate(w, "forgot.gohtml", seshData)
		if err != nil {
			panic(err)
		}
		return
	}

	timenow := time.Now().UTC()
	expired := usr.ResetPassExpiry.UTC().Before(timenow)

	if expired == true {
		w.WriteHeader(400)
		seshData.Error = append(seshData.Error,
			"Confirmation link has expired.")
		seshData.Token = tkn
		err := tpl.ExecuteTemplate(w, "reset.gohtml", seshData)
		if err != nil {
			panic(err)
		}
		return
	}

	// Encrypt the password using Bcrypt
	hashPass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)

	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	passwd := string(hashPass[:])

	models.UpdateUserPassword(usr.Email, passwd)

	sID := uuid.Must(uuid.NewV4())
	cookie := &http.Cookie{
		Name:     "session",
		Value:    sID.String(),
		Path:     "/",
		Expires:  time.Now().UTC().Add(time.Hour * 24),
		MaxAge:   86400,
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	activeTime := time.Now().UTC()
	models.CreateSession(usr.ID, usr.UserType, sID.String(), activeTime)

	seshData.Success = "Password Successfully Changed"
	seshData.UserCred.IsLogIn = true
	if usr.UserType == "streamer" {
		seshData.UserCred.IsStreamer = true
	}
	err = tpl.ExecuteTemplate(w, "reset.gohtml", seshData)
	if err != nil {
		panic(err)
		return
	}
	return
}
