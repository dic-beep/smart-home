package main

import (
	"Smart-Home/handler"
	"log"
	"net/http"
)


func main() {
	// Initialize the Router
	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("/dashboard", handler.Dashboard)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	port := ":8080"
	log.Println("Starting at port 8080")
	err := http.ListenAndServe(port, nil)
	log.Fatal(err)
}
