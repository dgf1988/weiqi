package weiqi

import (
	"flag"
	"net/http"
	"github.com/dgf1988/weiqi/h"
	"log"
)

const (
	Time_Format     = "2006年01月02日 15:04"
	Time_Def_Format = "2006-01-02 15:04:05"
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
	m.HandleFunc(sgfIdHandler, "/sgf/*", GET)
	m.HandleFunc(userSgfEditHandler, "/user/sgf/", GET)
	m.HandleFunc(userSgfEditHandler, "/user/sgf/*", GET)
	m.HandleFunc(handlerUserSgfAdd, "/user/sgf/add", POST)
	m.HandleFunc(handlerUserSgfDelete, "/user/sgf/del", POST)
	m.HandleFunc(handlerUserSgfUpdate, "/user/sgf/update", POST)
	/*

		mux.HandleFunc(postListHandler, "/post/", GET)
		mux.HandleFunc(postIdHandler, "/post/*", GET)
		mux.HandleFunc(userPostEidtHandler, "/user/post/", GET)
		mux.HandleFunc(userPostEidtHandler, "/user/post/*", GET)
		mux.HandleFunc(handlerUserPostAdd, "/user/post/add", POST)
		mux.HandleFunc(handlerUserPostUpdate, "/user/post/update", POST)
		mux.HandleFunc(handlerUserPostDelete, "/user/post/del", POST)

		mux.HandleFunc(playerListHandler, "/player/", GET)
		mux.HandleFunc(playerIdHandler, "/player/*", GET)
		mux.HandleFunc(userPlayerEditHandler, "/user/player/", GET)
		mux.HandleFunc(userPlayerEditHandler, "/user/player/*", GET)
		mux.HandleFunc(handlerUserPlayerAdd, "/user/player/add", POST)
		mux.HandleFunc(handlerUserPlayerDel, "/user/player/del", POST)
		mux.HandleFunc(handlerUserPlayerUpdate, "/user/player/update", POST)
	*/

	http.Handle("/", m)
	http.Handle("/static/", http.FileServer(http.Dir(config.BasePath)))
	log.Fatal(http.ListenAndServe(*port, nil))
}
