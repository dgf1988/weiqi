package weiqi

import (
	"database/sql"
	"fmt"
	"github.com/dgf1988/weiqi/h"
	"net/http"
)

//player list
func player_list_handler(w http.ResponseWriter, r *http.Request, args []string) {
	var err error
	var data = defData()
	data.User = getSessionUser(r)
	data.Head.Title = "棋手列表"
	data.Head.Desc = "围棋棋手列表"
	data.Head.Keywords = []string{"围棋", "棋手", "资料"}

	var players []Player
	if players, err = listPlayerOrderByRankDesc(40, 0); err != nil {
		h.ServerError(w, err)
		return
	}
	var cn = make([]Player, 0)
	var kr = make([]Player, 0)
	var jp = make([]Player, 0)
	var other = make([]Player, 0)
	for _, player := range players {
		switch player.Country {
		case "中国":
			cn = append(cn, player)
		case "日本":
			jp = append(jp, player)
		case "韩国":
			kr = append(kr, player)
		default:
			other = append(other, player)
		}
	}
	data.Content["Cn"] = cn
	data.Content["Jp"] = jp
	data.Content["Kr"] = kr
	data.Content["Other"] = other

	var html = defHtmlLayout().Append(defHtmlHead(), defHtmlHeader(), defHtmlFooter(), newHtmlContent("playerlist"))
	if err = html.Execute(w, data, nil); err != nil {
		logError("%s %s html.execute %s", r.Method, r.URL, err.Error())
		return
	}
}

//player id
func player_info_handler(w http.ResponseWriter, r *http.Request, args []string) {
	var err error
	var data = defData()
	data.User = getSessionUser(r)

	var player = new(Player)
	if err = Db.Player.Get(atoi(args[0])).Struct(player); err == sql.ErrNoRows {
		h.NotFound(w, "找不到棋手")
		return
	} else if err != nil {
		h.ServerError(w, err)
		return
	}
	data.Head.Title = player.Name
	data.Head.Desc = "围棋棋手"
	data.Head.Keywords = []string{"围棋", "棋手", "资料", player.Name}
	data.Content["Player"] = player

	var textid int64
	var text = new(Text)
	if err = Db.TextPlayer.Get(nil, player.Id).Scan(nil, nil, &textid); err == nil {
		if err = Db.Text.Get(textid).Struct(text); err != nil && err != sql.ErrNoRows {
			h.ServerError(w, err)
			return
		} else {
			text.Text = parseTextToHtml(text.Text)
		}
	} else if err != sql.ErrNoRows {
		h.ServerError(w, err)
		return
	}
	data.Content["Text"] = text

	var img Img
	if err = Db.Img.Get(nil, player.Name).Struct(&img); err == nil {
		data.Content["Img"] = img
	} else if err != sql.ErrNoRows {
		h.ServerError(w, err)
		return
	}

	var sgfs []Sgf
	if sgfs, err = listSgfByNameOrderByTimeDesc(player.Name); err != nil {
		h.ServerError(w, err)
		return
	}
	data.Content["Sgfs"] = sgfs

	var html = defHtmlLayout().Append(defHtmlHead(), defHtmlHeader(), defHtmlFooter(), newHtmlContent("playerid"))
	if err = html.Execute(w, data, nil); err != nil {
		logError("%s %s html.execute %s", r.Method, r.URL, err.Error())
		return
	}
}

func player_manage_handler(w http.ResponseWriter, r *http.Request, args []string) {
	var user = getSessionUser(r)
	if user == nil {
		h.SeeOther(w, r, "/login")
		return
	}

	var data = defData()
	data.User = user
	data.Head.Title = "棋手管理"
	data.Header.Navs = userNavItems()

	var players = make([]Player, 0)
	if rows, err := Db.Player.ListDesc(100, 0); err != nil {
		h.ServerError(w, err)
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			var player Player
			if err = rows.Struct(&player); err != nil {
				h.ServerError(w, err)
				return
			} else {
				players = append(players, player)
			}
		}
		if err = rows.Err(); err != nil {
			h.ServerError(w, err)
			return
		}
		var cn = make([]Player, 0)
		var kr = make([]Player, 0)
		var jp = make([]Player, 0)
		var other = make([]Player, 0)
		for _, player := range players {
			switch player.Country {
			case "中国":
				cn = append(cn, player)
			case "日本":
				jp = append(jp, player)
			case "韩国":
				kr = append(kr, player)
			default:
				other = append(other, player)
			}
		}
		data.Content["Cn"] = cn
		data.Content["Jp"] = jp
		data.Content["Kr"] = kr
		data.Content["Other"] = other
		var html = defHtmlLayout().Append(defHtmlHead(), defHtmlHeader(), defHtmlFooter(), newHtmlContent("userplayer"))
		if err = html.Execute(w, data, nil); err != nil {
			logError(err.Error())
			return
		}
	}
}

func player_editor_handler(w http.ResponseWriter, r *http.Request, args []string) {
	var user = getSessionUser(r)
	if user == nil {
		h.SeeOther(w, r, "/login")
		return
	}

	var playerid = atoi64(args[0])
	var player Player
	var err error
	if err = Db.Player.Get(playerid).Struct(&player); err == sql.ErrNoRows {
		h.NotFound(w, "找不到棋手")
		return
	} else if err != nil {
		h.ServerError(w, err)
		return
	}

	if r.Method == POST {
		//提交修改
		if err = r.ParseForm(); err != nil {
			h.ServerError(w, err)
			return
		}

		//获取资料
		var name = r.FormValue("name")
		if name == "" {
			h.NotFound(w, "姓名不能为空")
			return
		}
		var sex = atoi64(r.FormValue("sex"))
		var country = r.FormValue("country")
		var rank = atoi64(r.FormValue("rank"))
		var birth, _ = parseDate(r.FormValue("birth"))
		var text = r.FormValue("text")

		//更新棋手资料
		if _, err = Db.Player.Update(playerid).Values(nil, name, sex, country, rank, birth); err != nil {
			h.ServerError(w, err)
			return
		}

		var textid int64
		var textplayerid int64
		//是否有个人介绍
		if err = Db.TextPlayer.Get(nil, playerid).Scan(&textplayerid, nil, &textid); err == sql.ErrNoRows {
			//没有
			if text != "" {
				if textid, err = Db.Text.Add(nil, text); err != nil {
					h.ServerError(w, err)
					return
				}
				if _, err = Db.TextPlayer.Add(nil, playerid, textid); err != nil {
					h.ServerError(w, err)
					return
				}
			}
		} else if err == nil {
			//有，更新
			if text == "" {
				if _, err = Db.Text.Del(textid); err != nil {
					h.ServerError(w, err)
					return
				}
				if _, err = Db.TextPlayer.Del(textplayerid); err != nil {
					h.ServerError(w, err)
					return
				}
			} else if _, err = Db.Text.Update(textid).Values(nil, text); err != nil {
				h.ServerError(w, err)
				return
			}
		} else {
			h.ServerError(w, err)
			return
		}
		h.SeeOther(w, r, fmt.Sprint("/user/player/", playerid))
		return
	}

	var data = defData()
	data.User = user
	data.Head.Title = fmt.Sprintf("棋手编辑 - %s", player.Name)
	data.Header.Navs = userNavItems()
	data.Content["Player"] = player

	var textid int64
	var text = new(Text)
	if err = Db.TextPlayer.Get(nil, player.Id).Scan(nil, nil, &textid); err == sql.ErrNoRows {

	} else if err == nil {
		if err = Db.Text.Get(textid).Struct(text); err == sql.ErrNoRows {

		} else if err == nil {

		} else {
			h.ServerError(w, err)
			return
		}
	} else {
		h.ServerError(w, err)
		return
	}
	data.Content["Text"] = text

	if err = defHtmlLayout().Append(defHtmlHead(), defHtmlHeader(), defHtmlFooter(), newHtmlContent("userplayeredit")).Execute(w, data, nil); err != nil {
		logError(err.Error())
		return
	}
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
		if err := Db.Player.Get(p[0]).Struct(player); err == nil {
			var textid int64
			if err = Db.TextPlayer.Get(nil, player.Id).Scan(nil, nil, &textid); err == nil {
				if err = Db.Text.Get(textid).Struct(text); err != nil && err != sql.ErrNoRows {
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

	if rows, err := Db.Player.ListDesc(40, 0); err == nil {
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
	p.Rank = chineseToRank(r.FormValue("rank"))
	p.Birth, _ = parseDate(r.FormValue("birth"))
	return &p
}

//player add
func player_add_handler(w http.ResponseWriter, r *http.Request, args []string) {
	if getSessionUser(r) == nil {
		h.SeeOther(w, r, "/login")
		return
	}

	var err error
	if err = r.ParseForm(); err != nil {
		h.ServerError(w, err)
		return
	}

	var name = r.FormValue("name")
	if name == "" {
		h.NotFound(w, "姓名不能为空")
		return
	}

	var id int64
	if err = Db.Player.Get(nil, name).Scan(&id); err == nil {

	} else if err != sql.ErrNoRows {
		h.ServerError(w, err)
		return
	} else {
		var sex = atoi64(r.FormValue("sex"))
		var country = r.FormValue("country")
		var rank = atoi64(r.FormValue("rank"))
		var birth, _ = parseDate(r.FormValue("birth"))
		if id, err = Db.Player.Add(nil, name, sex, country, rank, birth); err != nil {
			h.ServerError(w, err)
			return
		} else {
		}
	}
	h.SeeOther(w, r, fmt.Sprint("/user/player/", id))
}

func player_del_handler(w http.ResponseWriter, r *http.Request, p []string) {
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
	err = Db.TextPlayer.Get(nil, playerid).Scan(&playertextid, nil, &textid)
	if err == nil {
		_, err = Db.Text.Del(textid)
		if err != nil {
			h.ServerError(w, err)
			return
		}
		_, err = Db.TextPlayer.Del(playertextid)
		if err != nil {
			h.ServerError(w, err)
			return
		}
	} else if err != sql.ErrNoRows {
		h.ServerError(w, err)
		return
	}
	var n int64
	n, err = Db.Player.Del(playerid)
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
