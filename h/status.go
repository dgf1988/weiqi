package h

import (
	"log"
	"net/http"
	"os"
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
	f, ferr := os.OpenFile(errFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if ferr != nil {
		panic(ferr.Error())
	}
	defer f.Close()
	log.New(f, "[ServerError: 500]", log.LstdFlags).Println(err.Error())
	textStatus(w, err.Error(), http.StatusInternalServerError)
}
