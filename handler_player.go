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

	err := playerListRender(w, u)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			h.NotFound(w, "找不到棋手")
		default:
			h.ServerError(w, err)
		}
	}
}

func playerListRender(w http.ResponseWriter, u *User) error {
	var playerlist = [40]Player{}
	n, err := Players.ListBy(&playerlist, 0)
	if err != nil {
		return err
	}
	data := defData()
	data.Head.Title = "棋手列表"
	data.Head.Desc = "围棋棋手列表"
	data.Head.Keywords = []string{"围棋", "棋手", "资料"}
	data.User = u
	data.Content["Players"] = playerlist[:n]
	return defHtmlLayout().Append(
		defHtmlHead(),
		defHtmlHeader(),
		defHtmlFooter(),
		newHtmlContent("playerlist"),
	).Execute(w, data, defFuncMap)
}

//player id
func playerIdHandler(w http.ResponseWriter, r *http.Request, args []string) {
	var u *User
	s := getSession(r)
	if s != nil {
		u = s.User
	}

	err := playerIdRender(w, u, args[0])
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			h.NotFound(w, "棋手不存在")
		default:
			h.ServerError(w, err)
		}
	}
}

func playerIdRender(w http.ResponseWriter, u *User, id interface{}) error {
	var player = new(Player)
	var text = new(Text)

	err := Players.GetStruct(id, player)
	if err != nil {
		return err
	}
	var textid int64
	err = PlayerText.Find(nil, player.Id).Scan(nil, nil, &textid)
	if err == sql.ErrNoRows {
		goto DATA
	}
	if err != nil {
		return err
	}
	err = Texts.GetStruct(textid, text)
	if err == sql.ErrNoRows {
		goto DATA
	}
	if err != nil {
		return err
	}

	DATA:
	data := defData()
	data.User = u
	data.Head.Title = player.Name
	data.Head.Desc = "围棋棋手"
	data.Head.Keywords = []string{"围棋", "棋手", "资料", player.Name}
	data.Content["Player"] = player
	text.Text = parseTextToHtml(text.Text)
	data.Content["Text"] = text
	return defHtmlLayout().Append(
		defHtmlHead(),
		defHtmlHeader(),
		defHtmlFooter(),
		newHtmlContent("playerid"),
	).Execute(w, data, defFuncMap)

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
		playerlist = [40]Player{}
		text = new(Text)
		err    error
	)
	if len(p) > 0 {
		action = "/user/player/update"
		err = Players.GetStruct(p[0], player)
		if err == sql.ErrNoRows {
			h.NotFound(w, "棋手不存在")
			return
		}
		if err != nil {
			h.ServerError(w, err)
			return
		}

		var textid int64
		err = PlayerText.Find(nil, player.Id).Scan(nil, nil, &textid)
		if err == sql.ErrNoRows {
			goto LISTPLAYERS
		}
		if err != nil {
			h.ServerError(w, err)
			return
		}

		err = Texts.GetStruct(textid, text)
		if err != nil && err != sql.ErrNoRows {
			h.ServerError(w, err)
			return
		}
	}

	LISTPLAYERS:
	n, err := Players.ListBy(&playerlist, 0)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	err = userPlayerEditRender(w, u, action, msg, player, text, playerlist[:n])
	if err != nil {
		h.ServerError(w, err)
	}
}

func userPlayerEditRender(w http.ResponseWriter, u *User, action, msg string, player *Player, text *Text, playerlist []Player) error {
	var editor = Editor{action, msg}
	data := defData()
	data.User = u
	data.Header.Navs = userNavItems()
	data.Content["Editor"] = editor
	data.Content["Player"] = player
	data.Content["Text"] = text
	data.Content["Players"] = playerlist

	return defHtmlLayout().Append(
		defHtmlHead(),
		defHtmlHeader(),
		defHtmlFooter(),
		newHtmlContent("userplayeredit"),
	).Execute(w, data, defFuncMap)
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
	text := r.FormValue("text")


	if checkPlayerFieldInput(p) != nil {
		h.SeeOther(w, r, "/user/player/?editormsg=输入不能为空")
		return
	}

	playerid, err := Players.Add(nil, p.Name, p.Sex, p.Country, p.Rank, p.Birth)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	textid, err := Texts.Add(nil, text)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	_, err = PlayerText.Add(nil, playerid, textid)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	h.SeeOther(w, r, fmt.Sprintf("/user/player/%d?editormsg=提交成功", playerid))
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
	var textid int64
	var pt_id int64
	err := PlayerText.Find(nil, id).Scan(&pt_id, nil, &textid)
	if err == sql.ErrNoRows {
		goto DELETEPLAYER
	}
	if err != nil  {
		h.ServerError(w, err)
		return

	}
	Texts.Del(textid)
	PlayerText.Del(pt_id)

	DELETEPLAYER:
	_, err = Players.Del(id)
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
	text := r.FormValue("text")

	if checkPlayerFieldInput(p) != nil {
		h.SeeOther(w, r, fmt.Sprintf("/user/player/%d?editormsg=输入不能为空", p.Id))
		return
	}

	_, err := Players.Set(p.Id, nil, p.Name, p.Sex, p.Country, p.Rank, p.Birth)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	var textid int64
	err = PlayerText.Find(nil, p.Id).Scan(nil, nil, &textid)
	if err == sql.ErrNoRows {
		textid, err = Texts.Add(nil, text)
		if err != nil {
			h.ServerError(w, err)
			return
		}
		_, err = PlayerText.Add(nil, p.Id, textid)
		if err != nil {
			h.ServerError(w, err)
			return
		}
	}
	if err != nil {
		h.ServerError(w, err)
		return
	} else {
		_, err = Texts.Set(textid, nil, text)
		if err != nil {
			h.ServerError(w, err)
			return
		}
	}
	h.SeeOther(w, r, fmt.Sprintf("/user/player/%d?editormsg=修改成功", p.Id))
}
