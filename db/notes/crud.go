// generic crud code
// do not edit. only generic code in this file. all customizations done in seperate methods in other go files

package notes

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/muly/go-common/errors" // TODO: need to remove dependency on this package and then delete that package
)

// Post posts the given list of records into the database collection
// returns list of errors (in the format errors.ErrMsgs) for all the failed records
func Post(ctx context.Context, dbConn *firestore.Client, list []Collection) ([]Collection, error) {
	var errs errors.ErrMsgs
	results := []Collection{}

	for _, r := range list {
		var err error
		r, err = post(ctx, dbConn, r)
		if err != nil {
			errs = append(errs, err.(errors.ErrMsg))
		} else {
			results = append(results, r)
		}
	}
	if len(errs) > 0 {
		return results, // results: the list with sucessfully posted records in first.
			errs // errs: the failed records in error
	}
	return results, nil
}

// Put updates the record
// if unique fields which being doc id is missing in the parameters, return error
// matches the record based on the doc id and updates the field with what is provided in the input struct
func Put(ctx context.Context, dbConn *firestore.Client, r Collection) error {
	// validate
	if r.DocID == "" { // TODO: use check on ID() after #19 is fixed
		return fmt.Errorf("key fields are missing: key %s", r.DocID)
	}

	// check of already existance
	exists, data := existsByID(ctx, dbConn, r.ID())
	if !exists {
		return fmt.Errorf("%v not found %w", r.ID(), ErrorNotFound)
	}

	if !r.existsByKeyFields(ctx, dbConn) {
		return fmt.Errorf("document doesn't already exists by key fields, %w", ErrorInvalidData)
	}
	r.CreatedDate = data.CreatedDate
	r.LastUpdate = time.Now()

	_, err := dbConn.Collection(CollectionName).Doc(r.DocID).Set(ctx, r)
	if err != nil {
		return err
	}

	return nil
}

// Patch updates the record with only provided fields
func Patch(ctx context.Context, dbConn *firestore.Client, id string, updates map[string]interface{}) (Collection, error) {
	if id == "" {
		return Collection{}, fmt.Errorf("key fields are missing: key %s", id)
	}

	exists, _ := existsByID(ctx, dbConn, id)
	if !exists {
		return Collection{}, errors.NewError(errors.ErrNotFound, id)
	}

	updates["last_update"] = time.Now()
	batch := dbConn.Batch()
	batch.Set(dbConn.Collection(CollectionName).Doc(id), updates, firestore.MergeAll)
	_, err := batch.Commit(ctx)
	if err != nil {
		return Collection{}, err
	}
	results, err := GetByID(ctx, dbConn, id)
	if err != nil {
		return Collection{}, err
	}

	return results, nil
}

// Delete deletes the record
// if doc id is blank in the input, returns generic error
// if doc id is not found in the database, returns not found error
// matches the record based on the doc id and deletes the record
func Delete(ctx context.Context, dbConn *firestore.Client, id string) error {
	if id == "" {
		return errors.NewError(errors.ErrEmptyInput, "id")
	}
	exists, _ := existsByID(ctx, dbConn, id)
	if !exists {
		return errors.NewError(errors.ErrNotFound, id)
	}

	_, err := dbConn.Collection(CollectionName).Doc(id).Delete(ctx)
	if err != nil {
		return errors.NewError(errors.ErrGeneric, err.Error())
	}

	return nil
}

// GetByID gets the record based on the doc id provided
// if doc id is blank in the input, returns generic error
// if record is not found, returns not found error
// Note: unlike Query(), Get doesn't apply Valid=True filter
func GetByID(ctx context.Context, dbConn *firestore.Client, id string) (Collection, error) {
	if id == "" {
		return Collection{}, errors.NewError(errors.ErrEmptyInput, "id")
	}

	found, results := existsByID(ctx, dbConn, id)
	if !found {
		return Collection{}, errors.NewError(errors.ErrNotFound, id)
	}

	return results, nil
}

// Get gets the records based on the keys and their values provided
func Get(ctx context.Context, dbConn *firestore.Client, filters url.Values) ([]Collection, error) {
	query := dbConn.Collection(CollectionName).Query
	for key, value := range filters {
		query = query.Where(key, "==", value[0])
	}
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return []Collection{}, err
	}

	results := []Collection{}
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

func post(ctx context.Context, dbConn *firestore.Client, r Collection) (Collection, error) {
	// check of already existance
	if r.existsByKeyFields(ctx, dbConn) {
		return Collection{}, fmt.Errorf("record already exists")
	}

	if !r.isValidPost() {
		return Collection{}, fmt.Errorf("record is invalid")
	}

	//
	r.CreatedDate = time.Now()
	r.LastUpdate = time.Now()
	r.DocID = r.ID()
	_, err := dbConn.Collection(CollectionName).Doc(r.DocID).Create(ctx, r)
	if err != nil {
		errType := errors.ErrGeneric
		if strings.Contains(err.Error(), "code = AlreadyExists desc = Document already exists") {
			errType = errors.ErrExists
		}
		return Collection{}, errors.NewError(errType, r.DocID).(errors.ErrMsg)
	}
	return r, nil
}

func existsByID(ctx context.Context, dbConn *firestore.Client, id string) (bool, Collection) {
	doc, err := dbConn.Collection(CollectionName).Doc(id).Get(ctx)
	if err != nil {
		return false, Collection{}
	}

	c := Collection{}
	doc.DataTo(&c)
	c.DocID = doc.Ref.ID
	return true, c
}
