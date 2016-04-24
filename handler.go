package weiqi

import (
	"github.com/dgf1988/weiqi/h"
	"net/http"
)

//默认处理器。处理首页访问
func handleDefault(w http.ResponseWriter, r *http.Request, args []string) {
	//从会话中获取用户信息，如果没登录，则为nil。
	u := getSessionUser(r)

	err := renderDefault(w, u)
	if err != nil {
		h.ServerError(w, err)
	}
}

func renderDefault(w http.ResponseWriter, u *User) error {
	data := defData()
	data.User = u

	var (
		posts   = make([]Post, 0)
		players = make([]Player, 0)
		sgfs    = make([]Sgf, 0)
	)

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
	}

	if rows, err := Posts.List(40, 0); err != nil {
		return err
	} else {
		defer rows.Close()
		for rows.Next() {
			var post Post
			err = rows.Struct(&post)
			if err != nil {
				return err
			} else {
				posts = append(posts, post)
			}
		}
	}

	if rows, err := Sgfs.List(40, 0); err != nil {
		return err
	} else {
		defer rows.Close()
		for rows.Next() {
			var sgf Sgf
			err = rows.Struct(&sgf)
			if err != nil {
				return err
			} else {
				sgfs = append(sgfs, sgf)
			}
		}
	}

	data.Content["Posts"] = posts
	data.Content["Sgfs"] = sgfs
	data.Content["Players"] = players

	return defHtmlLayout().Append(
		defHtmlHead(),
		defHtmlHeader(),
		defHtmlContent(),
		defHtmlFooter(),
	).Execute(w, data, defFuncMap)
}
