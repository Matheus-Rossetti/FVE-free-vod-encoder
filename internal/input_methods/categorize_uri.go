package input_methods

import (
	"net/url"
	"strings"
)

func CategorizeUri(uri string) string {

	//check if its url or path
	// if its path, validate file, push it to encodeQueue and return
	parsedUri, _ := url.Parse(uri)
	scheme := strings.ToLower(parsedUri.Scheme)

	switch scheme {
	case "http", "https":
		return "url"
	case "file", "":
		return "path"
	default:
		return "unsupported uri"
	}
}
