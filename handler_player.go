package weiqi

import (
	"database/sql"
	"fmt"
	"net/http"
	"github.com/dgf1988/weiqi/h"
)

//player list
func playerListHandler(w http.ResponseWriter, r *http.Request, p []string) {
	var u *U
	s := getSession(r)
	if s != nil {
		u = s.User
	}

	ps, err := dbListPlayer(40, 0)
	if err != nil && err != sql.ErrNoRows {
		h.ServerError(w, err)
		return
	}

	err = playerListHtml().Execute(w, playerListData(u, ps), defFuncMap)
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

func playerListData(u *U, playerlist []Player) *Data {
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
	var u *U
	s := getSession(r)
	if s != nil {
		u = s.User
	}

	id := atoi64(args[0])
	if id <= 0 {
		h.NotFound(w, "找不到棋手")
		return
	}

	p, err := dbGetPlayer(id)
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

func playerIdData(u *U, player *Player) *Data {
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
	var u *U
	s := getSession(r)
	if s != nil {
		u = s.User
	} else {
		h.SeeOther(w, r, "/login")
		return
	}
	r.ParseForm()
	var (
		action         = "/user/player/add"
		msg            = r.FormValue("editormsg")
		player *Player = nil
		err    error
	)

	if len(p) > 0 {
		action = "/user/player/update"
		player, err = dbGetPlayer(atoi64(p[0]))
		if err == sql.ErrNoRows || err == ErrPrimaryKey {
			h.NotFound(w, "找不到棋手")
			return
		}
		if err != nil {
			h.ServerError(w, err)
			return
		}
	} else {
		player = new(Player)
	}

	ps, err := dbListPlayer(40, 0)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	err = userPlayerEditHtml().Execute(w, userPlayerEditData(u, action, msg, player, ps), defFuncMap)
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

func userPlayerEditData(u *U, action, msg string, player *Player, players []Player) *Data {
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
	p.Sex = parseSex(r.FormValue("sex"))
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

	id, err := dbAddPlayer(p)
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
	_, err := dbDeletePlayer(id)
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

	_, err := dbUpdatePlayer(p)
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
