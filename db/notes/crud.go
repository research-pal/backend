// generic crud code
// do not edit. only generic code in this file. all customizations done in seperate methods in other go files

package notes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/muly/go-common/errors"
	"google.golang.org/api/iterator"
)

// Post posts the given list of records into the database collection
// returns list of errors (in the format errors.ErrMsgs) for all the failed records
func Post(ctx context.Context, dbConn *firestore.Client, list []Collection) error {
	var errs errors.ErrMsgs
	for _, r := range list {
		r.CreatedDate = time.Now()
		r.LastUpdate = time.Now()
		log.Printf("POST CRUD")
		_, err := dbConn.Collection(CollectionName).Doc(r.ID()).Create(ctx, r)
		if err != nil {
			log.Printf("POST CRUD error: %#v", err)
			errType := errors.ErrGeneric
			if strings.Contains(err.Error(), "code = AlreadyExists desc = Document already exists") {
				errType = errors.ErrExists
			}
			errs = append(errs, errors.NewError(errType, r.ID()).(errors.ErrMsg))
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// Put updates the record
// if unique fields which being doc id is missing in the parameters, return error
// matches the record based on the doc id and updates the field with what is provided in the input struct
func Put(ctx context.Context, dbConn *firestore.Client, docID string, r Collection) error {
	if docID == "" {
		return fmt.Errorf("key fields are missing: key %s", docID)
	}
	if !exists(ctx, dbConn, docID) {
		return fmt.Errorf("document does not exists to update: key %s", docID)
	}

	r.LastUpdate = time.Now()
	_, err := dbConn.Collection(CollectionName).Doc(docID).Set(ctx, r)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes the record
// if doc id is blank in the input, returns generic error
// if doc id is not found in the database, returns not found error
// matches the record based on the doc id and delete the record
func Delete(ctx context.Context, dbConn *firestore.Client, docID string) error {
	if docID == "" {
		return errors.NewError(errors.ErrEmptyInput, "docID")
	}
	if !exists(ctx, dbConn, docID) {
		return errors.NewError(errors.ErrNotFound, docID)
	}

	_, err := dbConn.Collection(CollectionName).Doc(docID).Delete(ctx)
	if err != nil {
		return errors.NewError(errors.ErrGeneric, err.Error())
	}

	return nil
}

// GetByID gets the record based on the doc id provided
// if doc id is blank in the input, return error
// if record is not found, error is returned
// Note: unlike Query(), Get doesn't apply Valid=True filter
func GetByID(ctx context.Context, dbConn *firestore.Client, docID string) (Collection, error) {
	if docID == "" {
		return Collection{}, fmt.Errorf("docID is missing, provide id")
	}

	r, err := dbConn.Collection(CollectionName).Doc(docID).Get(ctx)
	if err != nil {
		return Collection{}, err
	}
	v := Collection{}
	r.DataTo(&v)

	return v, nil
}

// Get gets the records based on the keys and their values provided
func Get(ctx context.Context, dbConn *firestore.Client, field string, fieldValue string) ([]Collection, error) {
	if field == "" {
		return []Collection{}, fmt.Errorf("field is missing, provide field")
	}

	vOne := Collection{}
	v := []Collection{}

	iter := dbConn.Collection(CollectionName).Where(field, "==", fieldValue).Documents(ctx)
	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}
		if err != nil {
			return []Collection{}, err
		}

		out, err := json.Marshal(doc.Data())
		if err != nil {
			return []Collection{}, err
		}
		err = json.Unmarshal(out, &vOne)
		if err != nil {
			return []Collection{}, err
		}

		v = append(v, vOne)
	}
	return v, nil
}

func exists(ctx context.Context, dbConn *firestore.Client, docID string) bool {
	_, err := dbConn.Collection(CollectionName).Doc(docID).Get(ctx)
	if err != nil {
		return false
	}
	return true
}
