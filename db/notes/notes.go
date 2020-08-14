package notes

import (
	"context"
	"errors"
	"log"
	"net/url"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
)

// Error reports when record is not matching
var (
	ErrorNotFound     = errors.New("Not Found")
	ErrorInvalidData  = errors.New("Invalid Data")
	ErrorAlreadyExist = errors.New("Already Exists")
)

// Collection holds the table fields
type Collection struct {
	DocID         string    `firestore:"-" json:"id"`
	Assignee      string    `firestore:"assignee" json:"assignee"`
	CreatedDate   time.Time `firestore:"created_date" json:"created_date"`
	Group         string    `firestore:"group" json:"group"`
	LastUpdate    time.Time `firestore:"last_update" json:"last_updated"`
	Notes         string    `firestore:"notes" json:"notes"`
	PriorityOrder string    `firestore:"priority_order" json:"priority_order"`
	Status        string    `firestore:"status" json:"status"`
	URL           string    `firestore:"url" json:"url"`
}

// ID generates the document id in the format desired if the DocID is not there
// will be used as the document id to save the record, and also to retrieve using it
func (r Collection) ID() string {
	if r.DocID != "" {
		return r.DocID
	}
	id, err := uuid.NewUUID()
	if err != nil {
		log.Printf("uuid error : %#v\n", err)
	}
	return id.String()
}

// CollectionName from the cloud
const (
	CollectionName = "notes"
)

// Unescape replaces the url escaped fields with the unescaped version.
// if there is any error while cleaning, save the original value and log the error
// the fields supported as of now are below:
// 1) url
func (r *Collection) Unescape() {
	c, err := url.QueryUnescape(r.URL)
	if err != nil {
		log.Printf("unable to unescape %s: %v", r.URL, err)
	} else {
		r.URL = c
	}
}

// TODO: need to remove this field after the #19 is addressed
func (r Collection) existsByKeyFields(ctx context.Context, dbConn *firestore.Client) bool {
	filters := url.Values{"url": []string{r.URL}}
	existing, err := Get(ctx, dbConn, filters)
	if err != nil {
		return false
	}
	if len(existing) > 0 {
		return true
	}
	return false
}

func (r Collection) isValid() (bool, []string) {
	invalidReasons := []string{}
	valid := true
	if r.URL == "" {
		invalidReasons = append(invalidReasons, `URL == ""`)
		valid = false
	}
	return valid, invalidReasons
}

func (r Collection) isValidPost() (bool, []string) {
	invalidReasons := []string{}
	valid := true
	if r.Status != "new" {
		invalidReasons = append(invalidReasons, `Status != "new"`)
		valid = false
	}
	if v, l := r.isValid(); !v {
		invalidReasons = append(invalidReasons, l...)
		valid = false
	}
	return valid, invalidReasons
}
