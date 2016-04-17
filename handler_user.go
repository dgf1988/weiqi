package weiqi

import (
	"fmt"
	"github.com/dgf1988/weiqi/h"
	"net/http"
)

//登录页面
func loginHandler(w http.ResponseWriter, r *http.Request, p []string) {

	//会话验证
	if getSession(r) != nil {
		h.SeeOther(w, r, "/user")
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
			h.SeeOther(w, r, "/user")
			return
		}
		if werr, ok := err.(*WeiqiError); ok {
			h.SeeOther(w, r, fmt.Sprint("/login?loginmsg=", werr.Msg))
		} else {
			h.ServerError(w, err)
		}
	} else if r.Method == GET {
		clearSessionMany()
		r.ParseForm()
		loginMsg := r.FormValue("loginmsg")
		registerMsg := r.FormValue("registermsg")

		err := loginHtml().Execute(w, loginData(loginMsg, registerMsg), defFuncMap)
		if err != nil {
			h.ServerError(w, err)
		}
	}
}

func loginHtml() *Html {
	return defHtmlLayout().Append(
		defHtmlHead(),
		defHtmlHeader(),
		defHtmlFooter(),
		newHtmlContent("login"),
	)
}

func loginData(loginmsg, registermsg string) *Data {
	data := defData()
	data.Head.Title = "登录"
	data.Content["LoginMsg"] = loginmsg
	data.Content["RegisterMsg"] = registermsg
	return data
}

func handlerLogout(w http.ResponseWriter, r *http.Request, p []string) {
	clearSession(w, r)
	h.SeeOther(w, r, "/login")
}

func handlerRegister(w http.ResponseWriter, r *http.Request, p []string) {

	//会话验证
	if getSession(r) != nil {
		h.SeeOther(w, r, "/user")
		return
	}

	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	password2 := r.FormValue("password2")
	email := r.FormValue("email")

	_, err := RegisterUser(username, password, password2, email, r.RemoteAddr)
	if err == nil {
		h.SeeOther(w, r, "/login?registermsg=注册成功")
		return
	}
	if werr, ok := err.(*WeiqiError); ok {
		h.SeeOther(w, r, fmt.Sprint("/login?registermsg=", werr.Msg))
	} else {
		h.ServerError(w, err)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request, p []string) {

	//会话验证
	s := getSession(r)
	if s == nil {
		h.SeeOther(w, r, "/login")
		return
	}

	err := userHtml().Execute(w, userData(s.User), defFuncMap)
	if err != nil {
		h.ServerError(w, err)
	}
}

func userHtml() *Html {
	return defHtmlLayout().Append(
		defHtmlHead(),
		defHtmlHeader(),
		defHtmlFooter(),
		newHtmlContent("user"),
	)
}

func userData(u *U) *Data {
	d := defData()
	d.User = u
	d.Header.Navs = userNavItems()
	return d
}
