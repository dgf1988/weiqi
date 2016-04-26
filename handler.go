package weiqi

import (
	"github.com/dgf1988/weiqi/h"
	"net/http"
)

//默认处理器。处理首页访问
func handleDefault(w http.ResponseWriter, r *http.Request, args []string) {
	//从会话中获取用户信息，如果没登录，则为nil。

	err := renderDefault(w, getSessionUser(r))
	if err != nil {
		h.ServerError(w, err)
	}
}

func renderDefault(w http.ResponseWriter, u *User) error {
	data := defData()
	data.User = u

	var (
		posts   = make([]Post, 0)
		players []Player
		sgfs    []Sgf
		err error
	)

	if players, err = listPlayerOrderRankDesc(40, 0); err != nil {
		return err
	}

	if rows, err := Db.Post.List(40, 0); err != nil {
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

	if sgfs, err = listSgfOrderTimeDesc(40, 0); err != nil {
		return err
	}

	data.Content["Posts"] = posts
	data.Content["Sgfs"] = sgfs
	data.Content["Players"] = players

	return defHtmlLayout().Append(
		defHtmlHead(),
		defHtmlHeader(),
		defHtmlContent(),
		defHtmlFooter(),
	).Execute(w, data, nil)
}
