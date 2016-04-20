package db

import (
	"fmt"
	"strings"
	"database/sql"
	"reflect"
	"errors"
	"time"
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

func parseValue(src interface{}) (interface{}, error) {
	switch src.(type) {
	case *sql.NullString:
		value := src.(*sql.NullString)
		if value.Valid {
			return value.String, nil
		}
		return nil, nil
	case *sql.NullBool:
		value := src.(*sql.NullBool)
		if value.Valid {
			return value.Bool, nil
		}
		return nil, nil
	case *sql.NullInt64:
		value := src.(*sql.NullInt64)
		if value.Valid {
			return value.Int64, nil
		}
		return nil, nil
	case *sql.NullFloat64:
		value := src.(*sql.NullFloat64)
		if value.Valid {
			return value.Float64, nil
		}
		return nil, nil
	case *NullTime:
		value := src.(*NullTime)
		if value.Valid {
			return value.Time, nil
		}
		return nil, nil
	case *NullBytes:
		value := src.(*NullBytes)
		if value.Valid {
			return value.Bytes, nil
		}
		return nil, nil
	case *string:
		return *src.(*string), nil
	case *int:
		return *src.(*int), nil
	case *int64:
		return *src.(*int64), nil
	case *float64:
		return *src.(*float64), nil
	case *bool:
		return *src.(*bool), nil
	case *time.Time:
		return *src.(*time.Time), nil
	case *[]byte:
		return *src.(*[]byte), nil
	case nil:
		return nil, nil
	}
	return nil, errors.New(fmt.Sprintf("db: unknow src type (%v)", reflect.TypeOf(src)))
}

func parseStruct(scans []interface{}, object interface{}) error {
	rv := reflect.ValueOf(object)
	if rv.Kind() != reflect.Ptr {
		return errors.New("db: the object must be a pointer which point to a struct 1")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return errors.New("db: the object must be a pointer which point to a struct 2")
	}
	for i := range scans {
		v, err := parseValue(scans[i])
		if err != nil {
			return err
		}
		rv.Field(i).Set(reflect.ValueOf(v))
	}
	return nil
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