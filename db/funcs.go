package db

import (
	"fmt"
	"strings"
	"database/sql"
	"reflect"
	"errors"
)

func copyMap(src map[string]interface{}) (dest map[string]interface{}) {
	dest = make(map[string]interface{})
	for k, v := range src {
		dest[k] = v
	}
	return
}

func formatMapToInsert(kvs map[string]interface{}) (string, []interface{}) {
	keys := make([]string, 0)
	args := make([]string, 0)
	values := make([]interface{}, 0)
	for k, v := range kvs {
		keys = append(keys, k)
		args = append(args, "?")
		values = append(values, v)
	}
	return fmt.Sprintf("(%s) VALUES (%s)", strings.Join(keys, ","), strings.Join(args, ",")), values
}

func formatMapToSet(kvs map[string]interface{}) (string, []interface{}) {
	sqlitems := make([]string, 0)
	values := make([]interface{}, 0)
	for k, v := range kvs {
		sqlitems = append(sqlitems, fmt.Sprint(k, "=?"))
		values = append(values, v)
	}
	return strings.Join(sqlitems, ","), values
}

func formatMapToWhere(kvs map[string]interface{}) (string, []interface{}) {
	sqlitems := make([]string, 0)
	values := make([]interface{}, 0)
	for k, v := range kvs {
		if v == nil {
			sqlitems = append(sqlitems, fmt.Sprint(k, " is null"))
			continue
		}
		sqlitems = append(sqlitems, fmt.Sprint(k, "=?"))
		values = append(values, v)
	}
	return strings.Join(sqlitems, " AND "), values
}



func scanToStruct(row *sql.Row, Scanner interface{}) error {
	vp := reflect.ValueOf(Scanner)
	if vp.Kind() == reflect.Ptr {
		//获取真实数据
		vp = vp.Elem()
	} else {
		return errors.New("db: the object must be a struct point")
	}
	//如果不是结构体，返回一个错误
	if vp.Kind() != reflect.Struct {
		return errors.New("db: the object must be a struct point")
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