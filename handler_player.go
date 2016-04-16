package weiqi

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

//player list
func playerListHandler(h *Http) {
	var u *U
	s := getSession(h.R)
	if s != nil {
		u = s.User
	}

	ps, err := dbListPlayer(40, 0)
	if err != nil && err != sql.ErrNoRows {
		h.ServerError(err.Error())
		return
	}

	err = playerListHtml().Execute(h.W, playerListData(u, ps), defFuncMap)
	if err != nil {
		h.ServerError(err.Error())
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
func playerIdHandler(h *Http) {
	var u *U
	s := getSession(h.R)
	if s != nil {
		u = s.User
	}

	id := atoi64(h.P[0])
	if id <= 0 {
		h.RequestError("参数错误").NotFound()
		return
	}

	p, err := dbGetPlayer(id)
	if err != nil {
		h.ServerError(err.Error())
		return
	}

	err = playerIdHtml().Execute(h.W, playerIdData(u, p), defFuncMap)
	if err != nil {
		h.ServerError(err.Error())
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
func userPlayerEditHandler(h *Http) {
	var u *U
	s := getSession(h.R)
	if s != nil {
		u = s.User
	} else {
		h.SeeOther("/login")
		return
	}
	h.R.ParseForm()
	var (
		action         = "/user/player/add"
		msg            = h.R.FormValue("editormsg")
		player *Player = nil
		err    error
	)

	if len(h.P) > 0 {
		action = "/user/player/update"
		player, err = dbGetPlayer(atoi64(h.P[0]))
		if err == sql.ErrNoRows || err == ErrPrimaryKey {
			h.RequestError("找不到棋手").NotFound()
			return
		}
		if err != nil {
			h.ServerError(err.Error())
			return
		}
	} else {
		player = new(Player)
	}

	ps, err := dbListPlayer(40, 0)
	if err != nil {
		h.ServerError(err.Error())
		return
	}

	err = userPlayerEditHtml().Execute(h.W, userPlayerEditData(u, action, msg, player, ps), defFuncMap)
	if err != nil {
		h.ServerError(err.Error())
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
	var err error
	p.Id = atoi64(r.FormValue("id"))
	p.Name = r.FormValue("name")
	p.Sex = parseSex(r.FormValue("sex"))
	p.Country = r.FormValue("country")
	p.Rank = r.FormValue("rank")
	p.Birth, err = time.Parse("2006-01-02", r.FormValue("birth"))
	if err != nil {
		p.Birth, err = time.Parse("2006年01月02日", r.FormValue("birth"))
		if err != nil {
			p.Birth, err = time.Parse("2006年1月2日", r.FormValue("birth"))
			if err != nil {
				p.Birth = time.Time{}
			}
		}
	}
	return &p
}

func checkPlayerFieldInput(p *Player) error {
	if p.Name == "" {
		return ErrInputEmpty
	}
	return nil
}

//player post
func handlerUserPlayerAdd(h *Http) {

	if getSession(h.R) == nil {
		h.SeeOther("/login")
		return
	}

	h.R.ParseForm()
	p := getPlayerFromRequest(h.R)

	if checkPlayerFieldInput(p) != nil {
		h.SeeOther("/user/player/?editormsg=输入不能为空")
		return
	}

	id, err := dbAddPlayer(p)
	if err == ErrPrimaryKey {
		h.RequestError("找不到棋手").NotFound()
		return
	}
	if err != nil {
		h.ServerError(err.Error())
		return
	}
	h.SeeOther(fmt.Sprintf("/user/player/%d?editormsg=提交成功", id))
}

func handlerUserPlayerDel(h *Http) {

	if getSession(h.R) == nil {
		h.SeeOther("/login")
		return
	}

	h.R.ParseForm()
	id := atoi64(h.R.FormValue("id"))
	if id <= 0 {
		h.RequestError("找不到棋手").NotFound()
		return
	}
	_, err := dbDeletePlayer(id)
	if err != nil {
		h.ServerError(err.Error())
		return
	}
	h.SeeOther(fmt.Sprint("/user/player/?editormsg=删除成功"))
}

func handlerUserPlayerUpdate(h *Http) {

	if getSession(h.R) == nil {
		h.SeeOther("/login")
		return
	}

	h.R.ParseForm()
	p := getPlayerFromRequest(h.R)

	if checkPlayerFieldInput(p) != nil {
		h.SeeOther(fmt.Sprintf("/user/player/%d?editormsg=输入不能为空", p.Id))
		return
	}

	_, err := dbUpdatePlayer(p)
	if err == ErrPrimaryKey {
		h.RequestError("找不到棋手").NotFound()
		return
	}
	if err != nil {
		h.ServerError(err.Error())
		return
	}
	h.SeeOther(fmt.Sprintf("/user/player/%d?editormsg=修改成功", p.Id))
}
