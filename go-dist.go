package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func FetchReadMe(w http.ResponseWriter, r *http.Request) {
	/* extract our vars out */
	vars := mux.Vars(r)

	prefix := "http://"
	if r.TLS != nil {
		prefix = "https://"
	}

	url := prefix + r.Host

	readme := "## Binaries\n"
	readme += "![Mac OSX](https://raw.githubusercontent.com/kcmerrill/go-dist/master/assets/apple_logo.png \"Mac OSX\") [386](" + url + "/" + vars["username"] + "/" + vars["project"] + "/mac/386) | [amd64](" + url + "/" + vars["username"] + "/" + vars["project"] + "/mac/amd64)"
	readme += "\n\n"
	readme += "![Windows](https://raw.githubusercontent.com/kcmerrill/go-dist/master/assets/windows_logo.png \"Windows\") [386](" + url + "/" + vars["username"] + "/" + vars["project"] + "/windows/386) | [amd64](" + url + "/" + vars["username"] + "/" + vars["project"] + "/windows/amd64)"
	readme += "\n\n"
	readme += "![Linux](https://raw.githubusercontent.com/kcmerrill/go-dist/master/assets/linux_logo.png \"Linux\") [386](" + url + "/" + vars["username"] + "/" + vars["project"] + "/linux/386) | [amd64](" + url + "/" + vars["username"] + "/" + vars["project"] + "/linux/amd64)"
	readme += "\n\n"

	w.Write([]byte(readme))
}

func FetchBinary(w http.ResponseWriter, r *http.Request) {
	/* extract our vars out */
	vars := mux.Vars(r)

	/* mac? really should be darwin ... dumb */
	if vars["os"] == "mac" {
		vars["os"] = "darwin"
	}

	/* used for our project name */
	project := "github.com/" + vars["username"] + "/" + vars["project"]

	/* filename, used to download */
	downloadFileName := vars["project"]
	directoryName := "/tmp/" + project
	actualFileName := directoryName + "/" + vars["project"] + "_" + vars["os"] + "_" + vars["arch"]

	if vars["os"] == "windows" {
		actualFileName = actualFileName + ".exe"
		downloadFileName = downloadFileName + ".exe"
	}

	/* does the folder exist? */
	info, err := os.Stat(actualFileName)

	if err != nil || (err == nil && time.Now().After(info.ModTime().Add(time.Hour*1))) {
		if dir_err := os.MkdirAll(directoryName, 0755); dir_err == nil {
			_, project_err := exec.Command("bash", "-c", "go get -u "+project).Output()
			if project_err == nil {
				_, build_err := exec.Command("bash", "-c", "gox --os="+vars["os"]+" --arch="+vars["arch"]+" --output="+strings.Replace(actualFileName, ".exe", "", -1)+" "+project).Output()
				if build_err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("error"))
					return
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("invalid project"))
				return
			}
		}
	}

	if data, read_err := ioutil.ReadFile(actualFileName); read_err == nil {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename="+downloadFileName)
		w.Header().Set("Content-Transfer-Encoding", "binary")
		w.Header().Set("Expires", "0")
		http.ServeContent(w, r, actualFileName, time.Now(), bytes.NewReader(data))
		return
	} else {
		fmt.Println("actual", actualFileName)
		fmt.Println("download", downloadFileName)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error2"))
		return
	}
}

func main() {
	r := mux.NewRouter()

	/* TODO: Get a proper regex for project */

	/* Grab text for use in your README.md */
	r.HandleFunc("/{username:[A-Za-z0-9\\-\\_]+}/{project:[A-Za-z0-9\\-\\_]+}", FetchReadMe)

	// Routes consist of a path and a handler function.
	r.HandleFunc("/{username:[A-Za-z0-9\\-\\_]+}/{project:[A-Za-z0-9\\-\\_]+}/{os:mac|windows|linux}/{arch:amd64|arm|386}", FetchBinary)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":80", r))
}
