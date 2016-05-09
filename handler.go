package weiqi

import (
	"github.com/dgf1988/weiqi/h"
	"net/http"
)

//默认处理器。处理首页访问
func defaultHandler(w http.ResponseWriter, r *http.Request, args []string) {
	//从会话中获取用户信息，如果没登录，则为nil。
	var err error
	var data = defData()
	data.User = getSessionUser(r)

	var posts []Post
	if posts, err = listPostByStatusOrderDesc(constStatusRelease, 40, 0); err != nil {
		h.ServerError(w, err)
		return
	}
	data.Content["Posts"] = posts

	var players []PlayerTable
	if players, err = listPlayerOrderByRankDesc(40, 0); err != nil {
		h.ServerError(w, err)
		return
	}
	data.Content["Players"] = players

	var sgfs []Sgf
	if sgfs, err = listSgfOrderByTimeDesc(40, 0); err != nil {
		h.ServerError(w, err)
		return
	}
	data.Content["Sgfs"] = sgfs

	var html = defHtmlLayout().Append(defHtmlHead(), defHtmlHeader(), defHtmlFooter(), defHtmlContent())
	if err = html.Execute(w, data, nil); err != nil {
		logError(err.Error())
	}
}
