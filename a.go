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
	m.HandleFunc(loginHandler, "/login", GET, POST)
	m.HandleFunc(handlerLogout, "/logout", GET)
	m.HandleFunc(handlerRegister, "/register", POST)
	m.HandleFunc(userHandler, "/user", GET)

	m.HandleFunc(handleListSgf, "/sgf/", GET)
	m.HandleFunc(handleSgfId, "/sgf/+", GET)
	m.HandleFunc(handleSgfEdit, "/user/sgf/", GET)
	m.HandleFunc(handleSgfEdit, "/user/sgf/+", GET)
	m.HandleFunc(handlerUserSgfAdd, "/user/sgf/add", POST)
	m.HandleFunc(handlerUserSgfDelete, "/user/sgf/del", POST)
	m.HandleFunc(handlerUserSgfUpdate, "/user/sgf/update", POST)

	m.HandleFunc(postListHandler, "/post/", GET)
	m.HandleFunc(postIdHandler, "/post/+", GET)
	m.HandleFunc(userPostEidtHandler, "/user/post/", GET)
	m.HandleFunc(userPostEidtHandler, "/user/post/+", GET)
	m.HandleFunc(handlerUserPostAdd, "/user/post/add", POST)
	m.HandleFunc(handlerUserPostUpdate, "/user/post/update", POST)
	m.HandleFunc(handlerUserPostDelete, "/user/post/del", POST)

	m.HandleFunc(handleListPlayer, "/player/", GET)
	m.HandleFunc(handlePlayerId, "/player/+", GET)
	m.HandleFunc(userPlayerEditHandler, "/user/player/", GET)
	m.HandleFunc(userPlayerEditHandler, "/user/player/+", GET)
	m.HandleFunc(handlerUserPlayerAdd, "/user/player/add", POST)
	m.HandleFunc(handlerUserPlayerDel, "/user/player/del", POST)
	m.HandleFunc(handlerUserPlayerUpdate, "/user/player/update", POST)

	m.HandleStd(http.FileServer(http.Dir(config.BasePath)), "/static/*", GET)

	http.Handle("/", m)
	log.Fatal(http.ListenAndServe(*port, nil))
}
