package db

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
)

// NewDBClient creates the new db connection
// it is callers responsibility to close the connection
// TODO: is it possible to for the caller to copy this session and close that copied session after use, so that this session stays active  through out the life of the app
func NewDBClient() *firestore.Client {
	projectID := os.Getenv("GCP_PROJECT") // export GCP_PROJECT=research-pal-2
	if projectID == "" {
		//Note: in case of app deployed to app engine on gcp, this line will detect the current project id, so this env need not be set in api.yaml
		projectID = firestore.DetectProjectID
	}

	// Get a Firestore client.
	dbClient, err := firestore.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatalf("Failed to create firestore client: %v", err)
	}

	return dbClient
}
