package weiqi

import (
	"errors"
	"github.com/dgf1988/weiqi/db"
	"log"
)

var (
	ErrPrimaryKey = errors.New("primary key error")
	Players       db.ITable
	Users         db.ITable
	Posts         db.ITable
	Sgfs          db.ITable

	Texts		  db.ITable
	PlayerText	  db.ITable
)

func init() {
	var err error
	db.Config(config.DbDriver, config.DbUsername, config.DbPassword, config.DbName, config.DbHost, config.DbPort, config.DbCharset)
	err = db.Connect()
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

	PlayerText, err = db.GetTable(config.DbName, "player_text")
	if err != nil {
		log.Fatal(err.Error())
	}
}
