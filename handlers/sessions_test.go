package handlers

import (
	"github.com/comebacknader/wink/config"
	"github.com/comebacknader/wink/models"
	"github.com/julienschmidt/httprouter"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

func init() {
	// Connect to the postgreSQL database
	config.NewDB("postgres://" + os.Getenv("WINK_DB_U") + ":" + os.Getenv("WINK_DB_P") +
		"@" + os.Getenv("WINK_DB_HOST") + "/" + os.Getenv("WINK_DB_NAME") + "?sslmode=disable")
}

// TestPostUser tests the Signup functionality.
func TestPostUser(t *testing.T) {
	const reqUrl = "/signup"

	// Test : Valid submission.
	validUser := url.Values{}
	validUser.Add("username", "skylarTest")
	validUser.Add("email", "skylar@test.com")
	validUser.Add("password", "tester")

	req := httptest.NewRequest("POST", reqUrl, strings.NewReader(validUser.Encode()))
	req.Form = validUser
	w := httptest.NewRecorder()
	ps := httprouter.Params{}

	PostSignup(w, req, ps)
	resp := w.Result()
	if resp.StatusCode != 303 {
		t.Error("Valid Submission | Expected Status Code: 303, Got: ", resp.StatusCode)
	}

	// Failing user submissions.

	// Test : If username is empty.
	emptyUsername := url.Values{}
	emptyUsername.Add("username", "")
	emptyUsername.Add("email", "test@invalid.com")
	emptyUsername.Add("password", "tester")

	req = httptest.NewRequest("POST", reqUrl, strings.NewReader(emptyUsername.Encode()))
	req.Form = emptyUsername
	w = httptest.NewRecorder()
	ps = httprouter.Params{}

	PostSignup(w, req, ps)
	resp = w.Result()
	if resp.StatusCode != 400 {
		t.Error("Empty Username | Expected Status Code: 400, Got: ", resp.StatusCode)
	}

	// Test : If username is too short.
	shortUsername := url.Values{}
	shortUsername.Add("username", "sky")
	shortUsername.Add("email", "test@invalid.com")
	shortUsername.Add("password", "tester")

	req = httptest.NewRequest("POST", reqUrl, strings.NewReader(shortUsername.Encode()))
	req.Form = shortUsername
	w = httptest.NewRecorder()

	PostSignup(w, req, ps)
	resp = w.Result()
	if resp.StatusCode != 400 {
		t.Error("Short Username | Expected Status Code: 400, Got: ", resp.StatusCode)
	}

	// Test : If username is too long.
	longUsername := url.Values{}
	longUsername.Add("username", "thisIsAReallyFuckingLongUsernameTooFuckingLongBro")
	longUsername.Add("email", "test@invalid.com")
	longUsername.Add("password", "tester")

	req = httptest.NewRequest("POST", reqUrl, strings.NewReader(longUsername.Encode()))
	req.Form = longUsername
	w = httptest.NewRecorder()

	PostSignup(w, req, ps)
	resp = w.Result()
	if resp.StatusCode != 400 {
		t.Error("Long Username | Expected Status Code: 400, Got: ", resp.StatusCode)
	}

	// Test : If username is already taken.
	takenUsername := url.Values{}
	takenUsername.Add("username", "skylarTest")
	takenUsername.Add("email", "test@invalid.com")
	takenUsername.Add("password", "tester")

	req = httptest.NewRequest("POST", reqUrl, strings.NewReader(takenUsername.Encode()))
	req.Form = takenUsername
	w = httptest.NewRecorder()

	PostSignup(w, req, ps)
	resp = w.Result()
	if resp.StatusCode != 400 {
		t.Error("Taken Username | Expected Status Code: 400, Got: ", resp.StatusCode)
	}

	// Test : If email is empty.
	emptyEmail := url.Values{}
	emptyEmail.Add("username", "skylar101")
	emptyEmail.Add("email", "")
	emptyEmail.Add("password", "tester")

	req = httptest.NewRequest("POST", reqUrl, strings.NewReader(emptyEmail.Encode()))
	req.Form = emptyEmail
	w = httptest.NewRecorder()

	PostSignup(w, req, ps)
	resp = w.Result()
	if resp.StatusCode != 400 {
		t.Error("EmptyEmail | Expected Status Code: 400, Got: ", resp.StatusCode)
	}

	// Test : If email is too short.
	shortEmail := url.Values{}
	shortEmail.Add("username", "skylar101")
	shortEmail.Add("email", "s@s.com")
	shortEmail.Add("password", "tester")

	req = httptest.NewRequest("POST", reqUrl, strings.NewReader(shortEmail.Encode()))
	req.Form = shortEmail
	w = httptest.NewRecorder()

	PostSignup(w, req, ps)
	resp = w.Result()
	if resp.StatusCode != 400 {
		t.Error("Short Email | Expected Status Code: 400, Got: ", resp.StatusCode)
	}

	// Test : If email is too long.
	longEmail := url.Values{}
	longEmail.Add("username", "skylar101")
	longEmail.Add("email", "thisshitbewaytoolonghomiewaytoolong@toolongbrahtoolong.com")
	longEmail.Add("password", "tester")

	req = httptest.NewRequest("POST", reqUrl, strings.NewReader(longEmail.Encode()))
	req.Form = longEmail
	w = httptest.NewRecorder()

	PostSignup(w, req, ps)
	resp = w.Result()
	if resp.StatusCode != 400 {
		t.Error("Long Email | Expected Status Code: 400, Got: ", resp.StatusCode)
	}

	// Test : If email is incorrect format.
	badFmtEmail := url.Values{}
	badFmtEmail.Add("username", "skylar101")
	badFmtEmail.Add("email", "skylar.com")
	badFmtEmail.Add("password", "tester")

	req = httptest.NewRequest("POST", reqUrl, strings.NewReader(badFmtEmail.Encode()))
	req.Form = badFmtEmail
	w = httptest.NewRecorder()

	PostSignup(w, req, ps)
	resp = w.Result()
	if resp.StatusCode != 400 {
		t.Error("Bad Format Email | Expected Status Code: 400, Got: ", resp.StatusCode)
	}

	// Test : If email already exists.
	takenEmail := url.Values{}
	takenEmail.Add("username", "skylar101")
	takenEmail.Add("email", "skylar@test.com")
	takenEmail.Add("password", "tester")

	req = httptest.NewRequest("POST", reqUrl, strings.NewReader(takenEmail.Encode()))
	req.Form = takenEmail
	w = httptest.NewRecorder()

	PostSignup(w, req, ps)
	resp = w.Result()
	if resp.StatusCode != 400 {
		t.Error("Taken Email | Expected Status Code: 400, Got: ", resp.StatusCode)
	}

	// Test : Password is empty.
	emptyPassword := url.Values{}
	emptyPassword.Add("username", "skylar101")
	emptyPassword.Add("email", "testShouldFail@test.com")
	emptyPassword.Add("password", "")

	req = httptest.NewRequest("POST", reqUrl, strings.NewReader(emptyPassword.Encode()))
	req.Form = emptyPassword
	w = httptest.NewRecorder()

	PostSignup(w, req, ps)
	resp = w.Result()
	if resp.StatusCode != 400 {
		t.Error("EmptyPassword | Expected Status Code: 400, Got: ", resp.StatusCode)
	}

	models.DelSessionByUsername("skylarTest")
	models.DeleteUser("skylarTest")
}

// TestPostLogin tests the Login functionality.
func TestPostLogin(t *testing.T) {
	const reqUrl = "/login"

	// Test : Valid submission.
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
		t.Error("Populating User | Expected Status Code: 303, Got: ", resp.StatusCode)
	}

	// Test : Valid login attempt with username.
	validUser = url.Values{}
	validUser.Add("credential", "skylarTest")
	validUser.Add("password", "tester")

	req = httptest.NewRequest("POST", reqUrl, strings.NewReader(validUser.Encode()))
	req.Form = validUser
	w = httptest.NewRecorder()
	ps = httprouter.Params{}

	PostLogin(w, req, ps)
	resp = w.Result()
	if resp.StatusCode != 303 {
		t.Error("Valid Login w/ Username | Expected Status Code: 303, Got: ", resp.StatusCode)
	}
	models.DelSessionByUsername("skylarTest")

	// Test : Valid login attempt with email.
	validUser = url.Values{}
	validUser.Add("credential", "skylar@test.com")
	validUser.Add("password", "tester")

	req = httptest.NewRequest("POST", reqUrl, strings.NewReader(validUser.Encode()))
	req.Form = validUser
	w = httptest.NewRecorder()
	ps = httprouter.Params{}

	PostLogin(w, req, ps)
	resp = w.Result()
	if resp.StatusCode != 303 {
		t.Error("Valid Login w/ Email | Expected Status Code: 303, Got: ", resp.StatusCode)
	}
	models.DelSessionByUsername("skylarTest")

	// Test : Empty credential.
	validUser = url.Values{}
	validUser.Add("credential", "")
	validUser.Add("password", "tester")

	req = httptest.NewRequest("POST", reqUrl, strings.NewReader(validUser.Encode()))
	req.Form = validUser
	w = httptest.NewRecorder()
	ps = httprouter.Params{}

	PostLogin(w, req, ps)
	resp = w.Result()
	if resp.StatusCode != 400 {
		t.Error("Empty Credential | Expected Status Code: 400, Got: ", resp.StatusCode)
	}

	// Test : Empty password.
	validUser = url.Values{}
	validUser.Add("credential", "skylarTest")
	validUser.Add("password", "")

	req = httptest.NewRequest("POST", reqUrl, strings.NewReader(validUser.Encode()))
	req.Form = validUser
	w = httptest.NewRecorder()
	ps = httprouter.Params{}

	PostLogin(w, req, ps)
	resp = w.Result()
	if resp.StatusCode != 400 {
		t.Error("Empty Password | Expected Status Code: 400, Got: ", resp.StatusCode)
	}

	// Test : Invalid credentials.
	validUser = url.Values{}
	validUser.Add("credential", "skylarTest2000000")
	validUser.Add("password", "tester")

	req = httptest.NewRequest("POST", reqUrl, strings.NewReader(validUser.Encode()))
	req.Form = validUser
	w = httptest.NewRecorder()
	ps = httprouter.Params{}

	PostLogin(w, req, ps)
	resp = w.Result()
	if resp.StatusCode != 400 {
		t.Error("Invalid Credential | Expected Status Code: 400, Got: ", resp.StatusCode)
	}

	// Tear Down : Delete username 'skylarTest' from database.
	models.DeleteUser("skylarTest")
}

// TestPostLogout tests the Logout functionality at /logout.
func TestPostLogout(t *testing.T) {
	resp := CreateMockAdmin(t)

	req := httptest.NewRequest("POST", "/logout", nil)
	w := httptest.NewRecorder()
	req.AddCookie(resp.Cookies()[0])
	ps := httprouter.Params{}

	PostLogout(w, req, ps)
	resp = w.Result()
	if resp.StatusCode != 303 {
		t.Error("LogOut User | Expected Status Code: 303, Got: ", resp.StatusCode)
	}

	models.DeleteUser("skylarTest")
}

// TestForgotPass tests submitting an email when user forgets password.
func TestForgotPass(t *testing.T) {
	// Test : Empty email doesn't go through.
	invalidEmail := url.Values{}
	invalidEmail.Add("email", "")

	req := httptest.NewRequest("POST", "/forgot", strings.NewReader(invalidEmail.Encode()))
	req.Form = invalidEmail
	w := httptest.NewRecorder()
	ps := httprouter.Params{}

	PostForgotPass(w, req, ps)
	resp := w.Result()
	if resp.StatusCode != 400 {
		t.Error("Invalid Email Forgot Pass | Expected Status Code: 400, Got: ", resp.StatusCode)
	}

	// Test : Valid email does go through.
	validEmail := url.Values{}
	validEmail.Add("email", "dragonandballz@gmail.com")

	req = httptest.NewRequest("POST", "/forgot", strings.NewReader(validEmail.Encode()))
	req.Form = validEmail
	w = httptest.NewRecorder()
	ps = httprouter.Params{}

	PostForgotPass(w, req, ps)
	resp = w.Result()
	if resp.StatusCode != 200 {
		t.Error("Valid Email Forgot Pass | Expected Status Code: 200, Got: ", resp.StatusCode)
	}
}

// TestResetPass tests resetting a password.
func TestResetPass(t *testing.T) {
	// Get the reset token from the database.
	token, err := models.GetTokenByEmail("dragonandballz@gmail.com")
	if err != nil || token == "" {
		t.Error("There was a problem retrieving ResetPassToken.")
	}

	// Test : The password and password confirmation do not match.
	nonConfPass := url.Values{}
	nonConfPass.Add("password", "tester123")
	nonConfPass.Add("confirmation", "tester124")
	nonConfPass.Add("token", token)
	req := httptest.NewRequest("POST", "/reset", strings.NewReader(nonConfPass.Encode()))
	req.Form = nonConfPass
	w := httptest.NewRecorder()
	ps := httprouter.Params{}

	PostResetPass(w, req, ps)
	resp := w.Result()
	if resp.StatusCode != 400 {
		t.Error("Pass Not Match Conf | Expected Status Code: 400, Got: ", resp.StatusCode)
	}

	// Test : Password too short.
	shortPass := url.Values{}
	shortPass.Add("password", "t")
	shortPass.Add("confirmation", "t")
	shortPass.Add("token", token)
	req = httptest.NewRequest("POST", "/reset", strings.NewReader(shortPass.Encode()))
	req.Form = shortPass
	w = httptest.NewRecorder()

	PostResetPass(w, req, ps)
	resp = w.Result()
	if resp.StatusCode != 400 {
		t.Error("Short Password | Expected Status Code: 400, Got: ", resp.StatusCode)
	}

	// Test : Valid password change.
	validChange := url.Values{}
	validChange.Add("password", "tester")
	validChange.Add("confirmation", "tester")
	validChange.Add("token", token)
	req = httptest.NewRequest("POST", "/reset", strings.NewReader(validChange.Encode()))
	req.Form = validChange
	w = httptest.NewRecorder()

	PostResetPass(w, req, ps)
	resp = w.Result()
	if resp.StatusCode != 200 {
		t.Error("Valid Password Change | Expected Status Code: 200, Got: ", resp.StatusCode)
	}

}
