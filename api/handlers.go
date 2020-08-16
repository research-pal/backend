package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"github.com/research-pal/backend/db/notes"
	"google.golang.org/appengine"
)

var dbConn *firestore.Client

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
	if err != nil {
		http.Error(w, err.Error()+": "+id, convertToHTTPStatus(err))
		return
	}

	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HandleNotesGetFiltered get the filtered data. below filters are supported:
// url, assignee, status, group, priority_order
func HandleNotesGetFiltered(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if valid, incorrectQueryParam := isAllowedQueryParam(params); !valid {
		http.Error(w, fmt.Sprintf("invalid query parameters: `%v`", incorrectQueryParam), http.StatusBadRequest)
		return
	}

	note, err := notes.Get(c, dbConn, params)
	if err != nil {
		http.Error(w, err.Error(), convertToHTTPStatus(err))
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

// HandleNotesPost saves the data in the db. below fields in json would do the action:
// {"assignee":"","group":"","notes":"","priority_order":"","status":"","url":""}
func HandleNotesPost(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	note := []notes.Collection{}

	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i := range note {
		note[i].Unescape()
	}

	results, err := notes.Post(c, dbConn, note)
	if err != nil {
		http.Error(w, err.Error(), convertToHTTPStatus(err))
		return
	}

	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HandleNotesPut gets the data by id provided and replaces the content given in below parameters.
// {"assignee":"","group":"","notes":"","priority_order":"","status":"","url":""}
func HandleNotesPut(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	note := notes.Collection{}

	id := mux.Vars(r)["id"]

	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if note.DocID == "" {
		note.DocID = id
	}

	if note.ID() != id {
		http.Error(w, "id in payload is incorrect", http.StatusBadRequest)
		return
	}

	note.Unescape()

	if err := notes.Put(c, dbConn, note); err != nil {
		http.Error(w, err.Error(), convertToHTTPStatus(err))
		return
	}

	note, err := notes.GetByID(c, dbConn, id)
	if err != nil {
		http.Error(w, err.Error(), convertToHTTPStatus(err))
		return
	}

	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HandleNotesDelete identifies by id provided and deletes it from db.
func HandleNotesDelete(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	id := mux.Vars(r)["id"]

	if err := notes.Delete(c, dbConn, id); err != nil {
		http.Error(w, err.Error(), convertToHTTPStatus(err))
		return
	}
}

// HandleNotesPatch updates only the given fields from the request body
// key fields are not allowed to be updated
func HandleNotesPatch(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	id := mux.Vars(r)["id"]

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "problem to read request body", http.StatusBadRequest)
	}

	data := map[string]interface{}{}
	if err = json.Unmarshal(content, &data); err != nil {
		http.Error(w, fmt.Sprintf("Unmarshal error.. %s: %v", content, err), http.StatusBadRequest)
		return
	}

	if valid, invalidFields := isValidPatchData(data); !valid {
		http.Error(w, fmt.Sprintf("invalid fields %v", invalidFields), http.StatusBadRequest)
		return
	}

	note, err := notes.Patch(c, dbConn, id, data)
	if err != nil {
		http.Error(w, err.Error(), convertToHTTPStatus(err))
		return
	}

	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
