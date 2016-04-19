package weiqi

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"reflect"
	"strings"
	"github.com/dgf1988/weiqi/db"
)

var (
	databases            *sql.DB
	ErrPrimaryKey = errors.New("primary key error")
	Players *db.Table
	User	*db.Table
	Posts	*db.Table
	Sgfs	*db.Table
)

func init() {
	conn, err := sql.Open(config.DbDriver, config.DbConnectString())
	if err != nil {
		log.Fatal(err.Error())
	}
	err = conn.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
	databases = conn

	db.Config(config.DbDriver, config.DbUsername, config.DbPassword, config.DbName, config.DbHost, config.DbPost, config.DbCharset)
	err = db.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}
	Players, err = db.GetTable(config.DbName, "player")
	if err != nil {
		log.Fatal(err.Error())
	}
	User, err = db.GetTable(config.DbName, "user")
	if err != nil {
		log.Fatal(err.Error())
	}
	Posts, err = db.GetTable(config.DbName, "post")
	if err != nil {
		log.Fatal(err.Error())
	}
	Sgfs, err = db.GetTable(config.DbName,  "sgf")
	if err != nil {
		log.Fatal(err.Error())
	}
}

func dbCount(tablename string) (int64, error) {
	var num int64
	row := databases.QueryRow("select count(*) as num from " + tablename)
	err := row.Scan(&num)
	if err != nil {
		return -1, err
	}
	return num, nil
}

func dbCountBy(tablename string, where string) (int64, error) {
	row := databases.QueryRow("select count(*) as num from " + tablename + " where " + where)
	var num int64
	err := row.Scan(&num)
	if err != nil {
		return -1, err
	}
	return num, nil
}

func dbDesc(tablename string) ([][6]string, error) {
	rows, err := databases.Query(fmt.Sprint("desc ", tablename))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	descs := make([][6]string, 0)
	scans := make([]interface{}, 6)
	for rows.Next() {
		for i := range scans {
			scans[i] = new(sql.NullString)
		}
		err := rows.Scan(scans...)
		if err != nil {
			return nil, err
		}
		var descrow [6]string
		for i := range scans {
			sqlnullvalue := scans[i].(*sql.NullString)
			if sqlnullvalue.Valid {
				descrow[i] = sqlnullvalue.String
			} else {
				descrow[i] = ""
			}
		}
		descs = append(descs, descrow)
	}
	return descs, rows.Err()
}

//保存数据
func dbUpdate(tablename string, id int64, datas map[string]interface{}) (int64, error) {
	sqlitems := make([]string, 0)
	sqlitems = append(sqlitems, "update", tablename, "set")
	keys := make([]string, 0)
	args := make([]interface{}, 0)
	for k, v := range datas {
		keys = append(keys, k+" = ?")
		args = append(args, v)
	}
	sqlitems = append(sqlitems, strings.Join(keys, ","))
	sqlitems = append(sqlitems, "where id = ? limit 1")
	args = append(args, id)
	updatesql := strings.Join(sqlitems, " ")
	res, err := databases.Exec(updatesql, args...)
	if err != nil {
		return 0, err
	} else {
		return res.RowsAffected()
	}
}

//
func dbClear(tablename string) (int64, error) {
	res, err := databases.Exec(fmt.Sprint("TRUNCATE table ", tablename))
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func scanRowToStruct(row sql.Row, object interface{}) error {
	vp := reflect.ValueOf(object)
	if vp.Kind() == reflect.Ptr {
		//获取真实数据
		vp = vp.Elem()
	}
	//如果不是结构体，返回一个错误
	if vp.Kind() != reflect.Struct {
		return errors.New("orm: the object must be a struct point")
	}
	scans := make([]interface{}, vp.NumField())
	//datas := make([][]byte, vp.NumField())
	for i := 0; i < vp.NumField(); i++ {
		scans[i] = vp.Field(i).Addr().Interface()
	}
	err := row.Scan(scans...)
	if err != nil {
		return err
	}
	return nil
}
