package weiqi

import (
	"database/sql"
	"fmt"
	"github.com/dgf1988/weiqi/h"
	"net/http"
	"sort"
)

//player list
func handlePlayerList(w http.ResponseWriter, r *http.Request, args []string) {

	err := playerListRender(w, getSessionUser(r))
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
	var players = make([]Player, 0)
	if rows, err := Players.List(40, 0); err != nil {
		return err
	} else {
		defer rows.Close()
		for rows.Next() {
			var player Player
			err = rows.Struct(&player)
			if err != nil {
				return err
			} else {
				players = append(players, player)
			}
		}
		if err = rows.Err(); err != nil {
			return err
		}
	}

	data := defData()
	data.Head.Title = "棋手列表"
	data.Head.Desc = "围棋棋手列表"
	data.Head.Keywords = []string{"围棋", "棋手", "资料"}
	data.User = u
	data.Content["Players"] = players
	return defHtmlLayout().Append(
		defHtmlHead(),
		defHtmlHeader(),
		defHtmlFooter(),
		newHtmlContent("playerlist"),
	).Execute(w, data, nil)
}

//player id
func handlePlayerId(w http.ResponseWriter, r *http.Request, args []string) {

	switch err := renderPlayerid(w, getSessionUser(r), args[0]); err {
	case nil:
	case sql.ErrNoRows:
		h.NotFound(w, "棋手不存在")
	default:
		h.ServerError(w, err)
	}
}

func renderPlayerid(w http.ResponseWriter, u *User, id interface{}) error {
	var player = new(Player)
	var text = new(Text)
	var err error

	if err = Players.Get(id).Struct(player); err != nil {
		return err
	} else {
		var textid int64
		if err = TextPlayer.Get(nil, player.Id).Scan(nil, nil, &textid); err == nil {
			if err = Texts.Get(textid).Struct(text); err != nil && err != sql.ErrNoRows {
				return err
			}
		} else if err != sql.ErrNoRows {
			return err
		}
	}

	var sgfs []Sgf
	if sgfs, err = listSgfByNamesOrderTimeDesc(player.Name); err != nil {
		return err
	}


	data := defData()
	data.User = u
	data.Head.Title = player.Name
	data.Head.Desc = "围棋棋手"
	data.Head.Keywords = []string{"围棋", "棋手", "资料", player.Name}
	data.Content["Player"] = player
	text.Text = parseTextToHtml(text.Text)
	data.Content["Text"] = text
	data.Content["Sgfs"] = sgfs
	return defHtmlLayout().Append(
		defHtmlHead(),
		defHtmlHeader(),
		defHtmlFooter(),
		newHtmlContent("playerid"),
	).Execute(w, data, nil)

}

//plaeyr edit
func handlePlayerEdit(w http.ResponseWriter, r *http.Request, p []string) {
	var user *User
	if user = getSessionUser(r); user == nil {
		h.SeeOther(w, r, "/login")
		return
	}


	r.ParseForm()
	var (
		action = "/user/player/add"
		msg    = r.FormValue("editormsg")

		player  = new(Player)
		text    = new(Text)
		players = make([]Player, 0)
	)
	if len(p) > 0 {
		action = "/user/player/update"
		if err := Players.Get(p[0]).Struct(player); err == nil {
			var textid int64
			if err = TextPlayer.Get(nil, player.Id).Scan(nil, nil, &textid); err == nil {
				if err = Texts.Get(textid).Struct(text); err != nil && err != sql.ErrNoRows {
					h.ServerError(w, err)
					return
				}
			} else if err != nil && err != sql.ErrNoRows {
				h.ServerError(w, err)
				return
			}
		} else if err == sql.ErrNoRows {
			h.NotFound(w, "棋手不存在")
			return
		} else {
			h.ServerError(w, err)
			return
		}
	}

	if rows, err := Players.List(40, 0); err == nil {
		defer rows.Close()
		for rows.Next() {
			var a Player
			if err = rows.Struct(&a); err == nil {
				players = append(players, a)
			} else {
				h.ServerError(w, err)
				return
			}
		}
		if err = rows.Err(); err != nil {
			h.ServerError(w, err)
			return
		}
	} else {
		h.ServerError(w, err)
		return
	}

	err := userPlayerEditRender(w, user, action, msg, player, text, players)
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
	).Execute(w, data, nil)
}

func getPlayerFromRequest(r *http.Request) *Player {
	var p Player
	p.Id = atoi64(r.FormValue("id"))
	p.Name = r.FormValue("name")
	p.Sex = chineseToSex(r.FormValue("sex"))
	p.Country = r.FormValue("country")
	p.Rank = r.FormValue("rank")
	p.Birth, _ = parseDate(r.FormValue("birth"))
	return &p
}

//player post
func handlePlayerAdd(w http.ResponseWriter, r *http.Request, args []string) {

	if getSession(r) == nil {
		h.SeeOther(w, r, "/login")
		return
	}

	r.ParseForm()
	p := getPlayerFromRequest(r)
	text := r.FormValue("text")

	if p.Name == "" {
		h.SeeOther(w, r, "/user/player/?editormsg=名字不能为空")
		return
	}

	playerid, err := Players.Add(nil, p.Name, p.Sex, p.Country, p.Rank, p.Birth)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	if text != "" {
		textid, err := Texts.Add(nil, text)
		if err != nil {
			h.ServerError(w, err)
			return
		}
		_, err = TextPlayer.Add(nil, playerid, textid)
		if err != nil {
			h.ServerError(w, err)
			return
		}
	}

	h.SeeOther(w, r, fmt.Sprintf("/user/player/%d?editormsg=提交成功", playerid))
}

func handlePlayerDel(w http.ResponseWriter, r *http.Request, p []string) {
	var err error

	if getSession(r) == nil {
		h.SeeOther(w, r, "/login")
		return
	}

	r.ParseForm()
	playerid := atoi64(r.FormValue("id"))
	if playerid < 0 {
		h.NotFound(w, "参数错误")
		return
	}

	var playertextid int64
	var textid int64
	err = TextPlayer.Get(nil, playerid).Scan(&playertextid, nil, &textid)
	if err == nil {
		_, err = Texts.Del(textid)
		if err != nil {
			h.ServerError(w, err)
			return
		}
		_, err = TextPlayer.Del(playertextid)
		if err != nil {
			h.ServerError(w, err)
			return
		}
	} else if err != sql.ErrNoRows {
		h.ServerError(w, err)
		return
	}
	var n int64
	n, err = Players.Del(playerid)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	if n == 0 {
		h.NotFound(w, "找不到棋谱")
		return
	}
	h.SeeOther(w, r, "/user/player/?editormsg=删除成功")
}

func handlePlayerUpdate(w http.ResponseWriter, r *http.Request, args []string) {
	var err error
	if getSession(r) == nil {
		h.SeeOther(w, r, "/login")
		return
	}

	r.ParseForm()
	player := getPlayerFromRequest(r)
	text := r.FormValue("text")

	if player.Name == "" {
		h.SeeOther(w, r, fmt.Sprintf("/user/player/%d?editormsg=名字不能为空", player.Id))
		return
	}

	//更新棋手
	_, err = Players.Update(player.Id).Values(nil, player.Name, player.Sex, player.Country, player.Rank, player.Birth)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	var textid int64
	var playertextid int64
	err = TextPlayer.Get(nil, player.Id).Scan(&playertextid, nil, &textid)
	if err == nil {
		err = Texts.Get(textid).Scan()
		if err == nil {
			_, err = Texts.Update(textid).Values(nil, text)
			if err != nil {
				h.ServerError(w, err)
				return
			}
		} else if err == sql.ErrNoRows {
			_, err = Texts.Add(nil, text)
			if err != nil {
				h.ServerError(w, err)
				return
			}
		} else {
			h.ServerError(w, err)
			return
		}
	} else if err == sql.ErrNoRows {
		if text != "" {
			textid, err = Texts.Add(nil, text)
			if err == nil {
				playertextid, err = TextPlayer.Add(nil, player.Id, textid)
				if err != nil {
					h.ServerError(w, err)
					return
				}
			} else {
				h.ServerError(w, err)
				return
			}
		}
	} else {
		h.ServerError(w, err)
		return
	}
	h.SeeOther(w, r, fmt.Sprintf("/user/player/%d?editormsg=修改成功", player.Id))
}
