package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Hello just a struct
type Hello struct {
	Hello string
}

// Version a struct to show the version
type Version struct {
	Version string
}

// our main function
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/hello", GetHello).Methods("GET")
	router.HandleFunc("/version", GetVersion).Methods("GET")
	fmt.Printf("%s\n", "Starting server on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

// GetHello return world
func GetHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s\n", "GetHello called")

	// create hello object
	hello := Hello{"world"}

	// to json
	js, err := json.Marshal(hello)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// GetVersion return the current version of the app
func GetVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s\n", "GetHello called")

	// create hello object
	version := Version{"1.0"}

	// to json
	js, err := json.Marshal(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
