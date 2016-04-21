package weiqi

import (
	"github.com/dgf1988/weiqi/h"
	"net/http"
)

//默认处理器。处理首页访问
func handleDefault(w http.ResponseWriter, r *http.Request, args []string) {
	//从会话中获取用户信息，如果没登录，则为nil。
	u := getSessionUser(r)

	err := render_default(w, u)
	if err != nil {
		h.ServerError(w, err)
	}
}

func render_default(w http.ResponseWriter, u *User) error {
	data := defData()
	data.User = u

	posts, err := Posts.ListMap(40, 0)
	if err != nil {
		return err
	}

	//
	players, err := Players.ListMap(40, 0)
	if err != nil {
		return err
	}

	//
	sgfs, err := Sgfs.ListMap(40, 0)
	if err != nil {
		return err
	}

	data.Content["Posts"] = posts
	data.Content["Sgfs"] = sgfs
	data.Content["Players"] = players

	return defHtmlLayout().Append(
		defHtmlHead(),
		defHtmlHeader(),
		defHtmlContent(),
		defHtmlFooter(),
	).Execute(w, data, defFuncMap)
}
