package db

import (
	"log"
	"testing"
	"time"
)

//Config("mysql", "weiqi", "tKWywchAVKxjLb4F", "www.weiqi163.com", 3306, "weiqi_2", "utf8")

func TestCount(t *testing.T) {
	log.SetPrefix("[Debug: db]")
	log.SetFlags(log.Ltime)
	Config("mysql", "root", "guofeng001", "weiqi2", "localhost", 3306, "utf8")
	//Config("mysql", "weiqi", "tKWywchAVKxjLb4F", "weiqi_163", "www.weiqi163.com", 3306, "utf8")
	err := Connect()
	if err != nil {
		t.Error(err.Error())
	}

	type Post struct {
		Id     int64
		Title  string
		Text   string
		Status int64
		Posted time.Time
		Update time.Time
	}

	type Player struct {
		Id      int64
		Name    string
		Sex     int64
		Country string
		Rank    string
		Birth   time.Time
	}

	//Posts, err := GetTable("weiqi2", "post")
	Players, err := GetTable("weiqi2", "player")

	rows, err := Players.Find(nil, nil, nil, "什么国")
	if err != nil {
		t.Error(err.Error())
	} else {
		defer rows.Close()
		for rows.Next() {
			player, err := rows.Slice()
			if err != nil {
				t.Error(err.Error())
			} else {
				t.Log(player)
				_, err = Players.Update(player[0]).Values(nil, nil, nil, "我国", "九段")
				t.Log(err)
			}
		}
	}
}
