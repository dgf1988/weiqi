package weiqi

import (
	"database/sql"
	"fmt"
	"github.com/dgf1988/weiqi/h"
	"net/http"
)

//post list
func postListHandler(w http.ResponseWriter, r *http.Request, p []string) {
	var u *User
	s := getSession(r)
	if s != nil {
		u = s.User
	}
	var posts = [40]Post{}
	n, err := Posts.ListBy(&posts, 0)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	cutPostTextMany(posts[:n])

	err = postListHtml().Execute(w, postListData(u, posts[:n]), defFuncMap)
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

func postListData(u *User, posts []Post) *Data {
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
	var u *User
	s := getSession(r)
	if s != nil {
		u = s.User
	}

	id := atoi64(args[0])
	if id <= 0 {
		h.NotFound(w, "找不到文章")
		return
	}
	var p = new(Post)
	err := Posts.GetBy(id, p)
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

func postIdData(u *User, post *Post) *Data {
	data := defData()
	data.User = u
	data.Head.Title = post.Title
	data.Head.Desc = "围棋文章"
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
		action = "/user/post/add"
		msg    = r.FormValue("editormsg")
		post   = new(Post)
		err    error
	)

	if len(args) > 0 {
		action = "/user/post/update"
		err := Posts.GetBy(args[0], post)
		if err == sql.ErrNoRows || err == ErrPrimaryKey {
			h.NotFound(w, "找不到文章")
			return
		}
		if err != nil {
			h.ServerError(w, err)
			return
		}
	}

	var posts = [40]Post{}
	n, err := Posts.ListBy(&posts, 0)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	err = userPostEditHtml().Execute(w, userPostEditData(s.User, action, msg, post, posts[:n]), defFuncMap)
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

func userPostEditData(u *User, action, msg string, post *Post, posts []Post) *Data {
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

	var p Post
	r.ParseForm()
	p.Title = r.FormValue("title")
	p.Text = r.FormValue("text")

	if len(p.Title) > 0 && len(p.Text) > 0 {
		_, err := Posts.Add(nil, p.Title, p.Text)
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
	var p Post
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
	_, err := Posts.Set(p.Id, nil, p.Title, p.Text, p.Pstatus)
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
	_, err := Posts.Del(id)
	if err != nil {
		h.Forbidden(w, fmt.Sprint(id, "删除失败", err.Error()))
	} else {
		h.SeeOther(w, r, "/user/post/")
	}
}
