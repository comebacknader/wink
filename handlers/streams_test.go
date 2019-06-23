package handlers

import (
	"github.com/comebacknader/wink/models"
	"bytes"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// CreateMockStreamer creates and logs in  mock streamer.
func CreateMockStreamer(t *testing.T) *http.Response {
	// Create a Mock Admin to do the updating
	resp := CreateMockAdmin(t)
	// Create 'testuser' to have in db
	userResp := CreateMockUser(t)
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

	DeleteMockAdmin()

	if resp.StatusCode != 200 {
		t.Error("Creating Mock Streamer| Expected: 200, Got:", resp.StatusCode)
	}

	return userResp
}

// DeleteMockStreamer deletes the mock streamer.
func DeleteMockStreamer() {
	//Need to also delete stream
	//Need to find the id of the mock streamer
	usr, _ := models.GetUserByName("testuser")
	models.DeleteStreamByUID(usr.ID)
	DeleteMockUser()
}

// TestGetDashboard tests the /dashboard/:username route.
func TestGetDashboard(t *testing.T) {
	resp := CreateMockStreamer(t)
	req := httptest.NewRequest("GET", "/dashboard/testuser", nil)
	cookie := resp.Cookies()[0]
	w := httptest.NewRecorder()
	req.AddCookie(cookie)
	param := httprouter.Param{}
	param.Key = "username"
	param.Value = "testuser"
	ps := httprouter.Params{param}

	GetDashboard(w, req, ps)

	resp = w.Result()

	if resp.StatusCode != 200 {
		t.Error("Get Dashboard | Expected: 200, Got:", resp.StatusCode)
	}

	DeleteMockStreamer()
}

// TestUpdateStream tests the /updatestream route.
func TestUpdateStream(t *testing.T) {
	resp := CreateMockStreamer(t)

	// Test Valid Submission
	form := url.Values{}
	form.Add("title", "New Stream Title")
	form.Add("game", "CS:GO")
	form.Add("twitter", "www.twitter.com/testUser")
	form.Add("siteone", "www.mytestusersite.com")
	req := httptest.NewRequest("POST", "/updatestream", strings.NewReader(form.Encode()))
	req.Form = form
	cookie := resp.Cookies()[0]
	w := httptest.NewRecorder()
	req.AddCookie(cookie)
	ps := httprouter.Params{}

	UpdateStream(w, req, ps)

	resp = w.Result()

	if resp.StatusCode != 200 {
		t.Error("Updating Stream Valid| Expected: 200, Got:", resp.StatusCode)
	}

	// Test Blank Title
	form = url.Values{}
	form.Add("title", "")
	form.Add("game", "CS:GO")
	form.Add("twitter", "www.twitter.com/testUser")
	form.Add("siteone", "www.mytestusersite.com")
	req = httptest.NewRequest("POST", "/updatestream", strings.NewReader(form.Encode()))
	req.Form = form
	w = httptest.NewRecorder()
	req.AddCookie(cookie)

	UpdateStream(w, req, ps)

	resp = w.Result()

	if resp.StatusCode != 400 {
		t.Error("Updating Stream Blank Title| Expected: 400, Got:", resp.StatusCode)
	}

	// Test Blank Game Title
	form = url.Values{}
	form.Add("title", "New Stream Title")
	form.Add("game", "")
	form.Add("twitter", "www.twitter.com/testUser")
	form.Add("siteone", "www.mytestusersite.com")
	req = httptest.NewRequest("POST", "/updatestream", strings.NewReader(form.Encode()))
	req.Form = form
	w = httptest.NewRecorder()
	req.AddCookie(cookie)

	UpdateStream(w, req, ps)

	resp = w.Result()

	if resp.StatusCode != 400 {
		t.Error("Updating Stream Blank Game| Expected: 400, Got:", resp.StatusCode)
	}

	// Test Title Too Long
	form = url.Values{}
	form.Add("title", `WhetherWhetherWhetherWhetherWhetherWhether
		WhetherWhetherWhetherWhetherWhetherWhetherWhetherWhetherWhether`)
	form.Add("game", "CS:GO")
	form.Add("twitter", "www.twitter.com/testUser")
	form.Add("siteone", "www.mytestusersite.com")
	req = httptest.NewRequest("POST", "/updatestream", strings.NewReader(form.Encode()))
	req.Form = form
	w = httptest.NewRecorder()
	req.AddCookie(cookie)

	UpdateStream(w, req, ps)

	resp = w.Result()

	if resp.StatusCode != 400 {
		t.Error("Updating Stream Long Title| Expected: 400, Got:", resp.StatusCode)
	}

	// Test Title Too Long
	form = url.Values{}
	form.Add("title", `WhetherWhetherWhetherWhetherWhetherWhether
		WhetherWhetherWhetherWhetherWhetherWhetherWhetherWhetherWhether`)
	form.Add("game", "CS:GO")
	form.Add("twitter", "www.twitter.com/testUser")
	form.Add("siteone", "www.mytestusersite.com")
	req = httptest.NewRequest("POST", "/updatestream", strings.NewReader(form.Encode()))
	req.Form = form
	w = httptest.NewRecorder()
	req.AddCookie(cookie)

	UpdateStream(w, req, ps)

	resp = w.Result()

	if resp.StatusCode != 400 {
		t.Error("Updating Stream Long Title| Expected: 400, Got:", resp.StatusCode)
	}

	DeleteMockStreamer()
}

type OnlineStatus struct {
	Online string
}

// TestUpdateOnline tests whether the streamer can update online status.
func TestUpdateOnline(t *testing.T) {
	resp := CreateMockStreamer(t)

	// Test status goes online
	form := url.Values{}
	form.Add("online", "online")
	req := httptest.NewRequest("POST", "/updateonline", strings.NewReader(form.Encode()))
	req.Form = form
	cookie := resp.Cookies()[0]
	w := httptest.NewRecorder()
	req.AddCookie(cookie)
	ps := httprouter.Params{}

	UpdateOnline(w, req, ps)

	resp = w.Result()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respStr := buf.String()

	if respStr != `{"status":"online"}` {
		t.Error(`Updating Stream Status| Expected {"status":"online"}`, respStr)
	}

	// Test status goes offline
	form = url.Values{}
	form.Add("online", "offline")
	req = httptest.NewRequest("POST", "/updateonline", strings.NewReader(form.Encode()))
	req.Form = form
	w = httptest.NewRecorder()
	req.AddCookie(cookie)

	UpdateOnline(w, req, ps)

	resp = w.Result()

	buf = new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respStr = buf.String()

	if respStr != `{"status":"offline"}` {
		t.Error(`Updating Stream Status| Expected {"status":"offline"}`, respStr)
	}

	DeleteMockStreamer()
}

// TestGetStreaming tests the /streaming/:username route.
func TestGetStreaming(t *testing.T) {
	resp := CreateMockStreamer(t)
	req := httptest.NewRequest("GET", "/streaming/testuser", nil)
	cookie := resp.Cookies()[0]
	w := httptest.NewRecorder()
	req.AddCookie(cookie)
	param := httprouter.Param{}
	param.Key = "username"
	param.Value = "testuser"
	ps := httprouter.Params{param}

	GetStreaming(w, req, ps)

	resp = w.Result()

	if resp.StatusCode != 200 {
		t.Error("Get Streaming | Expected: 200, Got:", resp.StatusCode)
	}

	DeleteMockStreamer()
}

// TestGetFrontpage tests the GET /frontpage route.
func TestGetFrontpage(t *testing.T) {
	resp := CreateMockAdmin(t)
	req := httptest.NewRequest("GET", "/frontpage", nil)
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

// TestPostFrontpage tests the POST /frontpage route.
func TestPostFrontpage(t *testing.T) {
	CreateMockStreamer(t)
	resp := CreateMockAdmin(t)

	// Testing adding streamer to frontpage.
	form := url.Values{}
	form.Add("streamer", "testuser")
	form.Add("addremove", "add")
	req := httptest.NewRequest("POST", "/frontpage", strings.NewReader(form.Encode()))
	req.Form = form
	cookie := resp.Cookies()[0]
	w := httptest.NewRecorder()
	req.AddCookie(cookie)
	ps := httprouter.Params{}

	PostFrontpage(w, req, ps)

	resp = w.Result()

	if resp.StatusCode != 200 {
		t.Error("PostFrontpage Add Streamer | Expected: 200, Got: ", resp.StatusCode)
	}

	// Testing adding streamer to frontpage.
	form = url.Values{}
	form.Add("streamer", "testuser")
	form.Add("addremove", "remove")
	req = httptest.NewRequest("POST", "/frontpage", strings.NewReader(form.Encode()))
	req.Form = form
	w = httptest.NewRecorder()
	req.AddCookie(cookie)

	PostFrontpage(w, req, ps)

	resp = w.Result()

	if resp.StatusCode != 200 {
		t.Error("PostFrontpage Remove Streamer | Expected: 200, Got: ", resp.StatusCode)
	}

	DeleteMockStreamer()
	DeleteMockAdmin()
}
