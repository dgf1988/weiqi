package weiqi

import (
	"flag"
	"net/http"
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
	mux := NewMux()
	mux.HandleFunc(defaultHandler, "/", GET)

	mux.HandleFunc(loginHandler, "/login", GET, POST)
	mux.HandleFunc(handlerLogout, "/logout", GET)
	mux.HandleFunc(handlerRegister, "/register", POST)

	mux.HandleFunc(postListHandler, "/post/", GET)
	mux.HandleFunc(postIdHandler, "/post/*", GET)
	mux.HandleFunc(userPostEidtHandler, "/user/post/", GET)
	mux.HandleFunc(userPostEidtHandler, "/user/post/*", GET)
	mux.HandleFunc(handlerUserPostAdd, "/user/post/add", POST)
	mux.HandleFunc(handlerUserPostUpdate, "/user/post/update", POST)
	mux.HandleFunc(handlerUserPostDelete, "/user/post/del", POST)

	mux.HandleFunc(sgfListHandler, "/sgf/", GET)
	mux.HandleFunc(sgfIdHandler, "/sgf/*", GET)
	mux.HandleFunc(userSgfEditHandler, "/user/sgf/", GET)
	mux.HandleFunc(userSgfEditHandler, "/user/sgf/*", GET)
	mux.HandleFunc(handlerUserSgfAdd, "/user/sgf/add", POST)
	mux.HandleFunc(handlerUserSgfDelete, "/user/sgf/del", POST)
	mux.HandleFunc(handlerUserSgfUpdate, "/user/sgf/update", POST)

	mux.HandleFunc(playerListHandler, "/player/", GET)
	mux.HandleFunc(playerIdHandler, "/player/*", GET)
	mux.HandleFunc(userPlayerEditHandler, "/user/player/", GET)
	mux.HandleFunc(userPlayerEditHandler, "/user/player/*", GET)
	mux.HandleFunc(handlerUserPlayerAdd, "/user/player/add", POST)
	mux.HandleFunc(handlerUserPlayerDel, "/user/player/del", POST)
	mux.HandleFunc(handlerUserPlayerUpdate, "/user/player/update", POST)

	mux.HandleFunc(userHandler, "/user", GET)

	http.Handle("/", mux)

	http.Handle("/static/", http.FileServer(http.Dir(config.BasePath)))

	http.ListenAndServe(*port, nil)
}
