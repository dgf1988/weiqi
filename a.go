package weiqi

import (
	"flag"
	"github.com/dgf1988/weiqi/h"
	"log"
	"net/http"
)

const (
	GET  = "GET"
	POST = "POST"
)

//Run ...
//程序的入口
func Run() {
	//
	port := flag.String("p", ":8080", "http server port")
	flag.Parse()

	//
	m := h.NewMux()
	m.HandleFunc(defaultHandler, "/", h.GET)

	m.HandleFunc(handleLogin, "/login", h.GET | h.POST)
	m.HandleFunc(handleLogout, "/logout", h.GET)
	m.HandleFunc(handleRegister, "/register", h.POST)
	m.HandleFunc(handleUser, "/user", h.GET)

	m.HandleFunc(handleSgfList, "/sgf/", h.GET)
	m.HandleFunc(handleSgfId, "/sgf/+", h.GET)
	m.HandleFunc(handleSgfEdit, "/user/sgf/", h.GET)
	m.HandleFunc(handleSgfEdit, "/user/sgf/+", h.GET)
	m.HandleFunc(sgf_remote_handler, "/user/sgf/remote", h.POST)
	m.HandleFunc(handleSgfAdd, "/user/sgf/add", h.POST)
	m.HandleFunc(handleSgfDel, "/user/sgf/del", h.POST)
	m.HandleFunc(handleSgfUpdate, "/user/sgf/update", h.POST)

	m.HandleFunc(handlePostList, "/post/", h.GET)
	m.HandleFunc(handlePostId, "/post/+", h.GET)
	m.HandleFunc(handlePostListPage, "/post/page/+", h.GET)
	m.HandleFunc(editPostHandler, "/user/post/", h.GET)
	m.HandleFunc(editPostHandler, "/user/post/+", h.GET)
	m.HandleFunc(handlePostAdd, "/user/post/add", h.POST)
	m.HandleFunc(handlePostStatus, "/user/post/status", h.POST)
	///user/post/status
	m.HandleFunc(handlePostUpdate, "/user/post/update", h.POST)
	m.HandleFunc(handlePostDel, "/user/post/del", h.POST)

	m.HandleFunc(player_list_handler, "/player/", h.GET)
	m.HandleFunc(player_info_handler, "/player/+", h.GET)
	m.HandleFunc(player_manage_handler, "/user/player/", h.GET)
	m.HandleFunc(player_editor_handler, "/user/player/+", h.GET | h.POST )
	m.HandleFunc(player_add_handler, "/user/player/add", h.POST)
	m.HandleFunc(player_del_handler, "/user/player/del", h.POST)

	m.HandleFunc(img_list_handler, "/user/img/", h.GET)
	m.HandleFunc(img_editor_handler, "/user/img/+", h.GET| h.POST)
	m.HandleFunc(img_upload_handler, "/user/img/upload", h.POST)
	m.HandleFunc(img_remove_handler, "/user/img/remove", h.POST)
	m.HandleFunc(img_remote_handler, "/user/img/remote", h.POST)
	m.HandleStd(http.FileServer(http.Dir(config.UploadPath)), "/img/*", h.GET)

	m.HandleStd(http.FileServer(http.Dir(config.BasePath)), "/static/*", h.GET)
	m.HandleStd(http.FileServer(http.Dir(config.BasePath+"static/site/")), "/+", h.GET)

	http.Handle("/", m)
	log.Fatal(http.ListenAndServe(*port, nil))
}