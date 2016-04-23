package db

import (
	"database/sql"
	"fmt"
)

var (
	db *sql.DB
)

func dbGetConnect(driver, user, password, host string, port int, database string) (*sql.DB, error) {
	conn, err := sql.Open(driver,
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true", user, password, host, port, database))
	if err != nil {
		return nil, err
	}
	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func dbQuery(sql string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(sql, args...)
}

func dbQueryRow(sql string, args ...interface{}) *sql.Row {
	return db.QueryRow(sql, args...)
}

func dbExec(sql string, args ...interface{}) (sql.Result, error) {
	return db.Exec(sql, args...)
}
