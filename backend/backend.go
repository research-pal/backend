package backend

import (
	"encoding/json"
	//"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
)

func init() {

	r := mux.NewRouter()
	//r.HandleFunc("/", HandleRoot).Methods("GET")
	r.HandleFunc("/notes/{encodedurl}", HandleNotesGet).Methods("GET")
	r.HandleFunc("/notes", HandleNotesPost).Methods("POST")
	r.HandleFunc("/notes", HandleNotesPut).Methods("PUT")

	http.Handle("/", r)

}

func HandleNotesGet(w http.ResponseWriter, r *http.Request) {

	//w.Write([]byte("in Notes Post"))

	c := appengine.NewContext(r)
	notes := Notes{}

	// read encoded url from uri
	params := mux.Vars(r)
	encodedURL, exists := params["encodedurl"]
	if !exists {
		http.Error(w, "url parameter is missing in URI", http.StatusBadRequest)
		return
	}

	// search in database
	if err := notes.get(encodedURL, c); err != nil && err != ErrorNoMatch {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// send response
	if err := json.NewEncoder(w).Encode(notes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func HandleNotesPut(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	notes := Notes{}

	// read from request body
	if err := json.NewDecoder(r.Body).Decode(&notes); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// save to database
	if err := notes.put(c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// send response
	if err := json.NewEncoder(w).Encode(notes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func HandleNotesPost(w http.ResponseWriter, r *http.Request) {

}

// TODO:
//// get/post api endpoint: if url (passed in request body) is already available in database, return that record, otherwise create a record for that url and return that record
//// put api endpoint: update/overwrite the record in the database for the given url
//// struct fields: url, notes (plain text for now, but will be html soon)

// Challanges:
