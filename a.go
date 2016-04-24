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
	m.HandleFunc(handleSgfAdd, "/user/sgf/add", POST)
	m.HandleFunc(handleSgfDel, "/user/sgf/del", POST)
	m.HandleFunc(handleSgfUpdate, "/user/sgf/update", POST)

	m.HandleFunc(handlePostList, "/post/", GET)
	m.HandleFunc(handlePostId, "/post/+", GET)
	m.HandleFunc(handlePostEdit, "/user/post/", GET)
	m.HandleFunc(handlePostEdit, "/user/post/+", GET)
	m.HandleFunc(handlePostAdd, "/user/post/add", POST)
	m.HandleFunc(handlePostUpdate, "/user/post/update", POST)
	m.HandleFunc(handlePostDel, "/user/post/del", POST)

	m.HandleFunc(handlePlayerList, "/player/", GET)
	m.HandleFunc(handlePlayerId, "/player/+", GET)
	m.HandleFunc(handlePlayerEdit, "/user/player/", GET)
	m.HandleFunc(handlePlayerEdit, "/user/player/+", GET)
	m.HandleFunc(handlePlayerAdd, "/user/player/add", POST)
	m.HandleFunc(handlePlayerDel, "/user/player/del", POST)
	m.HandleFunc(handlePlayerUpdate, "/user/player/update", POST)

	m.HandleStd(http.FileServer(http.Dir(config.BasePath)), "/static/*", GET)

	http.Handle("/", m)
	log.Fatal(http.ListenAndServe(*port, nil))
}
