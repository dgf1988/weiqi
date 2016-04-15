package main

import (
	"net/http"
)

type Http struct {
	//响应
	W http.ResponseWriter
	//请求
	R *http.Request
	//参数
	P []string
}

func (h *Http) RequestError(msg string) *ErrorHandler {
	return NewErrorHandler(h.W, msg)
}

func (h *Http) ServerError(msg string) {
	logError(h, msg)
	NewErrorHandler(h.W, msg).ServerError()
}

func (h *Http) SeeOther(url string) {
	http.Redirect(h.W, h.R, url, http.StatusSeeOther)
}
