package main

import (
	"database/sql"
	"html/template"
)

var(
	defFuncMap = template.FuncMap{
		"HasLogin": func(u *U) bool {
			return u != nil && u.Name != ""
		},
	}
)

type Handler interface {
	ServeHTTP(h *Http)
}

type HandlerFunc func(h *Http)

func (this HandlerFunc) ServeHTTP(h *Http) {
	this(h)
}

func defaultHandler(h *Http) {
	var u *U
	s := getSession(h.R)
	if s != nil {
		u = s.User
	}

	posts, err := dbListPostByPage(10, 0)
	if err != nil && err != sql.ErrNoRows {
		h.ServerError(err.Error())
		return
	}

	players, err := dbListPlayer(10, 0)
	if err != nil && err != sql.ErrNoRows {
		h.ServerError(err.Error())
		return
	}

	sgfs, err := dbListSgf(10, 0)
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
