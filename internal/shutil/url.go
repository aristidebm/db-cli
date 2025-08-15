package shutil

import (
	"errors"
	"net/url"
	"strings"
)

var URLParseError = errors.New("cannot not parse connection url")

func ParseURL(URL string) (*url.URL, error) {
	u, err := url.Parse(URL)
	if err != nil {
		// Handle non-standard URLs like sqlite3:////path
		if strings.HasPrefix(URL, "sqlite3:") {
			u.Scheme = "sqlite3"
			u.Path = strings.TrimPrefix(URL, "sqlite3:")
			// Remove leading slashes for absolute paths
			u.Path = strings.TrimLeft(u.Path, "/")
			if !strings.HasPrefix(u.Path, "/") {
				u.Path = "/" + u.Path
			}
			return nil, URLParseError
		}
		return nil, URLParseError
	}
	return u, nil
}
