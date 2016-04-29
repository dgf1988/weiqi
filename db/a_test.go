package db

import (
	"log"
	"testing"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

//Config("mysql", "weiqi", "tKWywchAVKxjLb4F", "www.weiqi163.com", 3306, "weiqi_2", "utf8")

func TestCount(t *testing.T) {
	log.SetPrefix("[Debug: db]")
	log.SetFlags(log.Ltime)
	err := Connect("mysql", "root", "guofeng001", "localhost", 3306, "weiqi_2")
	//err := Connect("mysql", "weiqi", "tKWywchAVKxjLb4F", "www.weiqi163.com", 3306, "weiqi_2")
	if err != nil {
		t.Error(err.Error())
		return
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
	Players, err := GetTable("weiqi_2", "player")
	if err != nil {
		t.Error(err.Error())
		return
	}

	if rows, err := Players.FindAll(nil, nil, nil, "中国"); err != nil {
		t.Error(err.Error())
		return
	} else {
		defer rows.Close()
		players := make([]Player, 0)
		before, _ := time.Parse("2006-01-02", "1990-01-01")
		for rows.Next() {
			var player Player
			players = append(players, player)
			if err := rows.Struct(&player); err != nil {
				t.Error(err.Error())
			} else {
				if player.Birth.After(before) {
					t.Log(player)
				}
			}
		}
	}
}
