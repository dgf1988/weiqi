package h

import "strings"

const (
	matchDefault = ""
	matchParam   = "+"
	matchAll     = "*"
	pathSplit    = "/"
)

type route struct {
	Pattern string
	Methods int
	Handler Handler

	DefRoute   *route
	ParamRoute *route
	AllRoute   *route
	Routes     map[string]*route
}

func newRoute() *route {
	return &route{Routes: make(map[string]*route)}
}

func (this *route) Handle(handler Handler, pattern string, methods int) {
	listPath := strings.Split(pattern, pathSplit)
	cr := this
	for i := range listPath {
		switch listPath[i] {
		//默认
		case matchDefault:
			if cr.DefRoute == nil {
				cr.DefRoute = newRoute()
			}
			cr = cr.DefRoute
		//参数
		case matchParam:
			if cr.ParamRoute == nil {
				cr.ParamRoute = newRoute()
			}
			cr = cr.ParamRoute
		case matchAll:
			if cr.AllRoute == nil {
				cr.AllRoute = newRoute()
			}
			cr = cr.AllRoute
		//静态
		default:
			r, ok := cr.Routes[listPath[i]]
			if !ok {
				r = newRoute()
				cr.Routes[listPath[i]] = r
			}
			cr = r
		}
	}
	cr.Pattern = pattern
	cr.Methods = methods
	cr.Handler = handler
}

func (this *route) Match(pattern string) (*route, []string) {
	listPath := strings.Split(pattern, pathSplit)
	listParam := make([]string, 0)
	cr := this
	for i := range listPath {
		if listPath[i] == matchDefault {
			cr = cr.DefRoute
		} else if r, ok := cr.Routes[listPath[i]]; ok {
			cr = r
		} else if cr.ParamRoute != nil {
			listParam = append(listParam, listPath[i])
			cr = cr.ParamRoute
		} else if cr.AllRoute != nil {
			return cr.AllRoute, listParam
		} else {
			return nil, nil
		}
		if cr == nil {
			return nil, nil
		}
	}
	return cr, listParam
}
