package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/dfroese-korewireless/continuous-demo/appinfo"
	"github.com/dfroese-korewireless/continuous-demo/sysinfo"
)

// Info contains all the information to be put into the web page
type Info struct {
	AppVersion, IPAddr string
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	sinfo, err := sysinfo.GetSystemInfo()

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ainfo, err := appinfo.GetAppInfo()

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	info := Info{AppVersion: ainfo.AppVersion, IPAddr: sinfo.IPAddress}

	t := template.New("index.html")
	t, err = t.ParseFiles("html/index.html")
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	t.Execute(os.Stdout, info)
	t.Execute(w, info)
}

func main() {
	m := http.NewServeMux()

	m.HandleFunc("/", serveIndex)

	log.Fatal(http.ListenAndServe(":80", m))
}
