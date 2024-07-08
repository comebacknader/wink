package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/comebacknader/wink/config"
	"github.com/comebacknader/wink/handlers"
	"github.com/comebacknader/wink/models"
	"github.com/julienschmidt/httprouter"
)

var tpl *template.Template

func init() {
	tpl = config.Tpl
}

type Server struct {
	r *httprouter.Router
}

func main() {
	// Connect to the postgreSQL database
	config.NewDB("postgres://" + os.Getenv("WINK_DB_U") + ":" + os.Getenv("WINK_DB_P") +
		"@" + os.Getenv("WINK_DB_HOST") + "/" + os.Getenv("WINK_DB_NAME") + "?sslmode=disable")

	// config.NewDB("user=" + os.Getenv("WINK_DB_U") +
	// 	" host=172.18.0.2" + " dbname=" + os.Getenv("WINK_DB_NAME") + " sslmode=disable")

	mux := httprouter.New()

	mux.GET("/", index)

	// Auth Handlers
	mux.GET("/signup", handlers.GetSignup)
	mux.POST("/signup", handlers.PostSignup)
	mux.GET("/login", handlers.GetLogin)
	mux.POST("/login", handlers.PostLogin)
	mux.POST("/logout", handlers.PostLogout)
	mux.GET("/forgot", handlers.ForgotPass)
	mux.POST("/forgot", handlers.PostForgotPass)
	mux.GET("/reset", handlers.GetResetPass)
	mux.POST("/reset", handlers.PostResetPass)

	// User Handlers
	mux.GET("/users", handlers.ShowAllUsers)
	mux.POST("/updatetype", handlers.UpdateUserType)
	mux.GET("/sesh", handlers.DelOldSessions)

	// Stream Handlers
	mux.GET("/dashboard/:username", handlers.GetDashboard)
	mux.POST("/updatestream", handlers.UpdateStream)
	mux.POST("/updateonline", handlers.UpdateOnline)
	mux.GET("/streaming/:username", handlers.GetStreaming)

	// Coin Handlers
	mux.GET("/coins", handlers.GetBuyCoins)
	mux.POST("/tip", handlers.PostTip)
	mux.GET("/purchases", handlers.GetPurchases)
	mux.POST("/purch", handlers.PostPurchase)

	// Frontpage Streamer Handlers
	mux.GET("/frontpage", handlers.GetFrontpage)
	mux.POST("/frontpage", handlers.PostFrontpage)

	// Websocket Routes
	go handlers.Hoob.Start()
	mux.GET("/ws", handlers.WebSockPage)

	// Serves the css files called by HTML files
	mux.ServeFiles("/assets/css/*filepath", http.Dir(config.Path+"assets/css/"))

	// Serves the javascript files called by HTML files
	mux.ServeFiles("/assets/js/*filepath", http.Dir(config.Path+"assets/js/"))

	// Serves the images called by HTML files
	mux.ServeFiles("/assets/img/*filepath", http.Dir(config.Path+"assets/img/"))

	mux.GET("/favicon.ico", Favicon)

	// Redirects 404 File Not Found errors to '/' route
	mux.NotFound = http.RedirectHandler("/", 301)

	env := os.Getenv("WINK_ENVIRON")
	if env == "production" {
		go http.ListenAndServe(":80", http.HandlerFunc(redirect))
		log.Fatal(http.ListenAndServeTLS(config.Port,
			config.Path+"www_wink_gg.pem", config.Path+"www_wink_gg.key", &Server{mux}))
	} else {
		log.Fatal(http.ListenAndServe(":8080", &Server{mux}))
		// log.Fatal(http.ListenAndServeTLS(config.Port, "cert.pem", "key.pem", &Server{mux}))
	}
}

// Sets up CORS for all requests
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials",
			"true")
	}
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		errors.New("Request method is OPTIONS")
	}
	s.r.ServeHTTP(w, r)
}

// redirect redirects all http requests to https requests
func redirect(w http.ResponseWriter, req *http.Request) {
	// For non-standard ports :8080 and :10433
	host := strings.Split(req.Host, ":")[0]
	port := ""
	if config.Port != ":443" {
		port = ":10443"
	}

	target := "https://" + host + port + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
	return
}

// Serves the index.gohtml file
func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get currentUser, which does the same as alreadyLoggedIn but returns user
	// with currentUser I could get the stream if I wanted to
	usr, loggedIn := handlers.GetCurrentUser(r)
	streamerData := handlers.StreamerData{}
	streamerData.UserCred.IsLogIn = false
	fpager, pagerErr := models.GetFrontpager()
	if pagerErr != nil {
		panic(pagerErr)
	}
	if loggedIn == true {
		streamerData.UserCred.IsLogIn = true
		if usr.UserType == "streamer" {
			streamerData.CurrUser = usr.Username
			streamerData.UserCred.IsStreamer = true
		}
		if usr.UserType == "admin" {
			streamerData.UserCred.IsAdmin = true
		}
	}
	//Get the stream associated with the user ID
	frontStream, strmErr := models.GetStreamByUID(fpager.ID)
	if strmErr != true || frontStream.Online == false {
		err := tpl.ExecuteTemplate(w, "index.gohtml", streamerData)
		config.HandleError(w, err)
		return
	}
	_, exists := models.GetUserById(fpager.ID)
	if exists == false {
		err := tpl.ExecuteTemplate(w, "home.gohtml", streamerData)
		config.HandleError(w, err)
		return
	}
	err := tpl.ExecuteTemplate(w, "home.gohtml", streamerData)
	config.HandleError(w, err)
	return
}

// Serves the index.gohtml file
func NotFounder(w http.ResponseWriter, req *http.Request) {
	err := tpl.ExecuteTemplate(w, "index.gohtml", nil)
	config.HandleError(w, err)
}

// Favicon : Serves the favicon.ico to the browser
func Favicon(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	http.ServeFile(w, req, config.Path+"assets/img/favicon.ico")
}
