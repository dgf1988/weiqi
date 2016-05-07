package mux

import (
	"net/http"
	"strings"
)

//SeeOther 303 跳转页面
func SeeOther(w http.ResponseWriter, r *http.Request, urlStr string) {
	http.Redirect(w, r, urlStr, http.StatusSeeOther)
}

//Unauthorized 401 没有权限
func Unauthorized(w http.ResponseWriter, msg string) {
	textStatus(w, msg, http.StatusUnauthorized)
}

//Forbidden 403 拒绝请求
func Forbidden(w http.ResponseWriter, msg string) {
	textStatus(w, msg, http.StatusForbidden)
}

//NotFound 404 找不到页面
func NotFound(w http.ResponseWriter, msg string) {
	textStatus(w, msg, http.StatusNotFound)
}

//MethodNotAllowed 405 方法错误
func MethodNotAllowed(w http.ResponseWriter, msg string, allows []string) {
	allow := strings.Join(allows, ",")
	w.Header().Set("Allow", allow)
	textStatus(w, allow+"\n"+msg, http.StatusMethodNotAllowed)
}

//ServerError 500 服务器错误
func ServerError(w http.ResponseWriter, err error) {
	errorlogger.Printf(err.Error())
	textStatus(w, err.Error(), http.StatusInternalServerError)
}
