package weiqi

import (
	"database/sql"
)

type Handler interface {
	ServeHTTP(h *Http)
}

type HandlerFunc func(h *Http)

func (this HandlerFunc) ServeHTTP(h *Http) {
	this(h)
}

func defaultHandler(h *Http) {
	u := getSessionUser(h.R)

	posts, err := dbListPostByPage(40, 0)
	if err != nil && err != sql.ErrNoRows {
		h.ServerError(err.Error())
		return
	}

	players, err := dbListPlayer(40, 0)
	if err != nil && err != sql.ErrNoRows {
		h.ServerError(err.Error())
		return
	}

	sgfs, err := dbListSgf(40, 0)
	if err != nil && err != sql.ErrNoRows {
		h.ServerError(err.Error())
		return
	}

	datamap := defaultData(u, posts, players, sgfs)
	err = defaultHtml().Execute(h.W, datamap, defFuncMap)
	if err != nil {
		h.ServerError(err.Error())
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
