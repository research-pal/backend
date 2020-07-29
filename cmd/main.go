package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/research-pal/backend/api"
	"github.com/research-pal/backend/db"
)

func main() {
	dbClient := db.NewDBClient()
	defer dbClient.Close()
	api.Init(dbClient)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting port to %s", port)
	}

	router := mux.NewRouter()

	router.HandleFunc("/notes", api.HandleNotesGetFiltered).Methods("GET")
	router.HandleFunc("/notes/{taskid}", api.HandleNotesGetByID).Methods("GET") // docid
	router.HandleFunc("/notes", api.HandleNotesPost).Methods("POST")
	router.HandleFunc("/notes", api.HandleNotesPut).Methods("PUT")

	log.Printf("Listening on port %s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Println("Listening error :", err)
	}
}
