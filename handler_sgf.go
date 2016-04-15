package weiqi

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

//sgf list
func sgfListHandler(h *Http) {
	var u *U
	s := getSession(h.R)
	if s != nil {
		u = s.User
	}

	sgfs, err := dbListSgf(40, 0)
	if err == sql.ErrNoRows {
		h.RequestError("page not found").NotFound()
		return
	}
	if err != nil {
		h.ServerError(err.Error())
		return
	}
	if err = sgfListHtml().Execute(h.W, sgfListData(u, sgfs), defFuncMap); err != nil {
		h.ServerError(err.Error())
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
func sgfIdHandler(h *Http) {
	var u *U
	s := getSession(h.R)
	if s != nil {
		u = s.User
	}

	id := atoi64(h.P[0])
	if id <= 0 {
		h.RequestError("sgf no found").NotFound()
		return
	}
	sgf, err := dbGetSgf(id)
	if err == sql.ErrNoRows {
		h.RequestError("找不到棋谱").NotFound()
		return
	}
	if err != nil {
		h.ServerError(err.Error())
		return
	}
	if err = sgfIdHtml().Execute(h.W, sgfIdData(u, sgf), defFuncMap); err != nil {
		h.ServerError(err.Error())
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
func userSgfEditHandler(h *Http) {
	s := getSession(h.R)
	if s == nil {
		h.SeeOther("/login")
		return
	}

	h.R.ParseForm()
	var (
		action      = "/user/sgf/add"
		msg         = h.R.FormValue("editormsg")
		sgf    *Sgf = nil
		err    error
	)

	if len(h.P) > 0 {
		action = "/user/sgf/update"
		sgf, err = dbGetSgf(atoi64(h.P[0]))
		if err == sql.ErrNoRows || err == ErrPrimaryKey {
			h.RequestError("找不到棋谱").NotFound()
			return
		}
		if err != nil {
			h.ServerError(err.Error())
			return
		}
	} else {
		sgf = new(Sgf)
	}

	sgfs, err := dbListSgf(40, 0)
	if err != nil {
		h.ServerError(err.Error())
		return
	}

	if err = userSgfEditHtml().Execute(h.W, userSgfEditData(s.User, action, msg, sgf, sgfs), defFuncMap); err != nil {
		h.ServerError(err.Error())
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
	var err error
	s.Id = atoi64(r.FormValue("id"))
	s.Time, err = time.Parse("2006-01-02", r.FormValue("time"))
	if err != nil {
		s.Time, _ = time.Parse("2006-01-02", "0000-00-00")
	}
	s.Event = r.FormValue("event")
	s.Place = r.FormValue("place")
	s.Black = r.FormValue("black")
	s.White = r.FormValue("white")
	s.Rule = r.FormValue("rule")
	s.Result = r.FormValue("result")
	s.Steps = r.FormValue("steps")
	return &s
}

func handlerUserSgfAdd(h *Http) {
	if getSession(h.R) == nil {
		h.SeeOther("/login")
		return
	}

	h.R.ParseForm()
	s := getSgfFromRequest(h.R)
	if s.Steps == "" {
		h.SeeOther(fmt.Sprint("/user/sgf/?editormsg=棋谱不能为空"))
		return
	}

	id, err := dbAddSgf(s)
	if err != nil {
		h.ServerError(err.Error())
		return
	}
	h.SeeOther(fmt.Sprint("/user/sgf/", id, "?editormsg=添加成功"))
}

func handlerUserSgfUpdate(h *Http) {
	if getSession(h.R) == nil {
		h.SeeOther("/login")
		return
	}

	h.R.ParseForm()
	s := getSgfFromRequest(h.R)
	if s.Id <= 0 {
		h.RequestError("参数错误").NotFound()
		return
	}
	if s.Steps == "" {
		h.SeeOther(fmt.Sprint("/user/sgf/?editormsg=棋谱不能为空"))
		return
	}

	_, err := dbUpdateSgf(s)
	if err != nil {
		h.ServerError(err.Error())
		return
	}
	h.SeeOther(fmt.Sprint("/user/sgf/", s.Id, "?editormsg=修改成功"))
}

func handlerUserSgfDelete(h *Http) {
	if getSession(h.R) == nil {
		h.SeeOther("/login")
		return
	}

	h.R.ParseForm()
	strid := h.R.FormValue("id")
	id := atoi64(strid)
	if id <= 0 {
		h.RequestError("参数错误").NotFound()
		return
	}

	_, err := dbDelSgf(id)
	if err != nil {
		h.ServerError(err.Error())
		return
	}
	h.SeeOther(fmt.Sprint("/user/sgf/", "?editormsg=删除成功"))
}
