package db

import (
	"database/sql"
	"reflect"
	"time"
)

// Table 保存表信息
type typeTable struct {
	//数据库名
	DatabaseName          string
	//表名
	Name                  string
	//字段结构信息
	Columns               []Column
	//字段数量
	ColumnNumbers         int
	//主键
	Primarykey            string
	//唯一键
	UniqueIndex           []string

	Fullname              string
	// 预备Sql执行语句
	sqlInsert             string

	sqlSelect             string

	sqlDelete			  string
	sqlUpdate             string

	sqlSelectCount		  string
	sqlArgMark			  []string
}

func newTable() *typeTable {
	return &typeTable{
		Columns: make([]Column, 0), UniqueIndex: make([]string, 0), sqlArgMark:make([]string, 0),
	}
}

func (t typeTable) makeScans() []interface{} {
	scans := make([]interface{}, t.ColumnNumbers)
	for i := range t.Columns {
		switch t.Columns[i].Type.Value {
		case typeInt, typeBigint:
			scans[i] = new(int64)
		case typeDate, typeDatetime, typeYear, typeTimestamp, typeTime:
			scans[i] = new(time.Time)
		case typeChar, typeVarchar, typeText, typeMediumTtext, typeLongtext:
			scans[i] = new(string)
		case typeFloat, typeDouble, typeDecimal:
			scans[i] = new(float64)

		default:
			scans[i] = new([]byte)
		}
	}
	return scans
}

func (t typeTable) makeNullableScans() []interface{} {
	scans := make([]interface{}, t.ColumnNumbers)
	for i := range t.Columns {
		switch t.Columns[i].Type.Value {
		case typeInt, typeBigint:
			scans[i] = new(sql.NullInt64)
		case typeDate, typeDatetime, typeYear, typeTimestamp, typeTime:
			scans[i] = new(NullTime)
		case typeChar, typeVarchar, typeText, typeMediumTtext, typeLongtext:
			scans[i] = new(sql.NullString)
		case typeFloat, typeDouble, typeDecimal:
			scans[i] = new(sql.NullFloat64)

		default:
			scans[i] = new(NullBytes)
		}
	}
	return scans
}

func (t typeTable) makeStructScans(object interface{}) ([]interface{}, error) {
	scans := make([]interface{}, t.ColumnNumbers)
	rv := reflect.ValueOf(object)
	if rv.Kind() != reflect.Ptr {
		return nil, NewErrorf("db: the object (%s) is not a pointer", rv.Kind())
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return nil, NewErrorf("db: the pointer (%s) can't point to a struct object", rv.Kind())
	}
	if rv.NumField() != t.ColumnNumbers {
		return nil, NewErrorf("db: the object field numbers (%d) not equals table column numbers (%d)", rv.NumField(), t.ColumnNumbers)
	}
	for i := range scans {
		scans[i] = rv.Field(i).Addr().Interface()
	}
	return scans, nil
}

func (t typeTable) parseSlice(scans []interface{}) []interface{} {
	data := make([]interface{}, t.ColumnNumbers)
	for i := range scans {
		data[i] = parseValue(scans[i])
	}
	return data
}

func (t typeTable) parseMap(scans []interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	for i := range t.Columns {
		data[t.Columns[i].Name] = parseValue(scans[i])
	}
	return data
}

func (typeTable) sqlQuery(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(query, args...)
}

func (typeTable) sqlQueryRow(query string, args ...interface{}) *sql.Row {
	return db.QueryRow(query, args...)
}

func (typeTable) sqlExec(query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
}
