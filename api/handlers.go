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

// Init initializes database connection
func Init(client *firestore.Client) {
	dbConn = client
}

// HandleNotesGetByID get the data by id provided.
func HandleNotesGetByID(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	params := mux.Vars(r)
	if params["id"] == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	id := params["id"]

	note, err := notes.GetByID(c, dbConn, id)
	if err != nil && err != notes.ErrorNoMatch {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if err == notes.ErrorNoMatch {
		http.Error(w, err.Error()+": "+id, http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HandleNotesGetFiltered get the filtered data. below filters are supported:
// encodedurl, assignee, status, group, priority_order
func HandleNotesGetFiltered(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	filters := map[string]string{}

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for k, v := range params {
		if !keyExists(k) {
			fmt.Fprintf(w, "given key `%s` is either doesn't exist or a typo, try correcting.", k)
			// http.Error(w, err.Error(), http.StatusBadRequest) // not woring
			return
		}
		filters[k] = v[0]
	}

	// Group Query params: encodedurl, status, group, assignee, priority_order

	note, err := notes.Get(c, dbConn, filters)
	if err != nil && err != notes.ErrorNoMatch {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if err == notes.ErrorNoMatch {
		http.Error(w, fmt.Sprintf("no records found with given filters: %v", filters), http.StatusNotFound)
		return
	}

	if len(note) == 0 {
		fmt.Fprintf(w, "There are no data for such query")
		return
	}

	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HandleNotesPost saves the data in the db. below paramaters in json would do the action:
// {"assignee":"","group":"","notes":"","priority_order":"","status":"","url":""}
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

// HandleNotesPut gets the data by id provided and replaces the content given in below parameters.
// {"assignee":"","group":"","notes":"","priority_order":"","status":"","url":""}
func HandleNotesPut(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	note := notes.Collection{}

	params := mux.Vars(r)
	if params["id"] == "" {
		http.Error(w, notes.ErrorMissing.Error(), http.StatusBadRequest)
		return
	}
	id := params["id"]

	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := notes.Put(c, dbConn, id, note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	note.DocID = id
	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HandleNotesDelete identifies by id provided and deletes it from db.
func HandleNotesDelete(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	params := mux.Vars(r)
	if params["id"] == "" {
		http.Error(w, notes.ErrorMissing.Error(), http.StatusBadRequest)
		return
	}
	id := params["id"]

	if err := notes.Delete(c, dbConn, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Document with id %s is deleted\n", id)
}

// HandleNotesPatch updates only give key value pairs
func HandleNotesPatch(w http.ResponseWriter, r *http.Request) { // TODO
	// c := appengine.NewContext(r)
	note := notes.Collection{}

	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, notes.ErrorMissing.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "%v", note)

	// if err := notes.Put(c, dbConn, id, note); err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// note.DocID = id
	// if err := json.NewEncoder(w).Encode(note); err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

}

func keyExists(k string) bool {
	fields := []string{"encodedurl", "assignee", "status", "group", "priority_order"}
	for i := range fields {
		if fields[i] == k {
			return true
		}
	}
	return false
}
