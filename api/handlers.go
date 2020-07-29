package api

import (
	"encoding/json"
	"net/http"
	"net/url"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"github.com/research-pal/backend/db/notes"
	"google.golang.org/appengine"
)

var (
	dbConn     *firestore.Client
	encodedurl string
)

func Init(client *firestore.Client) {
	dbConn = client
}

func HandleNotesGetByID(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	params := mux.Vars(r)
	// TODO: no query parameters..
	fieldValue := ""
	if params["docid"] != "" {
		fieldValue = params["docid"]
	}

	note, err := notes.GetByID(c, dbConn, fieldValue)
	if err != nil && err != notes.ErrorNoMatch {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if err == notes.ErrorNoMatch {
		http.Error(w, err.Error()+": "+fieldValue, http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleNotesGetFiltered(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// TODO: use mux functions to retrieve query params
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Group Query params: externalTaskURL(encodedurl), status, group, assignee

	field, fieldValue := "", ""
	if val, exists := params["encodedurl"]; exists {
		field, fieldValue = "url", val[0]
	} else if val, exists := params["status"]; exists {
		field, fieldValue = "status", val[0]
	} else if val, exists := params["assignee"]; exists {
		field, fieldValue = "assignee", val[0]
	} else if val, exists := params["group"]; exists {
		field, fieldValue = "group", val[0]
	} else {
		http.Error(w, "parameter is missing in URI", http.StatusBadRequest)
		return
	}

	note, err := notes.Get(c, dbConn, field, fieldValue)
	if err != nil && err != notes.ErrorNoMatch {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if err == notes.ErrorNoMatch {
		http.Error(w, err.Error()+": "+fieldValue, http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleNotesPost(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	note := []notes.Collection{}

	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := notes.Post(c, dbConn, note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleNotesPut(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	note := notes.Collection{}

	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := notes.Put(c, dbConn, note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
