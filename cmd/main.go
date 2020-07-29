package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"cloud.google.com/go/firestore"

	"github.com/research-pal/backend/api"
	"github.com/research-pal/backend/db"
)

func main() {

	dbClient := db.NewDBClient()
	defer dbClient.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting port to %s", port)
	}
	server := http.Server{
		Addr: fmt.Sprintf(":%v", port),
		Handler: &myHandler{
			dbConn: dbClient,
		},
	}
	log.Printf("Listening on port %s", port)
	log.Fatal(server.ListenAndServe())
}

type myHandler struct {
	dbConn *firestore.Client
}

func (m *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ServeHTTP() started with path", r.Method, r.URL.Path)

	w.Header().Set("Content-Type", "application/json")

	// params := mux.Vars(r) // TODO

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch r.URL.Path {
	case "/notes":
		switch r.Method {
		case http.MethodGet:
			if params["encodedurl"][0] != "" {
				api.HandleNotesGetByID(m.dbConn, w, r)
			} else {
				api.HandleNotesGetFiltered(m.dbConn, w, r)
			}
		case http.MethodPost:
			api.HandleNotesPost(m.dbConn, w, r)
		case http.MethodPut:
			api.HandleNotesPut(m.dbConn, w, r)
		}
	}

	// r.Handle("/notes", mw.ThenFunc(HandleNotesPut)).Methods("PUT")
	// r.Handle("/notes", mw.ThenFunc(HandleNotesDelete)).Methods("DELETE")
	// r.Handle("/notesall", mw.ThenFunc(HandleNotesGetAll)).Methods("GET") //TODO: /notesall is not looking like correct name. may be just /notes with a special query parameter like /notes?all=true

}
