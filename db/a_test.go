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
	table, err := GetTable("hoetom", "player")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(table.ToSql())
	t.Log(table.Get(2))
	list, err := table.List(50, 0)
	if err != nil {
		t.Error(err.Error())
	}
	for i := range list {
		t.Log(list[i])
	}
}