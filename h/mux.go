package h

import "net/http"

type Mux struct {
	Router *route
}

func NewMux() *Mux {
	return &Mux{Router: newRoute()}
}

func (this Mux) Handle(handler Handler, pattern string, methods ...string) {
	this.Router.Set(handler, pattern, methods...)
}

func (this Mux) HandleFunc(handler HandlerFunc, pattern string, methods ...string) {
	this.Router.Set(HandlerFunc(handler), pattern, methods...)
}

func (this Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, params := this.Router.Match(r.URL.Path)

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