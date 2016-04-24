package weiqi

import (
	"github.com/dgf1988/weiqi/db"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type weiqiDb struct {
	Player db.Table
	User db.Table
	Post db.Table
	Sgf db.Table

	Text db.Table
	TextPlayer db.Table
}

var Db weiqiDb = weiqiDb{}

var (
	Players db.Table
	Users   db.Table
	Posts   db.Table
	Sgfs    db.Table

	Texts      db.Table
	TextPlayer db.Table
)

func init() {
	var err error
	err = db.Connect(config.DbDriver, config.DbUsername, config.DbPassword, config.DbHost, config.DbPort, config.DbName)
	if err != nil {
		log.Fatal(err.Error())
	}
	Players, err = db.GetTable(config.DbName, "player")
	if err != nil {
		log.Fatal(err.Error())
	}
	Db.Player = Players

	Users, err = db.GetTable(config.DbName, "user")
	if err != nil {
		log.Fatal(err.Error())
	}
	Db.User = Users

	Posts, err = db.GetTable(config.DbName, "post")
	if err != nil {
		log.Fatal(err.Error())
	}
	Db.Post = Posts

	Sgfs, err = db.GetTable(config.DbName, "sgf")
	if err != nil {
		log.Fatal(err.Error())
	}
	Db.Sgf = Sgfs

	Texts, err = db.GetTable(config.DbName, "text")
	if err != nil {
		log.Fatal(err.Error())
	}
	Db.Text = Texts

	TextPlayer, err = db.GetTable(config.DbName, "text_player")
	if err != nil {
		log.Fatal(err.Error())
	}
	Db.TextPlayer = TextPlayer
}
