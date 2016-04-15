package weiqi

import (
	"database/sql"
	"fmt"
)

//post list
func postListHandler(h *Http) {
	var u *U
	s := getSession(h.R)
	if s != nil {
		u = s.User
	}

	posts, err := dbListPostByPage(40, 0)
	if err != nil {
		h.ServerError(err.Error())
		return
	}

	cutPostTextMany(posts)

	err = postListHtml().Execute(h.W, postListData(u, posts), defFuncMap)
	if err != nil {
		h.ServerError(err.Error())
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
func postIdHandler(h *Http) {
	var u *U
	s := getSession(h.R)
	if s != nil {
		u = s.User
	}

	id := atoi64(h.P[0])
	if id <= 0 {
		h.RequestError("page not found").NotFound()
		return
	}

	p, err := dbGetPost(id)
	if err == sql.ErrNoRows {
		h.RequestError("找不到文章").NotFound()
	} else if err != nil {
		h.ServerError(err.Error())
	} else {
		err = postIdHtml().Execute(h.W, postIdData(u, p), defFuncMap)
		if err != nil {
			h.ServerError(err.Error())
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
func userPostEidtHandler(h *Http) {

	//登录验证
	s := getSession(h.R)
	if s == nil {
		h.SeeOther("/login")
		return
	}

	h.R.ParseForm()
	var (
		action    = "/user/post/add"
		msg       = h.R.FormValue("editormsg")
		post   *P = nil
		err    error
	)

	if len(h.P) > 0 {
		action = "/user/post/update"
		post, err = dbGetPost(atoi64(h.P[0]))
		if err == sql.ErrNoRows || err == ErrPrimaryKey {
			h.RequestError("找不到文章").NotFound()
			return
		}
		if err != nil {
			h.ServerError(err.Error())
			return
		}
	} else {
		post = new(P)
	}

	posts, err := dbListPostByPage(40, 0)
	if err == sql.ErrNoRows {
		posts = make([]P, 0)
	} else if err != nil {
		h.ServerError(err.Error())
		return
	}
	err = userPostEditHtml().Execute(h.W, userPostEditData(s.User, action, msg, post, posts), defFuncMap)
	if err != nil {
		h.ServerError(err.Error())
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

func handlerUserPostAdd(h *Http) {

	//登录验证

	if getSession(h.R) == nil {
		h.SeeOther("/login")
		return
	}

	var p P
	h.R.ParseForm()
	p.Title = h.R.FormValue("title")
	p.Text = h.R.FormValue("text")

	if len(p.Title) > 0 && len(p.Text) > 0 {
		_, err := dbAddPost(&p)
		if err == nil {
			h.SeeOther(fmt.Sprint("/user/post/?editormsg=", p.Title, "提交成功"))
		} else {
			h.SeeOther("/user/post/?editormsg=" + err.Error())
		}
	} else {
		h.SeeOther("/user/post/?editormsg=标题或内容为空")
	}
}

func handlerUserPostUpdate(h *Http) {

	//登录验证

	if getSession(h.R) == nil {
		h.SeeOther("/login")
		return
	}

	h.R.ParseForm()
	var p P
	p.Id = atoi64(h.R.FormValue("id"))
	p.Title = h.R.FormValue("title")
	p.Text = h.R.FormValue("text")
	if p.Id <= 0 {
		h.RequestError("提交了错误参数").Forbidden()
		return
	}
	if len(p.Title) == 0 || len(p.Text) == 0 {
		h.SeeOther(fmt.Sprint("/user/post/", p.Id, "?editormsg=标题或内容为空"))
		return
	}
	err := dbUpdatePost(&p)
	if err != nil {
		h.SeeOther(fmt.Sprint("/user/post/", p.Id, "?editormsg=", err.Error()))
		return
	}
	h.SeeOther(fmt.Sprint("/user/post/", p.Id, "?editormsg=提交成功"))
}

func handlerUserPostDelete(h *Http) {

	//登录验证

	if getSession(h.R) == nil {
		h.SeeOther("/login")
		return
	}

	h.R.ParseForm()
	var id int64 = atoi64(h.R.FormValue("id"))
	if id <= 0 {
		h.RequestError("提交了错误参数").Forbidden()
		return
	}
	err := dbDeletePost(id)
	if err != nil {
		h.RequestError(fmt.Sprint(id, "删除失败", err.Error())).Forbidden()
	} else {
		h.SeeOther("/user/post/")
	}
}
