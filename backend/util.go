package backend

import (
	"net/url"
)

func CleanURL(URL string) (string, error) {

	u, err := url.Parse(URL)
	if err != nil {
		return "", err
	}

	cleanURL := u.Host + u.Path

	if string(cleanURL[len(cleanURL)-1]) == "/" {
		cleanURL = cleanURL[:len(cleanURL)-1]
	}

	return cleanURL, nil
}
