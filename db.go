package weiqi

import (
	"errors"
	"github.com/dgf1988/weiqi/db"
	"log"
)

var (
	ErrPrimaryKey = errors.New("primary key error")
	Players       *db.Table
	Users         *db.Table
	Posts         *db.Table
	Sgfs          *db.Table
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
}
