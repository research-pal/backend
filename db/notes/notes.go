package notes

import (
	"time"

	"errors"
)

var ErrorNoMatch = errors.New("No Matching Record")

type Collection struct {
	TaskID        string    `firestore:"task_id"`
	URL           string    `firestore:"url"`
	Notes         string    `firestore:"notes"`
	Status        string    `firestore:"status"`
	Assignee      string    `firestore:"assignee"`
	Group         string    `firestore:"group"`
	PriorityOrder string    `firestore:"priority_order"`
	CreatedDate   time.Time `firestore:"created_date"`
	LastUpdate    time.Time `firestore:"last_update"`
}

// ID generates the document id in the format desired
// will be used as the document id to save the record, and also to retrieve using it
// returns empty if any of the key fields are empty
func (r Collection) ID() string {
	if r.TaskID == "" {
		return ""
	}
	return r.TaskID
}

const (
	CollectionName = "notes"
)

// // Query queries the database to return the list of matching records.
// // Note: full text search is not supported at db layer. it is taken care at the service layer
// // returns only valid videos (Valid == true)
// func Query(ctx context.Context, dbConn *firestore.Client) ([]Collection, error) {
// 	video := dbConn.Collection(CollectionVideo)

// 	video.Query = video.Query.Where("valid", "==", true)

// 	iter := video.Documents(ctx)
// 	docs, err := iter.GetAll()
// 	if err != nil {
// 		return nil, err
// 	}

// 	results := []Collection{}
// 	for _, d := range docs {
// 		r := Collection{}
// 		d.DataTo(&r)
// 		// r.DBID = d.Ref.ID
// 		results = append(results, r)
// 	}
// 	return results, nil
// }

// // QueryAll supports querying all videos, including the invalid ones
// // supported filters:
// // 1) valid
// func QueryAll(ctx context.Context, dbConn *firestore.Client, valid *bool) ([]Collection, error) {
// 	video := dbConn.Collection(CollectionVideo)

// 	if valid != nil {
// 		video.Query = video.Query.Where("valid", "==", *valid)
// 	}

// 	iter := video.Documents(ctx)
// 	docs, err := iter.GetAll()
// 	if err != nil {
// 		return nil, err
// 	}

// 	results := []Collection{}
// 	for _, d := range docs {
// 		r := Collection{}
// 		d.DataTo(&r)
// 		// r.DBID = d.Ref.ID
// 		results = append(results, r)
// 	}
// 	return results, nil
// }
