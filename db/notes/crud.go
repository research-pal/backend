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
	checkURL := map[string]string{}

	for _, r := range list {
		if r.URL == "" {
			log.Printf("url is emtpy")
			return fmt.Errorf("record must have url value")
		}
		checkURL["url"] = r.URL
		existing, err := Get(ctx, dbConn, checkURL)
		if err != nil {
			log.Printf("error getting record by url: %v", err)
			return fmt.Errorf("document does not exists to update: url %s", r.URL)
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
		} else if existing[0].URL == r.URL {
			log.Printf("record already exists by url: %v", r.URL)
			return fmt.Errorf("record already exists with encodedurl %s", r.URL)
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
func Put(ctx context.Context, dbConn *firestore.Client, id string, r Collection) error {
	if id == "" {
		return fmt.Errorf("key fields are missing: key %s", id)
	}

	existing, err := GetByID(ctx, dbConn, id)
	if err != nil {
		log.Printf("error getting record by id: %v", err)
		return fmt.Errorf("document does not exists to update: key %s", id)
	}

	log.Printf("PUT CRUD")
	r.CreatedDate = existing.CreatedDate
	r.LastUpdate = time.Now()
	_, err = dbConn.Collection(CollectionName).Doc(id).Set(ctx, r)
	if err != nil {
		return err
	}

	return nil
}

// Patch updates the record
func Patch(ctx context.Context, dbConn *firestore.Client, id string, r Collection) error {
	if id == "" {
		return fmt.Errorf("key fields are missing: key %s", id)
	}

	existing, err := GetByID(ctx, dbConn, id)
	if err != nil {
		log.Printf("error getting record by id: %v", err)
		return fmt.Errorf("document does not exists to update: key %s", id)
	}

	log.Printf("PUT CRUD")
	r.CreatedDate = existing.CreatedDate
	r.LastUpdate = time.Now()
	_, err = dbConn.Collection(CollectionName).Doc(id).Set(ctx, r)
	if err != nil {
		return err
	}

	return nil
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

	log.Printf("GET BY ID CRUD")
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

	fields, fieldvals := []string{}, []string{}
	for k, v := range filters {
		if k == "encodedurl" {
			k = "url"
		}
		fields = append(fields, k)
		fieldvals = append(fieldvals, v)
	}

	vOne := Collection{}
	v := []Collection{}

	log.Printf("GET BY FILTER CRUD")
	var iter *firestore.DocumentIterator
	if len(filters) == 0 {
		iter = dbConn.Collection(CollectionName).Documents(ctx)
	} else if len(filters) == 1 { //TODO: use a for loop instead of hardcoding using else if
		iter = dbConn.Collection(CollectionName).Where(fields[0], "==", fieldvals[0]).Documents(ctx)
	} else if len(filters) == 2 {
		iter = dbConn.Collection(CollectionName).Where(fields[0], "==", fieldvals[0]).Where(fields[1], "==", fieldvals[1]).Documents(ctx)
	} else if len(filters) > 2 {
		return []Collection{}, fmt.Errorf("query params are %d, supports only 2 params", len(filters))
	}
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
		vOne.DocID = doc.Ref.ID
		v = append(v, vOne)
	}
	return v, nil
}

func exists(ctx context.Context, dbConn *firestore.Client, id string) bool {
	_, err := dbConn.Collection(CollectionName).Doc(id).Get(ctx)
	if err != nil {
		return false
	}
	return true
}
