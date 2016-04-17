package weiqi

import (
	"database/sql"
	"fmt"
	"github.com/dgf1988/weiqi/h"
	"net/http"
)

//sgf list
func sgfListHandler(w http.ResponseWriter, r *http.Request, p []string) {
	var u *U
	s := getSession(r)
	if s != nil {
		u = s.User
	}

	sgfs, err := dbListSgf(40, 0)
	if err == sql.ErrNoRows {
		h.NotFound(w, "sgf list not found")
		return
	}
	if err != nil {
		h.ServerError(w, err)
		return
	}
	if err = sgfListHtml().Execute(w, sgfListData(u, sgfs), defFuncMap); err != nil {
		h.ServerError(w, err)
		return
	}
}

func sgfListHtml() *Html {
	return defHtmlLayout().Append(
		defHtmlHead(),
		defHtmlHeader(),
		defHtmlFooter(),
		newHtmlContent("sgflist"),
	)
}

func sgfListData(u *U, sgfs []Sgf) *Data {
	data := defData()
	data.User = u
	data.Head.Title = "棋谱列表"
	data.Head.Desc = "围棋棋谱列表"
	data.Head.Keywords = []string{"围棋", "棋谱", "比赛"}
	data.Content["Sgfs"] = sgfs
	return data
}

//sgf id
func sgfIdHandler(w http.ResponseWriter, r *http.Request, p []string) {
	var u *U
	s := getSession(r)
	if s != nil {
		u = s.User
	}

	id := atoi64(p[0])
	if id <= 0 {
		h.NotFound(w, p[0]+" sgf not found")
		return
	}
	sgf, err := dbGetSgf(id)
	if err == sql.ErrNoRows {
		h.NotFound(w, p[0]+" sgf not found")
		return
	}
	if err != nil {
		h.ServerError(w, err)
		return
	}
	if err = sgfIdHtml().Execute(w, sgfIdData(u, sgf), defFuncMap); err != nil {
		h.ServerError(w, err)
		return
	}
}

func sgfIdHtml() *Html {
	return defHtmlLayout().Append(
		defHtmlHead(),
		defHtmlHeader(),
		defHtmlFooter(),
		newHtmlContent("sgfid"),
	)
}

func sgfIdData(u *U, sgf *Sgf) *Data {
	data := defData()
	data.User = u
	data.Head.Title = fmt.Sprintf("%s - %s VS %s", sgf.Event, sgf.Black, sgf.White)
	data.Head.Desc = "围棋棋谱之" + sgf.Event
	data.Head.Keywords = []string{"围棋", "棋谱", "比赛", sgf.Black, sgf.White}
	data.Content["Sgf"] = sgf
	return data
}

//sgf edit
func userSgfEditHandler(w http.ResponseWriter, r *http.Request, p []string) {
	s := getSession(r)
	if s == nil {
		h.SeeOther(w, r, "/login")
		return
	}

	r.ParseForm()
	var (
		action      = "/user/sgf/add"
		msg         = r.FormValue("editormsg")
		sgf    *Sgf = nil
		err    error
	)

	if len(p) > 0 {
		action = "/user/sgf/update"
		sgf, err = dbGetSgf(atoi64(p[0]))
		if err == sql.ErrNoRows || err == ErrPrimaryKey {
			h.NotFound(w, "sgf not found")
			return
		}
		if err != nil {
			h.ServerError(w, err)
			return
		}
	} else {
		sgf = new(Sgf)
	}

	sgfs, err := dbListSgf(40, 0)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	if err = userSgfEditHtml().Execute(w, userSgfEditData(s.User, action, msg, sgf, sgfs), defFuncMap); err != nil {
		h.ServerError(w, err)
		return
	}

}
func userSgfEditHtml() *Html {
	return defHtmlLayout().Append(
		defHtmlHead(),
		defHtmlHeader(),
		defHtmlFooter(),
		newHtmlContent("usersgfedit"),
	)
}

func userSgfEditData(u *U, action, msg string, sgf *Sgf, sgfs []Sgf) *Data {
	data := defData()
	data.User = u
	data.Header.Navs = userNavItems()
	data.Content["Editor"] = Editor{action, msg}
	data.Content["Sgf"] = sgf
	data.Content["Sgfs"] = sgfs
	return data
}

func getSgfFromRequest(r *http.Request) *Sgf {
	var s Sgf
	s.Id = atoi64(r.FormValue("id"))
	s.Time, _ = ParseDate(r.FormValue("time"))
	s.Event = r.FormValue("event")
	s.Place = r.FormValue("place")
	s.Black = r.FormValue("black")
	s.White = r.FormValue("white")
	s.Rule = r.FormValue("rule")
	s.Result = r.FormValue("result")
	s.Steps = r.FormValue("steps")
	return &s
}

func handlerUserSgfAdd(w http.ResponseWriter, r *http.Request, p []string) {
	if getSession(r) == nil {
		h.SeeOther(w, r, "/login")
		return
	}

	r.ParseForm()
	s := getSgfFromRequest(r)
	if s.Steps == "" {
		h.SeeOther(w, r, fmt.Sprint("/user/sgf/?editormsg=棋谱不能为空"))
		return
	}

	id, err := dbAddSgf(s)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	h.SeeOther(w, r, fmt.Sprint("/user/sgf/", id, "?editormsg=添加成功"))
}

func handlerUserSgfUpdate(w http.ResponseWriter, r *http.Request, p []string) {
	if getSession(r) == nil {
		h.SeeOther(w, r, "/login")
		return
	}

	r.ParseForm()
	s := getSgfFromRequest(r)
	if s.Id <= 0 {
		h.NotFound(w, "sgf id less than 0")
		return
	}
	if s.Steps == "" {
		h.SeeOther(w, r, fmt.Sprint("/user/sgf/?editormsg=棋谱不能为空"))
		return
	}

	_, err := dbUpdateSgf(s)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	h.SeeOther(w, r, fmt.Sprint("/user/sgf/", s.Id, "?editormsg=修改成功"))
}

func handlerUserSgfDelete(w http.ResponseWriter, r *http.Request, p []string) {
	if getSession(r) == nil {
		h.SeeOther(w, r, "/login")
		return
	}

	r.ParseForm()
	strid := r.FormValue("id")
	id := atoi64(strid)
	if id <= 0 {
		h.NotFound(w, "sgf id less than 0")
		return
	}

	_, err := dbDelSgf(id)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	h.SeeOther(w, r, fmt.Sprint("/user/sgf/", "?editormsg=删除成功"))
}
