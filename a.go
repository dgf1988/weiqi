package weiqi

import (
	"flag"
	"github.com/dgf1988/weiqi/h"
	"log"
	"net/http"
)

const (
	Time_Format     = "2006年01月02日 15:04"
	Time_Def_Format = "2006-01-02 15:04:05"
	GET             = "GET"
	POST            = "POST"
)

//Run ...
//程序的入口
func Run() {
	//
	port := flag.String("p", ":8080", "http server port")
	flag.Parse()

	//
	m := h.NewMux()
	m.HandleFunc(defaultHandler, "/", GET)
	m.HandleFunc(loginHandler, "/login", GET, POST)
	m.HandleFunc(handlerLogout, "/logout", GET)
	m.HandleFunc(handlerRegister, "/register", POST)
	m.HandleFunc(userHandler, "/user", GET)

	m.HandleFunc(sgfListHandler, "/sgf/", GET)
	m.HandleFunc(sgfIdHandler, "/sgf/+", GET)
	m.HandleFunc(userSgfEditHandler, "/user/sgf/", GET)
	m.HandleFunc(userSgfEditHandler, "/user/sgf/+", GET)
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

	m.HandleFunc(playerListHandler, "/player/", GET)
	m.HandleFunc(playerIdHandler, "/player/+", GET)
	m.HandleFunc(userPlayerEditHandler, "/user/player/", GET)
	m.HandleFunc(userPlayerEditHandler, "/user/player/+", GET)
	m.HandleFunc(handlerUserPlayerAdd, "/user/player/add", POST)
	m.HandleFunc(handlerUserPlayerDel, "/user/player/del", POST)
	m.HandleFunc(handlerUserPlayerUpdate, "/user/player/update", POST)

	m.HandleStd(http.FileServer(http.Dir(config.BasePath)), "/static/*", GET)

	http.Handle("/", m)
	log.Fatal(http.ListenAndServe(*port, nil))
}
