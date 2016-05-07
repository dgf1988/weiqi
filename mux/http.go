package mux

import (
	"net/http"
)

type Param  []string

type Http struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Param          Param
}

func NewHttp(w http.ResponseWriter, r *http.Request, p Param) *Http {
	return &Http{ResponseWriter:w, Request:r, Param:p}
}