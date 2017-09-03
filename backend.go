package backend

import (
	"github.com/gorilla/mux"
	"net/http"
)

func init() {
	r := mux.NewRouter()
	//r.HandleFunc("/", HandleRoot).Methods("GET")
	r.HandleFunc("/notes", HandleNotesGet).Methods("GET")
	r.HandleFunc("/notes", HandleNotesPut).Methods("PUT")
	r.HandleFunc("/notes", HandleNotesDelete).Methods("DELETE")
	r.HandleFunc("/notesall", HandleNotesGetAll).Methods("GET") //TODO: /notesall is not looking correct name

	http.Handle("/", r)
}
