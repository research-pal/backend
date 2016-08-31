package backend

import (
	"errors"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type Notes struct {
	URL        string
	Notes      string
	LastUpdate time.Time
}

var ErrorNoMatch = errors.New("No Matching Record")

func (n *Notes) get(URL string, c context.Context) error {

	cleanURL, err := CleanURL(URL)
	if err != nil {
		return err
	}

	key := datastore.NewKey(c, "Notes", cleanURL, 0, nil)

	if err = datastore.Get(c, key, n); err != nil && err.Error() == "datastore: no such entity" {
		err = ErrorNoMatch
	}

	return err
}

func (n *Notes) put(c context.Context) error {

	cleanURL, err := CleanURL(n.URL)
	if err != nil {
		return err
	}

	// generate the key
	key := datastore.NewKey(c, "Notes", cleanURL, 0, nil)
	n.LastUpdate = time.Now()

	// put the record into the database and capture the key
	if key, err = datastore.Put(c, key, n); err != nil {
		return err
	}

	// read from database into the same variable
	if err = n.get(n.URL, c); err != nil {
		return err
	}

	return nil

}
