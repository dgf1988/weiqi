package mux

import (
	"github.com/dgf1988/weiqi/logger"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var (
	errorlogger  = logger.New("error")
	accesslogger = logger.New("access")
	spiderlogger = logger.New("spider")
)

func init() {
	errorlogger.SetPrefix("[Error]")
	errorlogger.SetFlags(log.LstdFlags)

	accesslogger.SetPrefix("[Access]")
	accesslogger.SetFlags(log.LstdFlags)

	spiderlogger.SetPrefix("[Spider]")
	spiderlogger.SetFlags(log.LstdFlags)
}

const (
	GET = 1 << iota
	POST
)

func formatMethods(methods int) []string {
	var list_method = make([]string, 0)
	if methods&GET > 0 {
		list_method = append(list_method, "GET")
	}
	if methods&POST > 0 {
		list_method = append(list_method, "POST")
	}
	return list_method
}

func parseMethod(method string) int {
	switch strings.ToUpper(method) {
	case "GET":
		return GET
	case "POST":
		return POST
	default:
		return 0
	}
}

func getIp(r *http.Request) string {
	var ip = r.Header.Get("x-forwarded-for")
	if ip == "" {
		ip = r.RemoteAddr
		return regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}|::1`).FindString(ip)
	}
	return ip
}

//Mux
type Mux struct {
	router *route
}

func New() *Mux {
	return &Mux{router: newRoute()}
}

func (m Mux) Handle(methods int, pattern string, handler Handler) {
	m.router.Handle(methods, pattern, handler)
}

func (m Mux) HandleFunc(methods int, pattern string, f func(h *Http)) {
	m.router.Handle(methods, pattern, HandlerFunc(f))
}

func (m Mux) HandleStd(methods int, pattern string, handler http.Handler) {
	m.HandleFunc(methods, pattern, func(h *Http) {
		handler.ServeHTTP(h.ResponseWriter, h.Request)
	})
}

func (m Mux) HandleFuncStd(methods int, pattern string, f func(w http.ResponseWriter, h *http.Request)) {
	m.HandleFunc(methods, pattern, func(h *Http) {
		f(h.ResponseWriter, h.Request)
	})
}

func (m Mux) HandleFuncOld(methods int, pattern string, f func(w http.ResponseWriter, h *http.Request, args []string)) {
	m.HandleFunc(methods, pattern, func(h *Http) {
		f(h.ResponseWriter, h.Request, h.Param)
	})
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

	route, params := mux.router.Match(r.URL.Path)

	if route == nil || route.Handler == nil {
		NotFound(w, "page not found")
	} else {
		if (parseMethod(r.Method) & route.Methods) > 0 {
			route.Handler.ServeHTTP(newHttp(w, r, params))
		} else {
			MethodNotAllowed(w, "method not allowed", formatMethods(route.Methods))
		}
	}
}
