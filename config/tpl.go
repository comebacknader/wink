package config

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

var Tpl *template.Template
var Path string
var Port string

func init() {
	CreateFilePath()
	DeterminePort()
	Tpl = template.Must(template.ParseGlob(Path + "templates/*"))
}

// HandleError :  A generic error handler.
func HandleError(res http.ResponseWriter, err error) {
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}

// CreateFilePath sets the Path to the main.go file depending on env.
func CreateFilePath() {
	wd, _ := os.Getwd()
	Path = os.Getenv("WINK_PATH")
	env := os.Getenv("WINK_ENVIRON")
	if env == "production" {
		Path = "/home/ubuntu/gospace/src/github.com/comebacknader/wink/"
	} else {
		if wd == Path {
			Path = ""
		}
	}
}

// DeterminePort determines port :10433 for dev :433 for prod.
func DeterminePort() {
	env := os.Getenv("WINK_ENVIRON")
	if env == "production" {
		// Port = ":433"
		Port = ":443"
	} else {
		Port = ":10443"
	}
}
