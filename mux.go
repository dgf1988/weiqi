package main

import (
	"net/http"
)

type Mux struct {
	Router *Route
}

func NewMux() *Mux {
	return &Mux{Router: NewRoute()}
}

func (this Mux) Handle(handler Handler, pattern string, methods ...string) {
	this.Router.Add(handler, pattern, methods...)
}

func (this Mux) HandleFunc(handler func(h *Http), pattern string, methods ...string) {
	this.Router.Add(HandlerFunc(handler), pattern, methods...)
}

func (this Mux) HandleOld(handler http.Handler, pattern string, methods ...string) {
	this.HandleFunc(func(h *Http) {
		handler.ServeHTTP(h.W, h.R)
	}, pattern, methods...)
}

func (this Mux) HandleFuncOld(handler func(w http.ResponseWriter, r *http.Request), pattern string, methods ...string) {
	this.HandleFunc(func(h *Http) {
		handler(h.W, h.R)
	}, pattern, methods...)
}

func (this Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, params := this.Router.Get(r.URL.Path)
	h := &Http{w, r, params}

	if route == nil || route.Handler == nil || route.Methods == nil || len(route.Methods) == 0 {
		h.RequestError("page not found").NotFound()

	} else {
		for _, m := range route.Methods {
			if r.Method == m {
				logAccess(h)
				route.Handler.ServeHTTP(h)
				return
			}
		}
		h.RequestError("请求的方法不支持").MethodNotAllowed(route.Methods)
	}
}
