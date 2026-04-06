package core

import (
	nUrl "net/url"
	"strings"
)

type URIType int

const (
	Unsupported URIType = iota
	Url
	Path
)

func CategorizeUri(uri string) (URIType, string) {

	parsedUri, _ := nUrl.Parse(uri)
	scheme := strings.ToLower(parsedUri.Scheme)

	switch scheme {
	case "http", "https":
		return Url, scheme
	case "file", "c", "":
		return Path, scheme
	default:
		return Unsupported, scheme
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
