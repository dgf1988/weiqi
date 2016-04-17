package weiqi

import (
	"database/sql"
	"net/http"
	"github.com/dgf1988/weiqi/h"
)

func defaultHandler(w http.ResponseWriter, r *http.Request, p []string) {
	u := getSessionUser(r)

	posts, err := dbListPostByPage(40, 0)
	if err != nil && err != sql.ErrNoRows {
		h.ServerError(w, err)
		return
	}

	players, err := dbListPlayer(40, 0)
	if err != nil && err != sql.ErrNoRows {
		h.ServerError(w, err)
		return
	}

	sgfs, err := dbListSgf(40, 0)
	if err != nil && err != sql.ErrNoRows {
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
