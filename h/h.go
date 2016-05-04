package h

import (
	"github.com/dgf1988/weiqi/logger"
	"log"
	"strings"
	"net/http"
	"regexp"
)

var (
	errorlogger    = logger.New("error")
	accesslogger   = logger.New("access")
	spiderlogger   = logger.New("spider")
	notfoundlogger = logger.New("notfound")
)

func init() {
	errorlogger.SetPrefix("[Error: 500]")
	errorlogger.SetFlags(log.LstdFlags)

	accesslogger.SetPrefix("[Access]")
	accesslogger.SetFlags(log.LstdFlags)

	spiderlogger.SetPrefix("[Spider]")
	spiderlogger.SetFlags(log.LstdFlags)

	notfoundlogger.SetPrefix("[Notfound: 404]")
	notfoundlogger.SetFlags(log.LstdFlags)
}



const (
	GET = 1 << iota
	POST
)

func formatMethods(methods int) []string {
	var list_method = make([]string, 0)
	if methods & GET > 0 {
		list_method = append(list_method, "GET")
	}
	if methods & POST > 0 {
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
		panic("http: method error")
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