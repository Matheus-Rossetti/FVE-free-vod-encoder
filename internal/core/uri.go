package core

import (
	"errors"
	nUrl "net/url"
	"strings"
)

var ErrUnsupportedUri = errors.New("unsupported uri")

type URIType int

const (
	Unsupported URIType = iota
	Url
	Path
)

func CategorizeUri(uri string) (URIType, error) {

	parsedUri, _ := nUrl.Parse(uri)
	scheme := strings.ToLower(parsedUri.Scheme)

	switch scheme {
	case "http", "https":
		return Url, nil
	case "file", "":
		return Path, nil
	default:
		return Unsupported, ErrUnsupportedUri
	}
}

func (u URIType) String() string {
	switch u {
	case Url:
		return "url"
	case Path:
		return "path"
	default:
		return "unsupported"
	}
}
