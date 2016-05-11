package mysql

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "fmt"
)

type DB struct {
    *sql.DB
    DBName string
}

func Open(username, password, hostname string, port int, databasename string) (*DB, error) {
    sqldb, err := sql.Open("mysql",
        fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true", username, password, hostname, port, databasename))
    if err != nil {
        return nil, err
    }
    if err = sqldb.Ping(); err != nil {
        return nil, err
    }
    return &DB{DB: sqldb, DBName:databasename}, nil
}

func (db *DB) Use(databasename string) error {
    _, err := db.Exec(fmt.Sprintf("use %s", databasename))
    db.DBName = databasename
    return err
}

func (db *DB) ShowTables() ([]string, error) {
    rows, err := db.Query("show tables")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    tables := make([]string, 0)
    for rows.Next() {
        var tablename string
        err = rows.Scan(&tablename)
        if err != nil {
            return nil, err
        }
        tables = append(tables, tablename)
    }
    if err = rows.Close(); err != nil {
        return nil, err
    }
    return tables, nil
}

func (db *DB) DescTable(tablename string) (*sql.Rows, error) {
    return db.Query(fmt.Sprintf("desc %s", tablename))
}