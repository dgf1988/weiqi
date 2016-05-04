package weiqi

import (
	"database/sql"
	"fmt"
	"github.com/dgf1988/weiqi/h"
	"net/http"
)

//sgf list
func handleSgfList(w http.ResponseWriter, r *http.Request, p []string) {

	err := renderSgfList(w, getSessionUser(r))
	if err != nil {
		h.ServerError(w, err)
		return
	}
}

func renderSgfList(w http.ResponseWriter, u *User) error {
	var sgfs []Sgf
	var err error
	if sgfs, err = listSgfOrderByTimeDesc(40, 0); err != nil {
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
	).Execute(w, data, nil)
}

//sgf id
func handleSgfId(w http.ResponseWriter, r *http.Request, p []string) {

	var sgf = new(Sgf)
	var err error
	if err = Sgfs.Get(atoi(p[0])).Struct(sgf); err == sql.ErrNoRows {
		h.NotFound(w, "找不到棋谱")
		return
	} else if err != nil {
		h.ServerError(w, err)
		return
	}

	var black = new(Player)
	var white = new(Player)
	if err = Db.Player.Get(nil, sgf.Black).Struct(black); err == sql.ErrNoRows {
		black = nil
	} else if err != nil {
		h.ServerError(w, err)
		return
	}
	if err = Db.Player.Get(nil, sgf.White).Struct(white); err == sql.ErrNoRows {
		white = nil
	} else if err != nil {
		h.ServerError(w, err)
		return
	}

	err = defHtmlLayout().Append(
		defHtmlHead(),
		defHtmlHeader(),
		defHtmlFooter(),
		newHtmlContent("sgfid"),
	).Execute(w, sgfIDDAta(getSessionUser(r), sgf, black, white), nil)
	if err != nil {
		h.ServerError(w, err)
	}
}

func sgfIDDAta(u *User, sgf *Sgf, black, white *Player) *Data {
	data := defData()
	data.User = u
	data.Head.Title = fmt.Sprintf("%s - %s VS %s", sgf.Event, sgf.Black, sgf.White)
	data.Head.Desc = "围棋棋谱"
	data.Head.Keywords = []string{"围棋", "棋谱", "比赛", sgf.Black, sgf.White}
	data.Content["Sgf"] = sgf
	data.Content["Black"] = black
	data.Content["White"] = white
	return data
}

//sgf edit
func handleSgfEdit(w http.ResponseWriter, r *http.Request, p []string) {
	var user *User
	if user = getSessionUser(r); user == nil {
		h.SeeOther(w, r, "/login")
		return
	}

	r.ParseForm()
	var err error
	var sgfs = make([]Sgf, 0)
	if rows, err := Sgfs.ListDesc(40, 0); err != nil {
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
	var action string
	if len(p) > 0 {
		if err = Sgfs.Get(atoi(p[0])).Struct(sgf); err == sql.ErrNoRows {
			h.NotFound(w, "找不到棋手")
			return
		} else if err != nil {
			h.ServerError(w, err)
			return
		}
		action = "/user/sgf/update"
	} else {
		action = "/user/sgf/add"
	}

	if err = userSgfEditHtml().Execute(w, userSgfEditData(user, action, r.FormValue("editormsg"), sgf, sgfs), nil); err != nil {
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
	s.Time, _ = parseDate(r.FormValue("time"))
	s.Event = r.FormValue("event")
	s.Place = r.FormValue("place")
	s.Black = r.FormValue("black")
	s.White = r.FormValue("white")
	s.Rule = r.FormValue("rule")
	s.Result = r.FormValue("result")
	s.Steps = r.FormValue("steps")
	return &s
}

func handleSgfAdd(w http.ResponseWriter, r *http.Request, p []string) {
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

func handleSgfUpdate(w http.ResponseWriter, r *http.Request, p []string) {
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

func handleSgfDel(w http.ResponseWriter, r *http.Request, p []string) {
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

func sgf_remote_handler(w http.ResponseWriter, r *http.Request, args []string) {
	if getSessionUser(r) == nil {
		h.SeeOther(w, r, "/login")
		return
	}

	var err error
	if err = r.ParseForm(); err != nil {
		h.ServerError(w, err)
		return
	}

	var sgf *Sgf
	if sgf, err = remoteSgf(r.FormValue("src"), r.FormValue("charset")); err != nil {
		h.NotFound(w, err.Error())
		return
	}

	var id int64
	if id, err = Db.Sgf.Add(nil, sgf.Time, sgf.Place, sgf.Event, sgf.Black, sgf.White, sgf.Rule, sgf.Result, sgf.Steps); err != nil {
		h.ServerError(w, err)
		return
	}
	h.SeeOther(w, r, fmt.Sprint("/user/sgf/", id))
}
