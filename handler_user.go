package weiqi

import "fmt"

//登录页面
func loginHandler(h *Http) {

	//会话验证
	if getSession(h.R) != nil {
		h.SeeOther("/user")
		return
	}

	//post
	if h.R.Method == POST {
		h.R.ParseForm()
		username := h.R.FormValue("username")
		password := h.R.FormValue("password")

		u, err := loginUser(username, password)
		if err == nil && u != nil {
			newSession(u).Add(h.W)
			h.SeeOther("/user")
			return
		}
		if werr, ok := err.(*WeiqiError); ok {
			h.SeeOther(fmt.Sprint("/login?loginmsg=", werr.Msg))
		} else {
			h.ServerError(err.Error())
		}
	} else if h.R.Method == GET {
		clearSessionMany()
		h.R.ParseForm()
		loginMsg := h.R.FormValue("loginmsg")
		registerMsg := h.R.FormValue("registermsg")

		err := loginHtml().Execute(h.W, loginData(loginMsg, registerMsg), defFuncMap)
		if err != nil {
			h.ServerError(err.Error())
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

func handlerLogout(h *Http) {
	clearSession(h.W, h.R)
	h.SeeOther("/login")
}

func handlerRegister(h *Http) {

	//会话验证
	if getSession(h.R) != nil {
		h.SeeOther("/user")
		return
	}

	h.R.ParseForm()
	username := h.R.FormValue("username")
	password := h.R.FormValue("password")
	password2 := h.R.FormValue("password2")
	email := h.R.FormValue("email")

	_, err := RegisterUser(username, password, password2, email, h.R.RemoteAddr)
	if err == nil {
		h.SeeOther("/login?registermsg=注册成功")
		return
	}
	if werr, ok := err.(*WeiqiError); ok {
		h.SeeOther(fmt.Sprint("/login?registermsg=", werr.Msg))
	} else {
		h.ServerError(err.Error())
	}
}

func userHandler(h *Http) {

	//会话验证
	s := getSession(h.R)
	if s == nil {
		h.SeeOther("/login")
		return
	}

	err := userHtml().Execute(h.W, userData(s.User), defFuncMap)
	if err != nil {
		h.ServerError(err.Error())
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
