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
		sqlitems = append(sqlitems, k + "=?")
		values = append(values, v)
	}
	return strings.Join(sqlitems, ", "), values
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

func makeScans(length int) []interface{} {
	return make([]interface{}, length)
}

// 把源数据从指针提取出来，返回给使用的人。
func parseValue(src interface{}) (value interface{},err error) {
	switch src.(type) {
	case *sql.NullString:
		src_value := src.(*sql.NullString)
		if src_value.Valid {
			value = src_value.String
		}
	case *sql.NullBool:
		src_value := src.(*sql.NullBool)
		if src_value.Valid {
			value = src_value.Bool
		}
	case *sql.NullInt64:
		src_value := src.(*sql.NullInt64)
		if src_value.Valid {
			value = src_value.Int64
		}
	case *sql.NullFloat64:
		src_value := src.(*sql.NullFloat64)
		if src_value.Valid {
			value = src_value.Float64
		}
	case *NullTime:
		src_value := src.(*NullTime)
		if src_value.Valid {
			value = src_value.Time
		}
	case *NullBytes:
		src_value := src.(*NullBytes)
		if src_value.Valid {
			value = src_value.Bytes
		}
	case *string:
		value = *src.(*string)
	case *int:
		value = *src.(*int)
	case *int64:
		value = *src.(*int64)
	case *float64:
		value = *src.(*float64)
	case *bool:
		value = *src.(*bool)
	case *time.Time:
		value = *src.(*time.Time)
	case *[]byte:
		value = *src.(*[]byte)
	case nil:
		value = nil
	default:
		err = errors.New(fmt.Sprintf("db: unknow src value type (%v)", reflect.TypeOf(src)))
	}
	return
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