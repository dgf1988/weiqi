package weiqi

import (
	"database/sql"
	"fmt"
	"github.com/dgf1988/weiqi/h"
	"net/http"
)

//player list
func playerListHandler(w http.ResponseWriter, r *http.Request, args []string) {
	var u *User
	s := getSession(r)
	if s != nil {
		u = s.User
	}

	var ps [40]Player
	n, err := Players.ListBy(&ps, 0)
	if err != nil && err != sql.ErrNoRows {
		h.ServerError(w, err)
		return
	}

	err = playerListHtml().Execute(w, playerListData(u, ps[:n]), defFuncMap)
	if err != nil {
		h.ServerError(w, err)
		return
	}
}

func playerListHtml() *Html {
	return defHtmlLayout().Append(
		defHtmlHead(),
		defHtmlHeader(),
		defHtmlFooter(),
		newHtmlContent("playerlist"),
	)
}

func playerListData(u *User, playerlist []Player) *Data {
	data := defData()
	data.Head.Title = "棋手列表"
	data.Head.Desc = "围棋棋手列表"
	data.Head.Keywords = []string{"围棋", "棋手", "资料"}
	data.User = u
	data.Content["Players"] = playerlist
	return data
}

//player id
func playerIdHandler(w http.ResponseWriter, r *http.Request, args []string) {
	var u *User
	s := getSession(r)
	if s != nil {
		u = s.User
	}

	id := atoi64(args[0])
	if id <= 0 {
		h.NotFound(w, "找不到棋手")
		return
	}

	var p = new(Player)
	err := Players.GetBy(id, p)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	err = playerIdHtml().Execute(w, playerIdData(u, p), defFuncMap)
	if err != nil {
		h.ServerError(w, err)
		return
	}
}

func playerIdHtml() *Html {
	return defHtmlLayout().Append(
		defHtmlHead(),
		defHtmlHeader(),
		defHtmlFooter(),
		newHtmlContent("playerid"),
	)
}

func playerIdData(u *User, player *Player) *Data {
	data := defData()
	data.User = u
	data.Head.Title = player.Name
	data.Head.Desc = "围棋棋手"
	data.Head.Keywords = []string{"围棋", "棋手", "资料", player.Name}
	data.Content["Player"] = player
	return data
}

//plaeyr edit
func userPlayerEditHandler(w http.ResponseWriter, r *http.Request, p []string) {
	var u *User
	s := getSession(r)
	if s != nil {
		u = s.User
	} else {
		h.SeeOther(w, r, "/login")
		return
	}
	r.ParseForm()
	var (
		action = "/user/player/add"
		msg    = r.FormValue("editormsg")
		player = new(Player)
		err    error
	)

	if len(p) > 0 {
		action = "/user/player/update"
		err = Players.GetBy(atoi64(p[0]), player)
		if err == sql.ErrNoRows || err == ErrPrimaryKey {
			h.NotFound(w, "找不到棋手")
			return
		}
		if err != nil {
			h.ServerError(w, err)
			return
		}
	}

	var ps [40]Player
	n, err := Players.ListBy(&ps, 0)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	err = userPlayerEditHtml().Execute(w, userPlayerEditData(u, action, msg, player, ps[:n]), defFuncMap)
	if err != nil {
		h.ServerError(w, err)
		return
	}
}

func userPlayerEditHtml() *Html {
	return defHtmlLayout().Append(
		defHtmlHead(),
		defHtmlHeader(),
		defHtmlFooter(),
		newHtmlContent("userplayeredit"),
	)
}

func userPlayerEditData(u *User, action, msg string, player *Player, players []Player) *Data {
	data := defData()
	data.User = u
	data.Header.Navs = userNavItems()
	data.Content["Editor"] = Editor{action, msg}
	data.Content["Player"] = player
	data.Content["Players"] = players
	return data
}

func getPlayerFromRequest(r *http.Request) *Player {
	var p Player
	p.Id = atoi64(r.FormValue("id"))
	p.Name = r.FormValue("name")
	p.Sex = chineseToSex(r.FormValue("sex"))
	p.Country = r.FormValue("country")
	p.Rank = r.FormValue("rank")
	p.Birth, _ = ParseDate(r.FormValue("birth"))
	return &p
}

func checkPlayerFieldInput(p *Player) error {
	if p.Name == "" {
		return ErrInputEmpty
	}
	return nil
}

//player post
func handlerUserPlayerAdd(w http.ResponseWriter, r *http.Request, args []string) {

	if getSession(r) == nil {
		h.SeeOther(w, r, "/login")
		return
	}

	r.ParseForm()
	p := getPlayerFromRequest(r)

	if checkPlayerFieldInput(p) != nil {
		h.SeeOther(w, r, "/user/player/?editormsg=输入不能为空")
		return
	}

	id, err := Players.Add(nil, p.Name, p.Sex, p.Country, p.Rank, p.Birth)
	if err == ErrPrimaryKey {
		h.NotFound(w, "找不到棋手")
		return
	}
	if err != nil {
		h.ServerError(w, err)
		return
	}
	h.SeeOther(w, r, fmt.Sprintf("/user/player/%d?editormsg=提交成功", id))
}

func handlerUserPlayerDel(w http.ResponseWriter, r *http.Request, p []string) {

	if getSession(r) == nil {
		h.SeeOther(w, r, "/login")
		return
	}

	r.ParseForm()
	id := atoi64(r.FormValue("id"))
	if id <= 0 {
		h.NotFound(w, "找不到棋手")
		return
	}
	_, err := Players.Del(id)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	h.SeeOther(w, r, fmt.Sprint("/user/player/?editormsg=删除成功"))
}

func handlerUserPlayerUpdate(w http.ResponseWriter, r *http.Request, args []string) {

	if getSession(r) == nil {
		h.SeeOther(w, r, "/login")
		return
	}

	r.ParseForm()
	p := getPlayerFromRequest(r)

	if checkPlayerFieldInput(p) != nil {
		h.SeeOther(w, r, fmt.Sprintf("/user/player/%d?editormsg=输入不能为空", p.Id))
		return
	}

	_, err := Players.Set(p.Id, nil, p.Name, p.Sex, p.Country, p.Rank, p.Birth)
	if err == ErrPrimaryKey {
		h.NotFound(w, "找不到棋手")
		return
	}
	if err != nil {
		h.ServerError(w, err)
		return
	}
	h.SeeOther(w, r, fmt.Sprintf("/user/player/%d?editormsg=修改成功", p.Id))
}
