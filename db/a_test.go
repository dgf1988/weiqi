package db

import (
	"testing"
)

//Config("mysql", "weiqi", "tKWywchAVKxjLb4F", "www.weiqi163.com", 3306, "weiqi_2", "utf8")

func TestCount(t *testing.T) {
	Config("mysql", "root", "guofeng001", "weiqi2", "localhost", 3306, "utf8")
	//Config("mysql", "weiqi", "tKWywchAVKxjLb4F", "weiqi_163", "www.weiqi163.com", 3306, "utf8")
	err := Connect()
	if err != nil {
		t.Error(err.Error())
	}
	n := count("sgf")
	t.Log("num=", n)
}

func TestDesc(t *testing.T) {
	options, err := GetTable("weiqi", "sgf")
	if err != nil {
		t.Error(err.Error())
	}
	for i := range options.Columns {
		t.Log(options.Columns[i])
	}
	t.Log(options.sqlSelect)
	t.Log(options.sqlInsert)
	t.Log(options.sqlUpdate)
	t.Log(options.sqlDelete)

	t.Log(options.GetArray(3))
}
