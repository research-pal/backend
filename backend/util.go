package backend

import (
	"net/url"
)

func CleanURL(URL string) (string, error) {

	u, err := url.Parse(URL)
	if err != nil {
		return "", err
	}

	return u.Host + u.Path, nil
}
