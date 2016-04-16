package http

import "net/http"

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, p []string)
}

type