package db

import (
	"testing"
	"time"
)

//Config("mysql", "weiqi", "tKWywchAVKxjLb4F", "www.weiqi163.com", 3306, "weiqi_2", "utf8")

func TestCount(t *testing.T) {
	Config("mysql", "root", "guofeng001", "weiqi2", "localhost", 3306, "utf8")
	//Config("mysql", "weiqi", "tKWywchAVKxjLb4F", "weiqi_163", "www.weiqi163.com", 3306, "utf8")
	err := Connect()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestDesc(t *testing.T) {
	Players, err := GetTable("hoetom", "player")
	if err != nil {
		t.Error(err.Error())
	}
	type Player struct {
		ID int64
		PID int64
		Name string
		Sex int64
		Country int64
		Rank int64
		Cate string
		Birth time.Time
	}

	ps := make([]Player, 10)
	err = Players.ListBy(&ps, 0)
	if err != nil {
		t.Error(err.Error())
	} else {
		for i := range ps {
			t.Log(ps[i])
		}
	}
}
