package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"cloud.google.com/go/firestore"
	mapset "github.com/deckarep/golang-set"
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
		if errors.Is(err, notes.ErrorNotFound) {
			http.Error(w, err.Error()+": "+id, http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	if err != nil && err != notes.ErrorNoMatch {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if err == notes.ErrorNoMatch {
		http.Error(w, fmt.Sprintf("no records found with given filters: %v", params), http.StatusNotFound)
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
// {"assignee":"","group":"","notes":"","priority_order":"","status":"","encodedurl":""}
func HandleNotesPost(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	note := []notes.Collection{}

	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	results, err := notes.Post(c, dbConn, note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HandleNotesPut gets the data by id provided and replaces the content given in below parameters.
// {"assignee":"","group":"","notes":"","priority_order":"","status":"","encodedurl":""}
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

	if note.DocID == "" {
		note.DocID = id
	}

	if note.ID() != id {
		http.Error(w, "id in payload is incorrect", http.StatusBadRequest)
		return
	}

	if err := notes.Put(c, dbConn, note); err != nil {
		if errors.Is(err, notes.ErrorNoMatch) || errors.Is(err, notes.ErrorInvalidData) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	note, err := notes.GetByID(c, dbConn, id)
	if err != nil && err != notes.ErrorNoMatch {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if err == notes.ErrorNoMatch {
		http.Error(w, err.Error()+": "+id, http.StatusBadRequest)
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

	params := mux.Vars(r)
	if params["id"] == "" {
		http.Error(w, notes.ErrorMissing.Error(), http.StatusBadRequest)
		return
	}
	id := params["id"]

	if err := notes.Delete(c, dbConn, id); err != nil {
		// TODO: need to check if the error is of type errors.ErrNotFound, and return 400 instead
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HandleNotesPatch updates only the given fields from the request body
// key fields are not allowed to be updated
func HandleNotesPatch(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, notes.ErrorMissing.Error(), http.StatusBadRequest)
		return
	}

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "problem to read request body", http.StatusBadRequest)
	}

	data := map[string]interface{}{}
	if err = json.Unmarshal(content, &data); err != nil {
		http.Error(w, "Unmarshal error..", http.StatusBadRequest)
		return
	}

	if valid, invalidFields := isValidPatchData(data); !valid {
		http.Error(w, fmt.Sprintf("invalid fields %v", invalidFields), http.StatusBadRequest)
		return
	}

	note, err := notes.Patch(c, dbConn, id, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// isAllowedQueryParam validates if the given parameter list has any params not supported.
// if found any, returns false along with the list of invalid params
// else returns true
func isAllowedQueryParam(params url.Values) (bool, []string) {
	validParams := mapset.NewSetFromSlice([]interface{}{"encodedurl", "assignee", "status", "group", "priority_order"})
	incorrectQueryParam := []string{}

	for k := range params {
		if !validParams.Contains(k) {
			incorrectQueryParam = append(incorrectQueryParam, k)
		}
	}
	if len(incorrectQueryParam) > 0 {
		return false, incorrectQueryParam
	}
	return true, nil
}

func isValidPatchData(data map[string]interface{}) (bool, []string) {
	validFields := mapset.NewSetFromSlice([]interface{}{"assignee", "status", "group", "priority_order"})
	incorrectFields := []string{}
	for field := range data {
		if !validFields.Contains(field) {
			incorrectFields = append(incorrectFields, field)
		}
	}

	if len(incorrectFields) > 0 {
		return false, incorrectFields
	}
	return true, nil
}
