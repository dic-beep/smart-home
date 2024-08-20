package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
	"text/template"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
	//Setup default page
	if r.URL.Path != "/dashboard" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "POST" {
		//load HTML file
		tmpl, err := template.ParseFiles(path.Join("views", "dashboard.html"), path.Join("views", "setup.html"))
		if err != nil {
			log.Println(err)
			http.Error(w, "Something wrong, try again later", http.StatusInternalServerError)
			return
		}

		//execute HTML file
		tmpl.Execute(w, nil)
		if err != nil {
			log.Println(err)
			http.Error(w, "Something wrong, try again later", http.StatusInternalServerError)
			return
		}
		state := r.URL.Query().Get("state")
		var ledState bool
		if state == "on" {
			ledState = true
		} else if state == "off" {
			ledState = false
		} else {
			http.Error(w, "Invalid state", http.StatusBadRequest)
			return
		}

		// Update Firebase Realtime Database
		url := "https://smart-home-bdf74-default-rtdb.firebaseio.com/LED/digital.json"
		req, err := http.NewRequest("PUT", url, strings.NewReader(fmt.Sprintf("%v", ledState)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "LED state updated: %s", body)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//load HTML file
		tmpl, err := template.ParseFiles(path.Join("views", "login.html"), path.Join("views", "setup.html"))
		if err != nil {
			log.Println(err)
			http.Error(w, "something wrong, try again later", http.StatusInternalServerError)
			return
		}

		//execute the file
		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println(err)
			http.Error(w, "something wrong, try again later", http.StatusInternalServerError)
			return
		}
		return
	}
	http.Error(w, "HTTP methode must GET", http.StatusBadRequest)
}
