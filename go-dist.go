package main

import (
	"bytes"
	"encoding/json"
	"flag"
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

var cache *string

func ProjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		GitHubWebHookHandler(w, r)
	} else {
		FetchReadMeHandler(w, r)
	}
}

type GithubWebResponse struct {
	Ref string `json:"ref"`
}

func GitHubWebHookHandler(w http.ResponseWriter, r *http.Request) {
	/* extract our vars out */
	vars := mux.Vars(r)

	wh := &GithubWebResponse{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	err := json.Unmarshal(body, wh)

	if err != nil {
		fmt.Println(err.Error())
	}

	if strings.HasSuffix(wh.Ref, "/master") {
		/* used for our project name */
		project := "github.com/" + vars["username"] + "/" + vars["project"]
		directoryName := "/tmp/" + project
		os.RemoveAll(directoryName)
		fmt.Println("Resetting: ", directoryName)
	}
}

func FetchReadMeHandler(w http.ResponseWriter, r *http.Request) {
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

func FetchBinaryHandler(w http.ResponseWriter, r *http.Request) {
	cache := "90s"
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

	dur, dur_err := time.ParseDuration(cache)
	if dur_err != nil {
		dur, _ = time.ParseDuration("60m")
	}

	if err != nil || (err == nil && time.Now().After(info.ModTime().Add(dur))) {
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
	port := flag.Int("port", 80, "Port")
	username := flag.String("username", "[A-Za-z0-9\\-\\_]+", "")
	cache = flag.String("cache", "60m", "Cache for how many minutes")

	flag.Parse()

	fmt.Println("go-dist")
	r := mux.NewRouter()

	/* Grab text for use in your README.md */
	r.HandleFunc("/{username:"+*username+"}/{project:[A-Za-z0-9\\-\\_]+}", ProjectHandler)

	// Routes consist of a path and a handler function.
	r.HandleFunc("/{username:"+*username+"}/{project:[A-Za-z0-9\\-\\_]+}/{os:mac|windows|linux}/{arch:amd64|arm|386}", FetchBinaryHandler)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), r))
}
