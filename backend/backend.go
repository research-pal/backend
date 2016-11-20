package backend

import (
	"encoding/json"
	//"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
)

func init() {

	r := mux.NewRouter()
	//r.HandleFunc("/", HandleRoot).Methods("GET")
	r.HandleFunc("/notes", HandleNotesGet).Methods("GET")
	r.HandleFunc("/notes", HandleNotesPut).Methods("PUT")
	//r.HandleFunc("/notes/{encodedurl}", HandleNotesDelete).Methods("DELETE")

	http.Handle("/", r)

}

func HandleNotesGet(w http.ResponseWriter, r *http.Request) {

	//w.Write([]byte("in Notes Post"))

	c := appengine.NewContext(r)
	notes := Notes{}

	// read encoded url from uri
	//params := mux.Vars(r)
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	encodedURL := ""
	if val, exists := params["encodedurl"]; exists {
		encodedURL = val[0]
	} else {
		http.Error(w, "url parameter is missing in URI", http.StatusBadRequest)
		return
	}

	// search in database
	if err := notes.get(encodedURL, c); err != nil && err != ErrorNoMatch {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if err == ErrorNoMatch {
		http.Error(w, err.Error()+": "+encodedURL, http.StatusNotFound)
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

/*
func HandleNotesDelete(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	params := mux.Vars(r)

	encodedurl, exists := params["encodedurl"]
	if !exists {
		w.WriteHeader(http.StatusInternalServerError)
		// add error notesmessage that "goal parameter is missing"
		return
	}
	notes := Notes{}
	notes.URL = encodedurl

	err := notes.Delete(c)
	if err == ErrorNoMatch {
		http.Error(w, err.Error(), http.StatusOK)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func HandleNotesGetAll(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	if err := Notes.Get(c, goalFilter, offset, limit); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(goals); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
*/
// TODO:
//// get/post api endpoint: if url (passed in request body) is already available in database, return that record, otherwise create a record for that url and return that record
//// put api endpoint: update/overwrite the record in the database for the given url
//// struct fields: url, notes (plain text for now, but will be html soon)

// Challanges:
