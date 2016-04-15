package main

import (
	"fmt"
	"net/http"
	"strings"
)

func parseStatuLine(code int) string {
	return fmt.Sprint(code, " ", http.StatusText(code))
}

type ErrorHandler struct {
	w http.ResponseWriter
	Msg string
}

func NewErrorHandler(w http.ResponseWriter, msg string) *ErrorHandler {
	return &ErrorHandler{w, msg}
}

func (e ErrorHandler) SendText(code int) {
	http.Error(e.w, fmt.Sprint(parseStatuLine(code), "\n", e.Msg), code)
}

//500
func (e ErrorHandler) ServerError() {
	e.SendText(http.StatusInternalServerError)
}

//401
func (e ErrorHandler) Unauthorized() {
	e.SendText(http.StatusUnauthorized)
}

//403
func (e ErrorHandler) Forbidden() {
	e.SendText(http.StatusForbidden)
}

//404
func (e ErrorHandler) NotFound() {
	e.SendText(http.StatusNotFound)
}

//405
func (e ErrorHandler) MethodNotAllowed(methods []string){
	allows := strings.Join(methods, ",")
	e.Msg = fmt.Sprint(e.Msg, "\n", "Allow:", allows)
	e.w.Header().Add("Allow", allows)
	e.SendText(http.StatusMethodNotAllowed)
}