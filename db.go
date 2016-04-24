package weiqi

import (
	"github.com/dgf1988/weiqi/db"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var (
	Players db.ITable
	Users   db.ITable
	Posts   db.ITable
	Sgfs    db.ITable

	Texts      db.ITable
	TextPlayer db.ITable
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
	Users, err = db.GetTable(config.DbName, "user")
	if err != nil {
		log.Fatal(err.Error())
	}
	Posts, err = db.GetTable(config.DbName, "post")
	if err != nil {
		log.Fatal(err.Error())
	}
	Sgfs, err = db.GetTable(config.DbName, "sgf")
	if err != nil {
		log.Fatal(err.Error())
	}

	Texts, err = db.GetTable(config.DbName, "text")
	if err != nil {
		log.Fatal(err.Error())
	}

	TextPlayer, err = db.GetTable(config.DbName, "text_player")
	if err != nil {
		log.Fatal(err.Error())
	}
}
