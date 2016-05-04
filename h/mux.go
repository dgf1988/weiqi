package h

import (
	"net/http"
	"strings"
)

//Mux
type Mux struct {
	Router *route
}

func NewMux() *Mux {
	return &Mux{Router: newRoute()}
}

func (mux *Mux) Handle(h Handler, pattern string, methods int) {
	mux.Router.Handle(h, pattern, methods)
}

func (mux *Mux) HandleFunc(f HandlerFunc, pattern string, methods int) {
	mux.Router.Handle(HandlerFunc(f), pattern, methods)
}

func (mux *Mux) HandleStd(h http.Handler, pattern string, methods int) {
	mux.Router.Handle(HandlerFunc(func(w http.ResponseWriter, r *http.Request, p []string) {
		h.ServeHTTP(w, r)
	}), pattern, methods)
}

func (mux *Mux) HandleFuncStd(f http.HandlerFunc, pattern string, methods int) {
	mux.Router.Handle(HandlerFunc(func(w http.ResponseWriter, r *http.Request, p []string) {
		f(w, r)
	}), pattern, methods)
}

func (mux *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var useragent = strings.ToLower(r.UserAgent())
	var remoteaddr = getIp(r)
	//Slurp
	if strings.Contains(useragent, "bot") || strings.Contains(useragent, "spider") || strings.Contains(useragent, "slurp") {
		spiderlogger.Println(remoteaddr, r.Method, r.URL, r.UserAgent())
	} else {
		accesslogger.Println(remoteaddr, r.Method, r.URL, r.Header.Get("referer"), r.UserAgent())
	}

	route, params := mux.Router.Match(r.URL.Path)

	if route == nil || route.Handler == nil {
		NotFound(w, "page not found")
	} else {
		if ( parseMethod(r.Method) & route.Methods ) > 0 {
			route.Handler.ServeHTTP(w, r, params)
		} else {
			MethodNotAllowed(w, "method not allowed", formatMethods(route.Methods))
		}
	}
}

