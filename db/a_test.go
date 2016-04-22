package db

import (
	"testing"
	"time"
	"log"
)

//Config("mysql", "weiqi", "tKWywchAVKxjLb4F", "www.weiqi163.com", 3306, "weiqi_2", "utf8")

func TestCount(t *testing.T) {
	log.SetPrefix("[Debug: db]")
	log.SetFlags(log.Ltime)
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

	Posts, err := GetTable("weiqi_2", "post")

	if err != nil {
		t.Error(err.Error())
	} else {
		log.Println(Posts.ToSql())

		rows, err := Posts.Query("where post.title like ?", "%韩国%")
		if err != nil {
			t.Error(err.Error())
		} else {
			defer rows.Close()
			for rows.Next() {
				var id int64
				var title string
				err = rows.Scan(&id, &title)
				if err != nil {
					t.Error(err.Error())
				} else {
					t.Log(id, title)
				}
			}
		}
	}
}
