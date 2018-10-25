package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/dfroese-korewireless/continuous-demo/api"
	"github.com/dfroese-korewireless/continuous-demo/appinfo"
	"github.com/dfroese-korewireless/continuous-demo/storage"
	"github.com/dfroese-korewireless/continuous-demo/sysinfo"
	"github.com/gorilla/mux"
)

// Info contains all the information to be put into the web page
type Info struct {
	AppVersion, IPAddr, ContainerName string
}

var appInfo *Info

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if appInfo == nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	t := template.New("index.html")
	t, err := t.ParseFiles("html/index.html")
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	t.Execute(w, appInfo)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	return
}

func healthCheckJSON(w http.ResponseWriter, r *http.Request) {
	return
}

func main() {
	sinfo := sysinfo.GetSystemInfo()
	ainfo, err := appinfo.GetAppInfo()
	if err != nil {
		panic(err)
	}
	appInfo = &Info{AppVersion: ainfo.AppVersion, IPAddr: sinfo.IPAddress, ContainerName: sinfo.ContainerName}

	db, err := storage.New(ainfo.DBPath)
	if err != nil {
		panic(err)
	}
	apiCtx := api.New(db)

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	router.HandleFunc("/", serveIndex)
	router.HandleFunc("/health", healthCheck)
	apiRouter.HandleFunc("/health", healthCheckJSON).Methods("GET")
	apiRouter.HandleFunc("/messages", apiCtx.CreateMessage).Methods("POST")
	apiRouter.HandleFunc("/messages", apiCtx.GetMessages).Methods("GET")
	apiRouter.HandleFunc("/messages/{id}", apiCtx.GetMessage).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+ainfo.Port, router))
}
