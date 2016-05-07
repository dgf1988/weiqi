package weiqi

import (
	"flag"
	"github.com/dgf1988/weiqi/mux"
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
	var mymux = mux.New()
	mymux.HandleFuncOld(mux.GET, "/", defaultHandler)

	mymux.HandleFuncOld(mux.GET|mux.POST, "/login", handleLogin)
	mymux.HandleFuncOld(mux.GET, "/logout", handleLogout)
	mymux.HandleFuncOld(mux.POST, "/register", handleRegister)
	mymux.HandleFuncOld(mux.GET, "/user", handleUser)

	mymux.HandleFuncOld(mux.GET, "/sgf/", handleSgfList)
	mymux.HandleFuncOld(mux.GET, "/sgf/+", handleSgfId)
	mymux.HandleFuncOld(mux.GET, "/user/sgf/", handleSgfEdit)
	mymux.HandleFuncOld(mux.GET, "/user/sgf/+", handleSgfEdit)
	mymux.HandleFuncOld(mux.POST, "/user/sgf/remote", sgf_remote_handler)
	mymux.HandleFuncOld(mux.POST, "/user/sgf/add", handleSgfAdd)
	mymux.HandleFuncOld(mux.POST, "/user/sgf/del", handleSgfDel)
	mymux.HandleFuncOld(mux.POST, "/user/sgf/update", handleSgfUpdate)

	mymux.HandleFuncOld(mux.GET, "/post/", handlePostList)
	mymux.HandleFuncOld(mux.GET, "/post/+", handlePostId)
	mymux.HandleFuncOld(mux.GET, "/post/page/+", handlePostListPage)
	mymux.HandleFuncOld(mux.GET, "/user/post/", editPostHandler)
	mymux.HandleFuncOld(mux.GET, "/user/post/+", editPostHandler)
	mymux.HandleFuncOld(mux.POST, "/user/post/add", handlePostAdd)
	mymux.HandleFuncOld(mux.POST, "/user/post/status", handlePostStatus)
	mymux.HandleFuncOld(mux.POST, "/user/post/update", handlePostUpdate)
	mymux.HandleFuncOld(mux.POST, "/user/post/del", handlePostDel)

	mymux.HandleFuncOld(mux.GET, "/player/", player_list_handler)
	mymux.HandleFuncOld(mux.GET, "/player/+", player_info_handler)
	mymux.HandleFuncOld(mux.GET, "/user/player/", player_manage_handler)
	mymux.HandleFuncOld(mux.GET|mux.POST, "/user/player/+", player_editor_handler)
	mymux.HandleFuncOld(mux.POST, "/user/player/add", player_add_handler)
	mymux.HandleFuncOld(mux.POST, "/user/player/del", player_del_handler)

	mymux.HandleFuncOld(mux.GET, "/user/img/", img_list_handler)
	mymux.HandleFuncOld(mux.GET|mux.POST, "/user/img/+", img_editor_handler)
	mymux.HandleFuncOld(mux.POST, "/user/img/upload", img_upload_handler)
	mymux.HandleFuncOld(mux.POST, "/user/img/remove", img_remove_handler)
	mymux.HandleFuncOld(mux.POST, "/user/img/remote", img_remote_handler)

	mymux.HandleStd(mux.GET, "/img/*", http.FileServer(http.Dir(config.UploadPath)))
	mymux.HandleStd(mux.GET, "/static/*", http.FileServer(http.Dir(config.BasePath)))
	mymux.HandleStd(mux.GET, "/+", http.FileServer(http.Dir(config.BasePath+"static/site/")))

	http.Handle("/", mymux)
	log.Fatal(http.ListenAndServe(*port, nil))
}
