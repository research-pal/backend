package api

import (
	"encoding/json"
	"fmt"
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
	if params["docid"] == "" {
		http.Error(w, "docid is required", http.StatusBadRequest)
		return
	}
	docID := params["docid"]

	note, err := notes.GetByID(c, dbConn, docID)
	if err != nil && err != notes.ErrorNoMatch {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if err == notes.ErrorNoMatch {
		http.Error(w, err.Error()+": "+docID, http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleNotesGetFiltered(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	fields := []string{"encodedurl", "status", "assignee", "group"}
	field, fieldValue := "", ""

	// TODO: use mux functions to retrieve query params
	// mparams := mux.Vars(r)
	// if mparams ==nil{
	// 	fmt.Println("params are nil")
	// 	http.Error(w, "query params are required", http.StatusBadRequest)
	// 	return
	// }
	// fmt.Println(mparams, "\t:\t", r.URL.RequestURI())
	// for _, f = range fields {
	// 	if params[f] != "" {
	// 		field = f
	// 		if f == "encodedurl" {
	// 			field = "url"
	// 		}
	// 		fieldValue = params[f]
	// 	}
	// }

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Group Query params: externalTaskURL(encodedurl), status, group, assignee
	for _, f := range fields {
		if val, exists := params[f]; exists {
			field = f
			if f == "encodedurl" {
				field = "url"
			}
			fieldValue = val[0]
		}
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
	params := mux.Vars(r)

	// TODO: no query parameters..

	if params["docid"] == "" {
		http.Error(w, notes.ErrorMissing.Error(), http.StatusBadRequest)
		return
	}
	docID := params["docid"]

	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := notes.Put(c, dbConn, docID, note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleNotesDelete(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	params := mux.Vars(r)

	// TODO: no query parameters..

	if params["docid"] == "" {
		http.Error(w, notes.ErrorMissing.Error(), http.StatusBadRequest)
		return
	}
	docID := params["docid"]

	// if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	if err := notes.Delete(c, dbConn, docID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Document with docID %s is deleted\n", docID)
}
