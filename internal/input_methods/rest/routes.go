package rest

import (
	"net/http"
)

func (r *rest) AddRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", r.checkHealth)
	mux.HandleFunc("POST /encode", r.videoHandler())
}
