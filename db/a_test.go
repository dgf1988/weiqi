package db

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"testing"
)

//Config("mysql", "weiqi", "tKWywchAVKxjLb4F", "www.weiqi163.com", 3306, "weiqi_2", "utf8")

func TestCount(t *testing.T) {
	log.SetPrefix("[Debug: db]")
	log.SetFlags(log.Ltime)
	err := Connect("mysql", "root", "guofeng001", "localhost", 3306, "weiqi_2")
	//err := Connect("mysql", "weiqi", "tKWywchAVKxjLb4F", "www.weiqi163.com", 3306, "weiqi_2")
	if err != nil {
		t.Fatal(err.Error())
	}

	type Post struct {
		Id     float64
		Title  string
		Text   string
		Status int64
		Posted string
		Update string
	}

	Posts, err := GetTable("weiqi_2", "post")
	if err != nil {
		t.Fatal(err.Error())
	}
	var post = new(Post)
	if err = Posts.Get(9).Struct(post); err != nil {
		t.Fatal(err.Error())
	}
	t.Log(post.Id, post.Title, post.Status, post.Posted, post.Update)
}
