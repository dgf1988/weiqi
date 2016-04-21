package db

import (
	"testing"
	"time"
)

//Config("mysql", "weiqi", "tKWywchAVKxjLb4F", "www.weiqi163.com", 3306, "weiqi_2", "utf8")

func TestCount(t *testing.T) {
	//Config("mysql", "root", "guofeng001", "weiqi2", "localhost", 3306, "utf8")
	Config("mysql", "weiqi", "tKWywchAVKxjLb4F", "weiqi_163", "www.weiqi163.com", 3306, "utf8")
	err := Connect()
	if err != nil {
		t.Error(err.Error())
	}

	type Player struct {
		Id      int64
		Name    string
		Sex     int64
		Country string
		Rank    string
		Birth   time.Time
	}

	Players, err := GetTable("weiqi_2", "player")

	for i := range Players.Columns {
		t.Log(Players.Columns[i].FullName)
	}
	if err != nil {
		t.Error(err.Error())
	} else {
		var player = new(Player)
		t.Log(Players.FindBy(player, nil, "古力"))
		t.Log(player)
	}
}
