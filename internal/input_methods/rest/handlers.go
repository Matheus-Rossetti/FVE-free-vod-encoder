package rest

import (
	"fmt"
	"net/http"
)

func (rest *rest) videoHandler(w http.ResponseWriter, r *http.Request) {

	// might refactor to return a struct
	uriType, uri, keyStarter, status, err := validateBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	rest.pushJob(uriType, uri, keyStarter)

	fmt.Fprintf(w, "job added to internal queue!")
}

func (rest *rest) checkHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ok!")
}
