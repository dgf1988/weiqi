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
	m.HandleFunc(handleDefault, "/", GET)

	m.HandleFunc(handleLogin, "/login", GET, POST)
	m.HandleFunc(handleLogout, "/logout", GET)
	m.HandleFunc(handleRegister, "/register", POST)
	m.HandleFunc(handleUser, "/user", GET)

	m.HandleFunc(handleSgfList, "/sgf/", GET)
	m.HandleFunc(handleSgfId, "/sgf/+", GET)
	m.HandleFunc(handleSgfEdit, "/user/sgf/", GET)
	m.HandleFunc(handleSgfEdit, "/user/sgf/+", GET)
	m.HandleFunc(sgf_remote_handler, "/user/sgf/remote", POST)
	m.HandleFunc(handleSgfAdd, "/user/sgf/add", POST)
	m.HandleFunc(handleSgfDel, "/user/sgf/del", POST)
	m.HandleFunc(handleSgfUpdate, "/user/sgf/update", POST)

	m.HandleFunc(handlePostList, "/post/", GET)
	m.HandleFunc(handlePostId, "/post/+", GET)
	m.HandleFunc(handlePostListPage, "/post/page/+", GET)
	m.HandleFunc(handlePostEdit, "/user/post/", GET)
	m.HandleFunc(handlePostEdit, "/user/post/+", GET)
	m.HandleFunc(handlePostAdd, "/user/post/add", POST)
	m.HandleFunc(handlePostStatus, "/user/post/status", POST)
	///user/post/status
	m.HandleFunc(handlePostUpdate, "/user/post/update", POST)
	m.HandleFunc(handlePostDel, "/user/post/del", POST)

	m.HandleFunc(player_list_handler, "/player/", GET)
	m.HandleFunc(player_info_handler, "/player/+", GET)
	m.HandleFunc(handlePlayerEdit, "/user/player/", GET)
	m.HandleFunc(handlePlayerEdit, "/user/player/+", GET)
	m.HandleFunc(handlePlayerAdd, "/user/player/add", POST)
	m.HandleFunc(handlePlayerDel, "/user/player/del", POST)
	m.HandleFunc(handlePlayerUpdate, "/user/player/update", POST)

	m.HandleFunc(img_list_handler, "/user/img/", GET)
	m.HandleFunc(img_editor_handler, "/user/img/+", GET, POST)
	m.HandleFunc(img_upload_handler, "/user/img/upload", POST)
	m.HandleFunc(img_remove_handler, "/user/img/remove", POST)
	m.HandleFunc(img_remote_handler, "/user/img/remote", POST)
	m.HandleStd(http.FileServer(http.Dir(config.UploadPath)), "/img/*", GET)

	m.HandleStd(http.FileServer(http.Dir(config.BasePath)), "/static/*", GET)
	m.HandleStd(http.FileServer(http.Dir(config.BasePath+"static/site/")), "/+", GET)

	http.Handle("/", m)
	log.Fatal(http.ListenAndServe(*port, nil))
}
