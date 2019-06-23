package handlers

import (
	"github.com/comebacknader/wink/models"
	"github.com/julienschmidt/httprouter"
	_ "net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestGetBuyCoins tests the GET request to /coins.
func TestGetBuyCoins(t *testing.T) {
	resp := CreateMockUser(t)
	req := httptest.NewRequest("GET", "/coins", nil)
	cookie := resp.Cookies()[0]
	w := httptest.NewRecorder()
	req.AddCookie(cookie)
	ps := httprouter.Params{}

	GetBuyCoins(w, req, ps)

	resp = w.Result()

	if resp.StatusCode != 200 {
		t.Error("Get Buy Coins | Expected: 200, Got:", resp.StatusCode)
	}

	DeleteMockUser()
}

// TestPostTip tests sending a streamer coins.
func TestPostTip(t *testing.T) {
	resp := CreateMockUser(t)
	// Give the MockUser 50 coins
	models.AddCoinsToUser("testuser", 50)

	usr := models.User{}
	usr.Username = "dummyuser"
	usr.UserType = "streamer"
	usr.Hash = "somebullshithash"
	usr.Email = "dummy@dumb.com"

	CreateDummyUser(usr)

	// Test user doesn't have enough coins to give
	form := url.Values{}
	form.Add("sender", "testuser")
	form.Add("amount", "100")
	form.Add("receiver", "dummyuser")
	req := httptest.NewRequest("POST", "/tip", strings.NewReader(form.Encode()))
	req.Form = form
	cookie := resp.Cookies()[0]
	w := httptest.NewRecorder()
	req.AddCookie(cookie)
	ps := httprouter.Params{}

	PostTip(w, req, ps)

	resp = w.Result()

	if resp.StatusCode != 400 {
		t.Error("PostTip Not Enough Coins| Expected: 400, Got:", resp.StatusCode)
	}

	// Test seccessful tipping
	form = url.Values{}
	form.Add("sender", "testuser")
	form.Add("amount", "25")
	form.Add("receiver", "dummyuser")
	req = httptest.NewRequest("POST", "/tip", strings.NewReader(form.Encode()))
	req.Form = form
	w = httptest.NewRecorder()
	req.AddCookie(cookie)

	PostTip(w, req, ps)

	resp = w.Result()

	if resp.StatusCode != 200 {
		t.Error("PostTip Success | Expected 200, Got:", resp.StatusCode)
	}

	DeleteDummyUser(usr.Username)
	DeleteMockUser()
}

// TestPostPurchase tests the buying of coins.
