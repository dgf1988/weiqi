package h

import (
	"net/http"
)

type Param []string

func (p Param) Len() int {
	return len(p)
}

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, p []string)
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request, p []string)

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request, p []string) {
	h(w, r, p)
}
