package db

import "testing"

//Config("mysql", "weiqi", "tKWywchAVKxjLb4F", "www.weiqi163.com", 3306, "weiqi_2", "utf8")

func TestCount(t *testing.T) {
	Config("mysql", "root", "guofeng001", "localhost", "weiqi2", 3306, "utf8")
	err := Connect()
	if err != nil {
		t.Error(err.Error())
	}
	n := count("sgf")
	t.Log("num=", n)
}

func TestDesc(t *testing.T) {
	players, err := GetTable("weiqi2", "player")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(players.Add(map[string]interface{}{"pname": "我是谁。。。。。。。。。。。?"}))
}