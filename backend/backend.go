package main

import (
	"time"

	"github.com/gorilla/mux"
)

type Notes struct {
	URL        string
	Notes      string
	LastUpdate time.Time
}

func init() {

	r := mux.NewRouter()
	r.HandleFunc("/notes", HandleNotesPost).Methods("POST")
	r.HandleFunc("/notes", HandleNotesPut).Methods("PUT")
}

func HandleNotesPost(w http.ResponseWriter, r *http.Request) {}

func HandleNotesPut(w http.ResponseWriter, r *http.Request) {}

// TODO:
//// get/post api endpoint: if url (passed in request body) is already available in database, return that record, otherwise create a record for that url and return that record
//// put api endpoint: update/overwrite the record in the database for the given url
//// struct fields: url, notes (plain text for now, but will be html soon)

// Challanges:
