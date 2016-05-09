package weiqi

import (
	"github.com/dgf1988/weiqi/db"
	//mysql
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type weiqiDb struct {
	User        db.Table
	Post        db.Table
	Sgf         db.Table
    Img         db.Table

	Text        db.Table
	Item        db.Table

	Player      db.Table
	PlayerText  db.Table

	Project     db.Table
	ProjectItem db.Table

}

//Db 数据库
var Db weiqiDb

func init() {
	var err error
	err = db.Connect(config.DbDriver, config.DbUsername, config.DbPassword, config.DbHost, config.DbPort, config.DbName)
	if err != nil {
		log.Fatal(err.Error())
	}

    Db.User, err = db.GetTable(config.DbName, "user")
    if err != nil {
        log.Fatal(err.Error())
    }

    Db.Post, err = db.GetTable(config.DbName, "post")
    if err != nil {
        log.Fatal(err.Error())
    }

    Db.Sgf, err = db.GetTable(config.DbName, "sgf")
    if err != nil {
        log.Fatal(err.Error())
    }

    Db.Text, err = db.GetTable(config.DbName, "text")
    if err != nil {
        log.Fatal(err.Error())
    }

    Db.Img, err = db.GetTable(config.DbName, "img")
    if err != nil {
        log.Fatal(err.Error())
    }

    Db.Item, err = db.GetTable(config.DbName, "item")
    if err != nil {
        log.Fatal(err.Error())
    }
	Db.Player, err = db.GetTable(config.DbName, "player")
	if err != nil {
		log.Fatal(err.Error())
	}
    Db.PlayerText, err = db.GetTable(config.DbName, "playertext")
    if err != nil {
        log.Fatal(err.Error())
    }

	Db.Project, err = db.GetTable(config.DbName, "project")
	if err != nil {
		log.Fatal(err.Error())
	}

	Db.ProjectItem, err = db.GetTable(config.DbName, "projectitem")
	if err != nil {
		log.Fatal(err.Error())
	}
}
