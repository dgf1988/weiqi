package weiqi

import (
	"database/sql"
	"fmt"
	"github.com/dgf1988/weiqi/h"
	"net/http"
)

//sgf list
func sgfListHandler(w http.ResponseWriter, r *http.Request, p []string) {
	var u *User
	s := getSession(r)
	if s != nil {
		u = s.User
	}

	err := renderSgfList(w, u)
	if err != nil {
		h.ServerError(w, err)
		return
	}
}

func renderSgfList(w http.ResponseWriter, u *User) error {
	var sgfs = make([]Sgf, 0)

	if rows, err := Sgfs.List(40, 0); err == nil {
		defer rows.Close()
		for rows.Next() {
			var sgf Sgf
			if err = rows.Struct(&sgf); err == nil {
				sgfs = append(sgfs, sgf)
			} else {
				return err
			}
		}
		if err = rows.Err(); err != nil {
			return err
		}
	} else {
		return err
	}

	data := defData()
	data.User = u
	data.Head.Title = "棋谱列表"
	data.Head.Desc = "围棋棋谱列表"
	data.Head.Keywords = []string{"围棋", "棋谱", "比赛"}
	data.Content["Sgfs"] = sgfs

	return defHtmlLayout().Append(
		defHtmlHead(),
		defHtmlHeader(),
		defHtmlFooter(),
		newHtmlContent("sgflist"),
	).Execute(w, data, defFuncMap)
}

//sgf id
func sgfIdHandler(w http.ResponseWriter, r *http.Request, p []string) {
	var user *User
	if s := getSession(r); s != nil {
		user = s.User
	}

	var sgfid = atoi(p[0])
	var sgf = new(Sgf)
	var err error
	if err = Sgfs.Get(sgfid).Struct(sgf); err == nil {
		if err = sgfIdHtml().Execute(w, sgfIdData(user, sgf), defFuncMap); err != nil {
			h.ServerError(w, err)
		}
	} else if err == sql.ErrNoRows {
		h.NotFound(w, "找不到棋谱")
	} else {
		h.ServerError(w, err)
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

func sgfIdData(u *User, sgf *Sgf) *Data {
	data := defData()
	data.User = u
	data.Head.Title = fmt.Sprintf("%s - %s VS %s", sgf.Event, sgf.Black, sgf.White)
	data.Head.Desc = "围棋棋谱"
	data.Head.Keywords = []string{"围棋", "棋谱", "比赛", sgf.Black, sgf.White}
	data.Content["Sgf"] = sgf
	return data
}

//sgf edit
func userSgfEditHandler(w http.ResponseWriter, r *http.Request, p []string) {
	var s = getSession(r)
	if s == nil {
		h.SeeOther(w, r, "/login")
		return
	}

	r.ParseForm()
	var err error
	var sgfs = make([]Sgf, 0)
	if rows, err := Sgfs.List(40, 0); err != nil {
		h.ServerError(w, err)
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			var sgf Sgf
			if err = rows.Struct(&sgf); err == nil {
				sgfs = append(sgfs, sgf)
			} else {
				h.ServerError(w, err)
				return
			}
		}
	}

	var sgf = new(Sgf)
	var action = "/user/sgf/add"
	if len(p) > 0 {
		sgfid := atoi(p[0])
		if err = Sgfs.Get(sgfid).Struct(sgf); err == sql.ErrNoRows {
			h.NotFound(w, "找不到棋手")
			return
		} else if err != nil  {
			h.ServerError(w, err)
			return
		}
		action = "/user/sgf/update"
	}

	if err = userSgfEditHtml().Execute(w, userSgfEditData(s.User, action, r.FormValue("editormsg"), sgf, sgfs), defFuncMap); err != nil {
		h.ServerError(w, err)
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

func userSgfEditData(u *User, action, msg string, sgf *Sgf, sgfs []Sgf) *Data {
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

	id, err := Sgfs.Add(nil, s.Time, s.Place, s.Event, s.Black, s.White, s.Rule, s.Result, s.Steps)
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

	_, err := Sgfs.Update(s.Id).Values(nil, s.Time, s.Place, s.Event, s.Black, s.White, s.Rule, s.Result, s.Steps)
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

	_, err := Sgfs.Del(id)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	h.SeeOther(w, r, fmt.Sprint("/user/sgf/", "?editormsg=删除成功"))
}
