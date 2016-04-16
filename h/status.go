package h

import (
	"net/http"
	"strings"
)


//SeeOther 303
func SeeOther(w http.ResponseWriter, r *http.Request, urlStr string) {
	http.Redirect(w, r, urlStr, http.StatusSeeOther)
}

//Unauthorized 401
func Unauthorized(w http.ResponseWriter, msg string) {
	textStatus(w, msg, http.StatusUnauthorized)
}

//Forbidden 403
func Forbidden(w http.ResponseWriter, msg string) {
	textStatus(w, msg, http.StatusForbidden)
}

//NotFound 404
func NotFound(w http.ResponseWriter, msg string) {
	textStatus(w, msg, http.StatusNotFound)
}

//MethodNotAllowed 405
func MethodNotAllowed(w http.ResponseWriter, msg string, allows []string) {
	allow := strings.Join(allows, ",")
	w.Header().Set("Allow", allow)
	textStatus(w, allow + "\n" + msg, http.StatusMethodNotAllowed)
}

//ServerError 500
func ServerError(w http.ResponseWriter, err error) {
	textStatus(w, err.Error(), http.StatusInternalServerError)
}
