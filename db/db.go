package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db     *sql.DB
	config map[string]interface{} = make(map[string]interface{})
)

//Config 配置数据库
func Config(driver, user, password, database string, host string, port int, charset string) {
	config["driver"] = driver
	config["user"] = user
	config["password"] = password
	config["host"] = host
	config["port"] = port
	config["database"] = database
	config["charset"] = charset
}

//Connect 连接数据库
func Connect() error {
	var err error
	db, err = dbGetConnect(config["driver"].(string),
		config["user"].(string), config["password"].(string),
		config["host"].(string), config["port"].(int),
		config["database"].(string), config["charset"].(string))
	if err != nil {
		return err
	}
	return nil
}

func dbGetConnect(driver, user, password, host string, port int, database, charset string) (*sql.DB, error) {
	conn, err := sql.Open(driver, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
		user, password, host, port, database, charset))
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
