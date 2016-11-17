package backend

import (
	"net/url"
)

func CleanURL(rawurl string) (string, error) {

	// decode the URL if it is already encoded.
	URL, err := url.QueryUnescape(rawurl)
	if err != nil {
		return "", err
	}

	// parse the url
	u, err := url.Parse(URL)
	if err != nil {
		return "", err
	}

	// just get the host name and the rest of the path, ignoring the queries, bookmarks, schema etc
	cleanURL := u.Host + u.Path

	if string(cleanURL[len(cleanURL)-1]) == "/" {
		cleanURL = cleanURL[:len(cleanURL)-1]
	}

	return cleanURL, nil
}
