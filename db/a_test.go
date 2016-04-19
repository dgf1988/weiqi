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
	options, err := GetTable("weiqi", "option")
	if err != nil {
		t.Error(err.Error())
	}
	ops, err := options.ListArray(10, 0)
	if err != nil {
		t.Error(err.Error())
	}
	for i := range ops {
		t.Log(i, ops[i])
	}

	type Option struct{
		Id int
		Key string
		Value string
		Status int
	}
	var op Option
	err = options.Get(&op, 80)
	if err != nil {
		t.Error(err.Error())
	} else {
		t.Log(op)
	}

	var arrop = [10]Option{}
	err = options.List(&arrop, 0)
	if err != nil {
		t.Error(err.Error())
	}
	for i := range arrop {
		t.Log(i, arrop[i])
	}
}
