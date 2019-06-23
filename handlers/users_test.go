package handlers

import (
	"github.com/comebacknader/wink/models"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func CreateMockAdmin(t *testing.T) *http.Response {
	validUser := url.Values{}
	validUser.Add("username", "skylarTest")
	validUser.Add("email", "skylar@test.com")
	validUser.Add("password", "tester")

	req := httptest.NewRequest("POST", "/signup", strings.NewReader(validUser.Encode()))
	req.Form = validUser
	w := httptest.NewRecorder()
	ps := httprouter.Params{}

	PostSignup(w, req, ps)
	resp := w.Result()
	if resp.StatusCode != 303 {
		t.Error("Populating Mock Admin | Expected Status Code: 303, Got: ", resp.StatusCode)
	}

	usr, _ := models.GetUserByName("skylarTest")
	models.UpdateUserToAdmin(usr.ID)
	return resp
}

// CreateMockUser creates and logs in a mock user.
func CreateMockUser(t *testing.T) *http.Response {
	validUser := url.Values{}
	validUser.Add("username", "testuser")
	validUser.Add("email", "test@user.com")
	validUser.Add("password", "tester")

	req := httptest.NewRequest("POST", "/signup", strings.NewReader(validUser.Encode()))
	req.Form = validUser
	w := httptest.NewRecorder()
	ps := httprouter.Params{}

	PostSignup(w, req, ps)
	resp := w.Result()
	if resp.StatusCode != 303 {
		t.Error("Populating Mock User | Expected Status Code: 303, Got: ", resp.StatusCode)
	}
	return resp
}

// CreateDummyUser just creates a user and posts it to the db.
func CreateDummyUser(usr models.User) {
	models.PostUser(usr)
}

// DeleteDummyUser deletes the dummy user created.
func DeleteDummyUser(username string) {
	models.DeleteUser(username)
}

func DeleteMockAdmin() {
	models.DelSessionByUsername("skylarTest")
	models.DeleteUser("skylarTest")
}

func DeleteMockUser() {
	models.DelSessionByUsername("testuser")
	models.DeleteUser("testuser")
}

// TestShowAllUsers tests /users route.
func TestShowAllUsers(t *testing.T) {
	resp := CreateMockAdmin(t)
	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	req.AddCookie(resp.Cookies()[0])
	ps := httprouter.Params{}

	ShowAllUsers(w, req, ps)

	resp = w.Result()

	if resp.StatusCode != 200 {
		t.Error("Expected: 200, Got:", resp.StatusCode)
	}
	DeleteMockAdmin()
}

// TestUpdateUserType tests /updatetype route.
func TestUpdateUserType(t *testing.T) {
	resp := CreateMockAdmin(t)
	_ = CreateMockUser(t)
	form := url.Values{}
	form.Add("username", "testuser")
	form.Add("usertype", "streamer")
	req := httptest.NewRequest("POST", "/updatetype", strings.NewReader(form.Encode()))
	req.Form = form
	cookie := resp.Cookies()[0]
	w := httptest.NewRecorder()
	req.AddCookie(cookie)
	ps := httprouter.Params{}

	UpdateUserType(w, req, ps)

	resp = w.Result()

	if resp.StatusCode != 200 {
		t.Error("Updating testuser to Streamer| Expected: 200, Got:", resp.StatusCode)
	}

	form = url.Values{}
	form.Add("username", "testuser")
	form.Add("usertype", "user")
	req = httptest.NewRequest("POST", "/updatetype", strings.NewReader(form.Encode()))
	req.Form = form

	w = httptest.NewRecorder()
	req.AddCookie(cookie)
	ps = httprouter.Params{}

	UpdateUserType(w, req, ps)

	resp = w.Result()

	if resp.StatusCode != 200 {
		t.Error("Updating testuser back to User | Expected: 200, Got:", resp.StatusCode)
	}

	DeleteMockAdmin()
	DeleteMockUser()
}
