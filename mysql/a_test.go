package mysql

import (
    "testing"
    "database/sql"
)

var db *DB

func init() {
    var err error
    db, err = Open("root", "guofeng001", "localhost", 3306, "weiqi_2")
    if err != nil {
        panic(err.Error())
    }
}

func TestOpen(t *testing.T) {
    var err error
    var id int
    var key string
    var value string
    err = db.QueryRow("select * from item limit 1").Scan(&id, &key, &value)
    if err != nil {
        t.Fatal(err.Error())
    } else {
        t.Log(id, key, value)
    }
}

func TestDB_Use(t *testing.T) {
    err := db.Use("weiqi_2")
    if err != nil {
        t.Fatal(err.Error())
    } else {
        t.Log(db.DBName)
    }
}

func TestDB_ShowTables(t *testing.T) {
    tables, err := db.ShowTables()
    if err != nil {
        t.Fatal(err.Error())
    } else {
        for i := range tables {
            t.Log(tables[i])
        }
    }
}

func TestDB_DescTable(t *testing.T) {
    rows, err := db.DescTable("sgf")
    if err != nil {
        t.Fatal(err.Error())
    }
    defer rows.Close()
    for rows.Next() {
        descrow := make([]interface{}, 6)
        for i := range descrow {
            descrow[i] = &sql.NullString{}
        }
        err = rows.Scan(descrow...)
        if err != nil {
            t.Fatal(err.Error())
        }
        var strs = make([]string, 0)
        for i := range descrow {
            str := descrow[i].(*sql.NullString)
            if str.Valid {
                strs = append(strs, str.String)
            } else {
                strs = append(strs, "NULL")
            }
        }
        t.Log(strs)
    }
    err = rows.Close()
    if err != nil {
        t.Fatal(err.Error())
    }
}
