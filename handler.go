package weiqi

import (
	"github.com/dgf1988/weiqi/h"
	"net/http"
)

//默认处理器。处理首页访问
func defaultHandler(w http.ResponseWriter, r *http.Request, args []string) {
	//从会话中获取用户信息，如果没登录，则为nil。
	u := getSessionUser(r)

	//
	posts, err := dbListPostByPage(40, 0)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	//
	players, err := dbListPlayer(40, 0)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	//
	sgfs, err := dbListSgf(40, 0)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	datamap := defaultData(u, posts, players, sgfs)
	err = defaultHtml().Execute(w, datamap, defFuncMap)
	if err != nil {
		h.ServerError(w, err)
		return
	}
}

func defaultHtml() *Html {
	return defHtmlLayout().Append(
		defHtmlHead(),
		defHtmlHeader(),
		defHtmlContent(),
		defHtmlFooter(),
	)
}

func defaultData(u *U, posts []P, players []Player, sgfs []Sgf) *Data {
	data := defData()
	data.User = u
	data.Content["Posts"] = posts
	data.Content["Sgfs"] = sgfs
	data.Content["Players"] = players
	return data
}
