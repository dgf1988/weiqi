package weiqi

import (
	"database/sql"
	"fmt"
	"github.com/dgf1988/weiqi/h"
	"net/http"
)

//post list
func postListHandler(w http.ResponseWriter, r *http.Request, p []string) {
	var u *U
	s := getSession(r)
	if s != nil {
		u = s.User
	}

	posts, err := dbListPostByPage(40, 0)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	cutPostTextMany(posts)

	err = postListHtml().Execute(w, postListData(u, posts), defFuncMap)
	if err != nil {
		h.ServerError(w, err)
	}
}

func postListHtml() *Html {
	return defHtmlLayout().Append(
		defHtmlHead(),
		defHtmlHeader(),
		defHtmlFooter(),
		newHtmlContent("postlist"),
	)
}

func postListData(u *U, posts []P) *Data {
	data := defData()
	data.User = u
	data.Head.Title = "文章列表"
	data.Head.Desc = "围棋文章列表"
	data.Head.Keywords = []string{"围棋", "文章", "新闻", "资料"}
	data.Content["Posts"] = posts
	return data
}

//post id
func postIdHandler(w http.ResponseWriter, r *http.Request, args []string) {
	var u *U
	s := getSession(r)
	if s != nil {
		u = s.User
	}

	id := atoi64(args[0])
	if id <= 0 {
		h.NotFound(w, "找不到文章")
		return
	}

	p, err := dbGetPost(id)
	if err == sql.ErrNoRows {
		h.NotFound(w, "找不到文章")
	} else if err != nil {
		h.ServerError(w, err)
	} else {
		err = postIdHtml().Execute(w, postIdData(u, p), defFuncMap)
		if err != nil {
			h.ServerError(w, err)
		}
	}
}

func postIdHtml() *Html {
	return defHtmlLayout().Append(
		defHtmlHead(),
		defHtmlHeader(),
		defHtmlFooter(),
		newHtmlContent("postid"),
	)
}

func postIdData(u *U, post *P) *Data {
	data := defData()
	data.User = u
	data.Head.Title = post.Title
	data.Head.Desc = "围棋文章之" + post.Title
	data.Head.Keywords = []string{"围棋", "文章", "新闻", "资料"}
	post.Text = parseTextToHtml(post.Text)
	data.Content["Post"] = post
	return data
}

//post edit
func userPostEidtHandler(w http.ResponseWriter, r *http.Request, args []string) {

	//登录验证
	s := getSession(r)
	if s == nil {
		h.SeeOther(w, r, "/login")
		return
	}

	r.ParseForm()
	var (
		action    = "/user/post/add"
		msg       = r.FormValue("editormsg")
		post   *P = nil
		err    error
	)

	if len(args) > 0 {
		action = "/user/post/update"
		post, err = dbGetPost(atoi64(args[0]))
		if err == sql.ErrNoRows || err == ErrPrimaryKey {
			h.NotFound(w, "找不到文章")
			return
		}
		if err != nil {
			h.ServerError(w, err)
			return
		}
	} else {
		post = new(P)
	}

	posts, err := dbListPostByPage(40, 0)
	if err == sql.ErrNoRows {
		posts = make([]P, 0)
	} else if err != nil {
		h.ServerError(w, err)
		return
	}
	err = userPostEditHtml().Execute(w, userPostEditData(s.User, action, msg, post, posts), defFuncMap)
	if err != nil {
		h.ServerError(w, err)
		return
	}
}

func userPostEditHtml() *Html {
	return defHtmlLayout().Append(
		defHtmlHead(),
		defHtmlHeader(),
		defHtmlFooter(),
		newHtmlContent("userpostedit"),
	)
}

func userPostEditData(u *U, action, msg string, post *P, posts []P) *Data {
	data := defData()
	data.User = u
	data.Header.Navs = userNavItems()
	data.Content["Editor"] = Editor{action, msg}
	data.Content["Post"] = post
	data.Content["Posts"] = posts
	return data
}

func handlerUserPostAdd(w http.ResponseWriter, r *http.Request, args []string) {
	//登录验证

	if getSession(r) == nil {
		h.SeeOther(w, r, "/login")
		return
	}

	var p P
	r.ParseForm()
	p.Title = r.FormValue("title")
	p.Text = r.FormValue("text")

	if len(p.Title) > 0 && len(p.Text) > 0 {
		_, err := dbAddPost(&p)
		if err == nil {
			h.SeeOther(w, r, fmt.Sprint("/user/post/?editormsg=", p.Title, "提交成功"))
		} else {
			h.SeeOther(w, r, "/user/post/?editormsg="+err.Error())
		}
	} else {
		h.SeeOther(w, r, "/user/post/?editormsg=标题或内容为空")
	}
}

func handlerUserPostUpdate(w http.ResponseWriter, r *http.Request, args []string) {

	//登录验证

	if getSession(r) == nil {
		h.SeeOther(w, r, "/login")
		return
	}

	r.ParseForm()
	var p P
	p.Id = atoi64(r.FormValue("id"))
	p.Title = r.FormValue("title")
	p.Text = r.FormValue("text")
	if p.Id <= 0 {
		h.Forbidden(w, "错误的参数")
		return
	}
	if len(p.Title) == 0 || len(p.Text) == 0 {
		h.SeeOther(w, r, fmt.Sprint("/user/post/", p.Id, "?editormsg=标题或内容为空"))
		return
	}
	err := dbUpdatePost(&p)
	if err != nil {
		h.SeeOther(w, r, fmt.Sprint("/user/post/", p.Id, "?editormsg=", err.Error()))
		return
	}
	h.SeeOther(w, r, fmt.Sprint("/user/post/", p.Id, "?editormsg=提交成功"))
}

func handlerUserPostDelete(w http.ResponseWriter, r *http.Request, args []string) {

	//登录验证

	if getSession(r) == nil {
		h.SeeOther(w, r, "/login")
		return
	}

	r.ParseForm()
	var id int64 = atoi64(r.FormValue("id"))
	if id <= 0 {
		h.NotFound(w, "找不找文章")
		return
	}
	err := dbDeletePost(id)
	if err != nil {
		h.Forbidden(w, fmt.Sprint(id, "删除失败", err.Error()))
	} else {
		h.SeeOther(w, r, "/user/post/")
	}
}
