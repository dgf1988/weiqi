package weiqi

import (
	"fmt"
	"github.com/dgf1988/weiqi/mux"
	"net/http"
)

//登录页面
func handleLogin(w http.ResponseWriter, r *http.Request, p []string) {

	//会话验证
	if getSession(r) != nil {
		mux.SeeOther(w, r, "/user")
		return
	}

	//post
	if r.Method == POST {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		u, err := loginUser(username, password)
		if err == nil && u != nil {
			newSession(u).Add(w)
			mux.SeeOther(w, r, "/user")
			return
		}
		if werr, ok := err.(*WeiqiError); ok {
			mux.SeeOther(w, r, fmt.Sprint("/login?loginmsg=", werr.Msg))
		} else {
			mux.ServerError(w, err)
		}
	} else if r.Method == GET {
		gcSession()
		r.ParseForm()
		loginMsg := r.FormValue("loginmsg")
		registerMsg := r.FormValue("registermsg")

		if err := renderLogin(w, loginMsg, registerMsg); err != nil {
			mux.ServerError(w, err)
		}
	}
}

func renderLogin(w http.ResponseWriter, loginmsg, registermsg string) error {
	data := defData()
	data.Head.Title = "登录"
	data.Content["LoginMsg"] = loginmsg
	data.Content["RegisterMsg"] = registermsg

	return defHtmlLayout().Append(
		defHtmlHead(),
		defHtmlHeader(),
		defHtmlFooter(),
		newHtmlContent("login"),
	).Execute(w, data, nil)
}

func handleLogout(w http.ResponseWriter, r *http.Request, p []string) {
	clearSession(w, r)
	mux.SeeOther(w, r, "/login")
}

func handleRegister(w http.ResponseWriter, r *http.Request, p []string) {

	//会话验证
	if getSession(r) != nil {
		mux.SeeOther(w, r, "/user")
		return
	}

	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	password2 := r.FormValue("password2")
	email := r.FormValue("email")

	_, err := registerUser(username, password, password2, email, r.RemoteAddr)
	if err == nil {
		mux.SeeOther(w, r, "/login?registermsg=注册成功")
		return
	}
	if werr, ok := err.(*WeiqiError); ok {
		mux.SeeOther(w, r, fmt.Sprint("/login?registermsg=", werr.Msg))
	} else {
		mux.ServerError(w, err)
	}
}

func handleUser(w http.ResponseWriter, r *http.Request, p []string) {

	//会话验证
	var user *User
	if user = getSessionUser(r); user == nil {
		mux.SeeOther(w, r, "/login")
		return
	}

	err := renderUser(w, user)
	if err != nil {
		mux.ServerError(w, err)
	}
}

func renderUser(w http.ResponseWriter, u *User) error {
	d := defData()
	d.User = u
	d.Header.Navs = userNavItems()
	return defHtmlLayout().Append(
		defHtmlHead(),
		defHtmlHeader(),
		defHtmlFooter(),
		newHtmlContent("user"),
	).Execute(w, d, nil)
}
