package db

import (
	"log"
	"testing"
	_ "github.com/go-sql-driver/mysql"
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

	if rows, err := db.Query("show tables"); err != nil {
		t.Fatal(err.Error())
	} else {
		defer rows.Close()
		for rows.Next() {
			var table string
			if err = rows.Scan(&table); err != nil {
				t.Fatal(err.Error())
			} else {
				t.Log(table)
			}
		}
		if err = rows.Err(); err != nil {
			t.Fatal(err.Error())
		}
	}
}
