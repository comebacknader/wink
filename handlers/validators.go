package handlers

import (
	"net/http"

	"github.com/badoux/checkmail"
	"github.com/comebacknader/wink/config"
	"github.com/comebacknader/wink/models"
)

// Helper function to validate User fields
// return 1 is error
// return 0 is non-error
func ValidateUserFields(w http.ResponseWriter, usr models.User, errors SessionData) int {
	// errors := CredErrors{}
	// errors.Error = "Username cannot be blank"
	if usr.Username == "" {
		errors.Error = append(errors.Error, "Username cannot be blank")
		w.WriteHeader(400)
		err := tpl.ExecuteTemplate(w, "signup.gohtml", errors)
		config.HandleError(w, err)
		return 0
	}

	if len(usr.Username) < 6 || len(usr.Username) > 30 {
		errors.Error = append(errors.Error, "Username must be between 6 and 30 characters.")
		w.WriteHeader(400)
		err := tpl.ExecuteTemplate(w, "signup.gohtml", errors)
		config.HandleError(w, err)
		return 0
	}

	if usr.Email == "" {
		errors.Error = append(errors.Error, "Email cannot be blank.")
		w.WriteHeader(400)
		err := tpl.ExecuteTemplate(w, "signup.gohtml", errors)
		config.HandleError(w, err)
		return 0
	}

	if len(usr.Email) < 8 || len(usr.Email) > 40 {
		errors.Error = append(errors.Error, "Email must be between 7 and 40 characters.")
		w.WriteHeader(400)
		err := tpl.ExecuteTemplate(w, "signup.gohtml", errors)
		config.HandleError(w, err)
		return 0
	}

	err := checkmail.ValidateFormat(usr.Email)
	if err != nil {
		errors.Error = append(errors.Error, "Email is not correct format.")
		w.WriteHeader(400)
		err := tpl.ExecuteTemplate(w, "signup.gohtml", errors)
		config.HandleError(w, err)
		return 0
	}

	if usr.Hash == "" {
		errors.Error = append(errors.Error, "Password cannot be blank.")
		w.WriteHeader(400)
		err := tpl.ExecuteTemplate(w, "signup.gohtml", errors)
		config.HandleError(w, err)
		return 0
	}

	if len(usr.Hash) < 6 || len(usr.Hash) > 50 {
		errors.Error = append(errors.Error, "Password must be between 6 and 50 characters.")
		w.WriteHeader(400)
		err := tpl.ExecuteTemplate(w, "signup.gohtml", errors)
		config.HandleError(w, err)
		return 0
	}

	return 1
}
