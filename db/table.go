package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// Table 保存表信息
type Table struct {
	//数据库名
	DatabaseName             string
	//表名
	Name                     string
	//字段结构信息
	Columns                  []Column
	//字段数量
	ColumnNumbers            int
	//主键
	Primarykey               string
	//唯一键
	UniqueIndex              []string

	fullName                 string
	// 预备Sql执行语句
	sqlInsert                string
	sqlDeleteByPrimarykey    string

	sqlSelect                string
	sqlSelectByPrimarykey    string

	sqlUpdate                string

	sqlOrderByPrimarykey     string
}

func newTable() *Table {
	return &Table{
		Columns: make([]Column, 0), UniqueIndex: make([]string, 0),
	}
}

func (t Table) makeScans() []interface{} {
	return makeScans(t.ColumnNumbers)
}

func (t Table) makeNullableScans() []interface{} {
	scans := makeScans(t.ColumnNumbers)
	for i := range t.Columns {
		switch t.Columns[i].Type.Value {
		case typeInt:
			scans[i] = new(sql.NullInt64)
		case typeDate, typeDatetime, typeYear, typeTimestamp, typeTime:
			scans[i] = new(NullTime)
		case typeChar, typeVarchar, typeText, typeMediumTtext, typeLongtext:
			scans[i] = new(sql.NullString)
		case typeFloat:
			scans[i] = new(sql.NullBool)
		default:
			scans[i] = new(NullBytes)
		}
	}
	return scans
}

func (t Table) makeStructScans(object interface{}) ([]interface{}, error) {
	scans := t.makeScans()
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

func (t Table) parseSlice(scans []interface{}) ([]interface{}, error) {
	var err error
	data := make([]interface{}, t.ColumnNumbers)
	for i := range scans {
		data[i], err = parseValue(scans[i])
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (t Table) parseMap(scans []interface{}) (map[string]interface{}, error) {
	var err error
	data := make(map[string]interface{})
	for i := range t.Columns {
		data[t.Columns[i].Name], err = parseValue(scans[i])
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (t Table) add(args ...interface{}) (int64, error) {
	cols := make([]string, 0)
	marks := make([]string, 0)
	values := make([]interface{}, 0)
	for i := range args {
		if args[i] == nil {
			continue
		}
		cols = append(cols, t.Columns[i].FullName)
		marks = append(marks, "?")
		values = append(values, args[i])
	}
	res, err := dbExec(fmt.Sprintf("%s (%s) VALUES (%s)", t.sqlInsert, strings.Join(cols, ", "), strings.Join(marks, ", ")), values...)
	if err != nil {
		return -1, err
	}
	return res.LastInsertId()
}

func (t Table) del(key interface{}) (int64, error) {
	res, err := dbExec(t.sqlDeleteByPrimarykey, key)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (t Table) set(key interface{}, args ...interface{}) (int64, error) {
	items := make([]string, 0)
	values := make([]interface{}, 0)
	for i := range args {
		if args[i] == nil {
			continue
		}
		items = append(items, t.Columns[i].FullName + "=?")
		values = append(values, args[i])
	}
	res, err := dbExec(fmt.Sprintf("%s SET %s WHERE %s = ? limit 1", t.sqlUpdate, strings.Join(items, ", "), t.Primarykey), append(values, key)...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (t Table) get(key interface{}, scans ...interface{}) error {
	return dbQueryRow(t.sqlSelectByPrimarykey, key).Scan(scans...)
}

func (t Table) find(where string, args ...interface{}) *sql.Row {
	return dbQueryRow(fmt.Sprintf("%s WHERE %s limit 1", t.sqlSelect, where), args...)
}

func (t Table) query(query string, args ...interface{}) (*sql.Rows, error) {
	return dbQuery(fmt.Sprintf("%s %s", t.sqlSelect, query), args...)
}

func (t Table) update(key interface{}, datas map[string]interface{}) (int64, error) {
	sqlSets, valus := formatMapToSet(datas)
	res, err := dbExec(fmt.Sprintf("%s SET %s WHERE %s = ? limit 1", t.sqlUpdate, sqlSets, t.Primarykey), append(valus, key)...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (t Table) udpateMany(datas map[string]interface{}, query string, args ...interface{}) (int64, error) {
	sqlSets, values := formatMapToSet(datas)
	res, err := dbExec(fmt.Sprintf("%s SET %s %s", t.sqlUpdate, sqlSets, query), append(values, args...)...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (Table) sqlQuery(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(query, args...)
}

func (Table) sqlQueryRow(query string, args ...interface{}) *sql.Row {
	return db.QueryRow(query, args...)
}

func (Table) sqlExec(query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
}