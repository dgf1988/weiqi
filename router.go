package weiqi

import (
	"strings"
)

const (
	ROUTER_DEFAULT_PATTERN = ""
	ROUTER_PARAM_PATTERN   = "*"
)

const (
	POST = "POST"
	GET  = "GET"
)

type Route struct {
	Pattern string
	Methods []string
	Handler Handler

	DefRoute   *Route
	ParamRoute *Route
	Routes     map[string]*Route
}

func NewRoute() *Route {
	return &Route{Routes: make(map[string]*Route)}
}

func (this *Route) Add(handler Handler, pattern string, methods ...string) {
	listPath := strings.Split(pattern, "/")
	cr := this
	for i := range listPath {
		switch listPath[i] {
		//默认
		case ROUTER_DEFAULT_PATTERN:
			if cr.DefRoute == nil {
				cr.DefRoute = NewRoute()
			}
			cr = cr.DefRoute
		//参数
		case ROUTER_PARAM_PATTERN:
			if cr.ParamRoute == nil {
				cr.ParamRoute = NewRoute()
			}
			cr = cr.ParamRoute
		//静态
		default:
			r, ok := cr.Routes[listPath[i]]
			if !ok {
				r = NewRoute()
				cr.Routes[listPath[i]] = r
			}
			cr = r
		}
	}
	cr.Pattern = pattern
	cr.Methods = methods
	cr.Handler = handler
}

func (this *Route) Get(pattern string) (*Route, []string) {
	listPath := strings.Split(pattern, "/")
	listParam := make([]string, 0)
	cr := this
	for i := range listPath {
		if listPath[i] == ROUTER_DEFAULT_PATTERN {
			cr = cr.DefRoute
		} else {
			r, ok := cr.Routes[listPath[i]]
			if ok {
				cr = r
			} else {
				listParam = append(listParam, listPath[i])
				cr = cr.ParamRoute
			}
		}
		if cr == nil {
			return nil, nil
		}
	}
	return cr, listParam
}
