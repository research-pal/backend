// generic crud code
// do not edit. only generic code in this file. all customizations done in seperate methods in other go files

package notes

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/muly/go-common/errors"
)

// Post posts the given list of records into the database collection
// returns list of errors (in the format errors.ErrMsgs) for all the failed records
func Post(ctx context.Context, dbConn *firestore.Client, list []Collection) error {
	var errs errors.ErrMsgs
	checkURL := map[string]string{}

	for _, r := range list {
		if r.EncodedURL == "" {
			log.Printf("url is emtpy")
			return fmt.Errorf("record must have url value")
		}
		checkURL["encodedurl"] = r.EncodedURL
		existing, err := Get(ctx, dbConn, checkURL)
		if err != nil {
			log.Printf("error getting record by url: %v", err)
			return fmt.Errorf("document does not exists to update: url %s", r.EncodedURL)
		}
		if len(existing) == 0 {
			r.CreatedDate = time.Now()
			r.LastUpdate = time.Now()
			r.Status = "new"
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
		} else if existing[0].EncodedURL == r.EncodedURL {
			log.Printf("record already exists by url: %v", r.EncodedURL)
			return fmt.Errorf("record already exists with encodedurl %s", r.EncodedURL)
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
func Put(ctx context.Context, dbConn *firestore.Client, r Collection) error {
	if r.DocID == "" {
		return fmt.Errorf("key fields are missing: key %s", r.DocID)
	}

	existing, err := GetByID(ctx, dbConn, r.DocID)
	if err != nil {
		log.Printf("error getting record by id: %v", err)
		return fmt.Errorf("document does not exists to update: key %s", r.DocID)
	}

	log.Printf("PUT CRUD")
	r.CreatedDate = existing.CreatedDate
	r.LastUpdate = time.Now()
	// TODO : PUT and PATCH, when id provided in the request and is different than in the URL, we should throw error
	_, err = dbConn.Collection(CollectionName).Doc(r.DocID).Set(ctx, r)
	if err != nil {
		return err
	}

	return nil
}

// Patch updates the record with only provided fields
func Patch(ctx context.Context, dbConn *firestore.Client, id string, r map[string]interface{}) (Collection, error) {
	batch := dbConn.Batch()

	if id == "" {
		return Collection{}, fmt.Errorf("key fields are missing: key %s", id)
	}

	v := Collection{}
	log.Printf("PATCH CRUD")
	if exists(ctx, dbConn, id) {
		r["last_update"] = time.Now()
		// TODO:find better method to update instead of using batch approach
		// TODO : PUT and PATCH, when id provided in the request and is different than in the URL, we should throw error
		batch.Set(dbConn.Collection(CollectionName).Doc(id), r, firestore.MergeAll)
		_, err := batch.Commit(ctx)
		if err != nil {
			return v, err
		}
		v, err = GetByID(ctx, dbConn, id)
		if err != nil {
			return Collection{}, err
		}
	} else {
		return Collection{}, fmt.Errorf("document does not exists to update: key %s", id)
	}

	return v, nil
}

// Delete deletes the record
// if doc id is blank in the input, returns generic error
// if doc id is not found in the database, returns not found error
// matches the record based on the doc id and delete the record
func Delete(ctx context.Context, dbConn *firestore.Client, id string) error {
	if id == "" {
		return errors.NewError(errors.ErrEmptyInput, "id")
	}
	if !exists(ctx, dbConn, id) {
		return errors.NewError(errors.ErrNotFound, id)
	}

	log.Printf("DELETE CRUD")
	_, err := dbConn.Collection(CollectionName).Doc(id).Delete(ctx)
	if err != nil {
		return errors.NewError(errors.ErrGeneric, err.Error())
	}

	return nil
}

// GetByID gets the record based on the doc id provided
// if doc id is blank in the input, return error
// if record is not found, error is returned
// Note: unlike Query(), Get doesn't apply Valid=True filter
func GetByID(ctx context.Context, dbConn *firestore.Client, id string) (Collection, error) {
	if id == "" {
		return Collection{}, fmt.Errorf("id is missing, provide id")
	}

	log.Printf("GETBYID CRUD")
	r, err := dbConn.Collection(CollectionName).Doc(id).Get(ctx)
	if err != nil {
		return Collection{}, err
	}
	v := Collection{}
	v.DocID = id
	r.DataTo(&v)

	return v, nil
}

// Get gets the records based on the keys and their values provided
func Get(ctx context.Context, dbConn *firestore.Client, filters map[string]string) ([]Collection, error) {
	if filters == nil {
		return []Collection{}, fmt.Errorf("required parameter is missing in URI")
	}

	results := []Collection{}

	log.Printf("GETBYFILTER CRUD")
	query := dbConn.Collection(CollectionName).Query
	for key, value := range filters {
		query = query.Where(key, "==", value)
	}
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return []Collection{}, err
	}

	for _, doc := range docs {
		r := Collection{}
		if err := doc.DataTo(&r); err != nil {
			return []Collection{}, err
		}
		r.DocID = doc.Ref.ID
		results = append(results, r)
	}

	return results, nil
}

func exists(ctx context.Context, dbConn *firestore.Client, id string) bool {
	_, err := dbConn.Collection(CollectionName).Doc(id).Get(ctx)
	if err != nil {
		return false
	}
	return true
}
