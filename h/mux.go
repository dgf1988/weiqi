package h

import "net/http"

//Mux
type Mux struct {
	Router *route
}

func NewMux() *Mux {
	return &Mux{Router: newRoute()}
}

func (mux *Mux) Handle(h Handler, pattern string, methods ...string) {
	mux.Router.Handle(h, pattern, methods...)
}

func (mux *Mux) HandleFunc(f HandlerFunc, pattern string, methods ...string) {
	mux.Router.Handle(HandlerFunc(f), pattern, methods...)
}

func (mux *Mux) HandleStd(h http.Handler, pattern string, methods ...string) {
	mux.Router.Handle(HandlerFunc(func(w http.ResponseWriter, r *http.Request, p []string) {
		h.ServeHTTP(w, r)
	}), pattern, methods...)
}

func (mux *Mux) HandleFuncStd(f http.HandlerFunc, pattern string, methods ...string) {
	mux.Router.Handle(HandlerFunc(func(w http.ResponseWriter, r *http.Request, p []string) {
		f(w, r)
	}), pattern, methods...)
}

func (mux *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, params := mux.Router.Match(r.URL.Path)

	if route == nil || route.Handler == nil || route.Methods == nil || len(route.Methods) == 0 {
		NotFound(w, "page not found")
	} else {
		for _, m := range route.Methods {
			if r.Method == m {
				route.Handler.ServeHTTP(w, r, params)
				return
			}
		}
		MethodNotAllowed(w, "method not allowed", route.Methods)
	}
}