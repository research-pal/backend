package main

import (
	"fmt"
	"log"
	"net/http"
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
		log.Printf("Defaulting to port %s", port)
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
	switch r.URL.Path {
	case "/notes":
		switch r.Method {
		case http.MethodGet:
			api.HandleNotesGet(m.dbConn, w, r)
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
