package db

import (
	"database/sql"
	"fmt"
	"strings"
	"reflect"
	"errors"
	"time"
)

// Table 保存表信息
type Table struct {
	//数据库名
	DatabaseName  string
	//表名
	Name          string
	//字段结构信息
	Columns       []Column
	//字段数量
	ColumnNumbers int
	//主键
	Primarykey    string
	//唯一键
	UniqueIndex   []string

	fullName 	  string
	// 预备Sql执行语句
	sqlSelect     string
	sqlInsert     string
	sqlDelete	  string
	sqlUpdate 	  string

}

// ToSql 输出表结构Sql语句。
func (t Table) ToSql() string {
	stritems := make([]string, 0)
	stritems = append(stritems, fmt.Sprintf("CREATE TABLE `%s` (", t.Name))
	colitems := make([]string, 0)
	for i := range t.Columns {
		colitems = append(colitems, "\t"+t.Columns[i].ToSql())
	}
	if t.Primarykey != "" {
		colitems = append(colitems, fmt.Sprintf("\tPRIMARY KEY (`%s`)", t.Primarykey))
	}
	for i := range t.UniqueIndex {
		colitems = append(colitems, fmt.Sprintf("\tUNIQUE KEY `%s_%d` (`%s`)", t.UniqueIndex[i], i, t.UniqueIndex[i]))
	}
	stritems = append(stritems, strings.Join(colitems, ",\n"), ") ENGINE=InnoDB DEFAULT CHARSET=utf8")
	return strings.Join(stritems, "\n")
}

func (t Table) getScanner() []interface{} {
	scans := make([]interface{}, t.ColumnNumbers)
	for i := range t.Columns {
		switch t.Columns[i].Type.Name {
		case "int":
			scans[i] = new(sql.NullInt64)
		case "date", "year", "datetime", "time", "timestamp":
			scans[i] = new(NullTime)
		default:
			scans[i] = new(sql.NullString)
		}
	}
	return scans
}

func (t Table) getScanner2() []interface{} {
	scanner := make([]interface{}, t.ColumnNumbers)
	for i := range t.Columns {
		scanner[i] = reflect.New(t.Columns[i].Type.Type).Interface()
	}
	return scanner
}

func (t Table) getStructScanner(obj interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return nil, errors.New("db: the object must be a struct point")
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return nil, errors.New("db: the object must be a struct point")
	}
	scans := make([]interface{}, t.ColumnNumbers)
	for i := range scans {
		scans[i] = v.Field(i).Addr().Interface()
	}
	return scans, nil
}

func (t Table) setScanner(dest []interface{}, obj interface{}) error {
	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		v := reflect.ValueOf(obj).Elem()
		if v.Kind() == reflect.Struct {
			for i := range dest {
				dest[i] = v.Field(i).Addr().Interface()
			}
			return nil
		}
	}
	return errors.New("db: the object can't point to a struct interface")
}

func (t Table) parseScansArray(scans []interface{}) []interface{} {
	data := make([]interface{}, 0)
	for i := range scans {
		switch scans[i].(type) {
		case *sql.NullString:
			nullstr := scans[i].(*sql.NullString)
			if nullstr.Valid {
				data = append(data, nullstr.String)
			} else {
				data = append(data, nil)
			}
		case *sql.NullInt64:
			nullint := scans[i].(*sql.NullInt64)
			if nullint.Valid {
				data = append(data, nullint.Int64)
			} else {
				data = append(data, nil)
			}
		case *NullTime:
			nulltime := scans[i].(*NullTime)
			if nulltime.Valid {
				data = append(data, nulltime.Time)
			} else {
				data = append(data, nil)
			}
		case *string:
			data = append(data, *scans[i].(*string))
		case *int:
			data = append(data, *scans[i].(*int))
		case *time.Time:
			data = append(data, *scans[i].(*time.Time))
		default:
			data = append(data, nil)

		}
	}
	return data
}

func (t Table) parseScansMap(scans []interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	for i := range t.Columns {
		switch scans[i].(type) {
		case *sql.NullString:
			nullstr := scans[i].(*sql.NullString)
			if nullstr.Valid {
				data[t.Columns[i].Name] = nullstr.String
			} else {
				data[t.Columns[i].Name] = nil
			}
		case *sql.NullInt64:
			nullint := scans[i].(*sql.NullInt64)
			if nullint.Valid {
				data[t.Columns[i].Name] = nullint.Int64
			} else {
				data[t.Columns[i].Name] = nil
			}
		case *NullTime:
			nulltime := scans[i].(*NullTime)
			if nulltime.Valid {
				data[t.Columns[i].Name] = nulltime.Time
			} else {
				data[t.Columns[i].Name] = nil
			}
		}
	}
	return data
}

func (t Table) Get(obj interface{}, key interface{}) error {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return errors.New("db: the object must be a pointer which point to a struct")
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("db: the object must be a pointer which point to a struct")
	}
	dest := make([]interface{}, t.ColumnNumbers)
	for i := range dest {
		dest[i] = v.Field(i).Addr().Interface()
	}
	row := queryRow(fmt.Sprintf("select * from %s.%s where %s = ? limit 1", t.DatabaseName, t.Name, t.Primarykey), key)
	err := row.Scan(dest...)
	if err != nil {
		return err
	}
	return nil
}

// GetArray 按主键取数据
func (t Table) GetArray(key interface{}) ([]interface{}, error) {
	dest := t.getScanner2()
	row := queryRow(fmt.Sprintf("select * from %s.%s where %s = ? limit 1", t.DatabaseName, t.Name, t.Primarykey), key)
	err := row.Scan(dest...)
	if err != nil {
		return nil, err
	}
	return t.parseScansArray(dest), nil
}

// GetMap	按主键取数据，输出字典
func (t Table) GetMap(key interface{}) (map[string]interface{}, error) {
	row := queryRow(fmt.Sprintf("select * from %s.%s where %s = ? limit 1", t.DatabaseName, t.Name, t.Primarykey), key)
	dest := t.getScanner()
	err := row.Scan(dest...)
	if err != nil {
		return nil, err
	}
	return t.parseScansMap(dest), nil
}

func (t Table) Find(obj interface{}, where string, args ...interface{}) error {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return errors.New("db: the object must be a pointer which point to a struct")
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("db: the object must be a pointer which point to a struct")
	}
	dest := make([]interface{}, t.ColumnNumbers)
	for i := range dest {
		dest[i] = v.Field(i).Addr().Interface()
	}
	row := queryRow(fmt.Sprintf("select * from %s.%s where %s limit 1", t.DatabaseName, t.Name, where), args...)
	err := row.Scan(dest...)
	if err != nil {
		return err
	}
	return nil
}

// FindArray 查找数据
func (t Table) FindArray(where string, args ...interface{}) ([]interface{}, error) {
	row := queryRow(fmt.Sprintf("select * from %s.%s where %s limit 1", t.DatabaseName, t.Name, where), args...)
	scans := t.getScanner()
	err := row.Scan(scans...)
	if err != nil {
		return nil, err
	}
	return t.parseScansArray(scans), nil
}

// FindMap	查找数据，输出字典
func (t Table) FindMap(where string, args ...interface{}) (map[string]interface{}, error) {
	row := queryRow(fmt.Sprintf("select * from %s.%s where %s limit 1", t.DatabaseName, t.Name, where), args...)
	scans := t.getScanner()
	err := row.Scan(scans...)
	if err != nil {
		return nil, err
	}
	return t.parseScansMap(scans), nil
}

func (t Table) List(objs interface{}, skip int) error {
	vs := reflect.ValueOf(objs)
	if vs.Kind() == reflect.Ptr{
		vs = vs.Elem()
		if vs.Kind() != reflect.Array {
			return errors.New("db: the object must be a pointer which point to array of struct")
		}
	} else if vs.Kind() != reflect.Slice {
		return errors.New("db: the object must be a pointer which point to array of struct")
	}
	rows, err := query(fmt.Sprintf("select * from %s.%s order by %s desc limit ?,?", t.DatabaseName, t.Name, t.Primarykey), skip, vs.Len())
	if err != nil {
		return err
	}
	defer rows.Close()
	dest := make([]interface{}, t.ColumnNumbers)
	for i := 0; i < vs.Len(); i ++ {
		v := vs.Index(i)
		if v.Kind() == reflect.Struct {
			v = v.Addr().Elem()
		} else if v.Kind() == reflect.Ptr {
			v = v.Elem()
			if v.Kind() != reflect.Struct {
				return errors.New("db: the object must be a pointer which point to array of struct")
			}
		} else {
			return errors.New("db: the object must be a pointer which point to array of struct")
		}
		for i := 0; i < v.NumField(); i ++ {
			dest[i] = v.Field(i).Addr().Interface()
		}
		if !rows.Next() {
			break
		}
		err = rows.Scan(dest...)
		if err != nil {
			return err
		}
	}
	return rows.Err()
}

// ListArray 列出数据
func (t Table) ListArray(take, skip int) ([][]interface{}, error) {
	rows, err := query(fmt.Sprintf("select * from %s.%s order by %s desc limit ?,?", t.DatabaseName, t.Name, t.Primarykey), skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dest := t.getScanner()
	results := make([][]interface{}, 0)
	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			return nil, err
		}
		results = append(results, t.parseScansArray(dest))
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// List 列出数据，输出字典
func (t Table) ListMap(take, skip int) ([]map[string]interface{}, error) {
	rows, err := query(fmt.Sprintf("select * from %s.%s order by %s desc limit ?,?", t.DatabaseName, t.Name, t.Primarykey), skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dest := t.getScanner()
	results := make([]map[string]interface{}, 0)
	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			return nil, err
		}
		results = append(results, t.parseScansMap(dest))
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// Query 查询数据
func (t Table) QueryArray(take, skip int, where string, args ...interface{}) ([][]interface{}, error) {
	rows, err := query(fmt.Sprintf("select * from %s.%s where %s order by %s desc limit %d,%d", t.DatabaseName, t.Name, where, t.Primarykey, skip, take), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dest := t.getScanner()
	results := make([][]interface{}, 0)
	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			return nil, err
		}
		results = append(results, t.parseScansArray(dest))
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// QueryMap 查询数据，输出字典
func (t Table) QueryMap(take, skip int, where string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := query(fmt.Sprintf("select * from %s.%s where %s order by %s desc limit %d,%d", t.DatabaseName, t.Name, where, t.Primarykey, skip, take), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dest := t.getScanner()
	results := make([]map[string]interface{}, 0)
	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			return nil, err
		}
		results = append(results, t.parseScansMap(dest))
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// Set 按主键修改数据
func (t Table) Set(key interface{}, kvs map[string]interface{}) (int64, error) {
	setsql, values := formatMapToSet(kvs)
	res, err := exec(fmt.Sprintf("update %s.%s set %s where %s=? limit 1", t.DatabaseName, t.Name, setsql, t.Primarykey), append(values, key)...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

// SetMany 批量修改数据
func (t Table) SetMany(kvs map[string]interface{}, where string, args ...interface{}) (int64, error) {
	setsql, values := formatMapToSet(kvs)
	for i := range args {
		values = append(values, args[i])
	}
	res, err := exec(fmt.Sprintf("update %s.%s set %s where %s", t.DatabaseName, t.Name, setsql, where), values...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

// Add 添加数据
func (t Table) Add(kvs map[string]interface{}) (int64, error) {
	addsql, values := formatMapToInsert(kvs)
	res, err := exec(fmt.Sprintf("insert into %s.%s %s", t.DatabaseName, t.Name, addsql), values...)
	if err != nil {
		return -1, err
	}
	return res.LastInsertId()
}

// Del 按主键删除数据
func (t Table) Del(key interface{}) (int64, error) {
	res, err := exec(fmt.Sprintf("delete from %s.%s where %s=? limit 1", t.DatabaseName, t.Name, t.Primarykey), key)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

// Save 保存数据
func (t Table) Save(kvs map[string]interface{}) (int64, error) {
	key, ok := kvs[t.Primarykey]
	if ok {
		sets := copyMap(kvs)
		delete(sets, t.Primarykey)
		return t.Set(key, sets)
	} else {
		return t.Add(kvs)
	}
}

// Remove 移除数据
func (t Table) Remove(kvs map[string]interface{}) (int64, error) {
	wheresql, values := formatMapToWhere(kvs)
	res, err := exec(fmt.Sprintf("delete from %s.%s where %s limit 1", t.DatabaseName, t.Name, wheresql), values...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

// Count 统计
func (t Table) Count() int64 {
	return count(t.fullName)
}

// Count 条件统计
func (t Table) CountBy(where string, args ...interface{}) int64 {
	return countBy(t.fullName, where, args...)
}

func newTable() *Table {
	return &Table{
		Columns: make([]Column, 0), UniqueIndex: make([]string, 0),
	}
}

func GetTable(databasename, tablename string) (*Table, error) {
	cols, err := GetColumns(databasename, tablename)
	if err != nil {
		return nil, err
	}
	var t = newTable()
	//设置表名
	t.DatabaseName = databasename
	t.Name = tablename
	t.fullName = fmt.Sprint(t.DatabaseName, ".", t.Name)

	//保存字段
	t.Columns = cols
	t.ColumnNumbers = len(cols)


	//遍历字段，读取其它信息。
	itemsSelect := make([]string, 0)
	itemsInsertKey := make([]string, 0)
	itemsInsertArgs := make([]string, 0)
	for _, col := range cols {
		itemsSelect = append(itemsSelect, col.FullName)
		itemsInsertArgs = append(itemsInsertArgs, "?")
		itemsInsertKey = append(itemsInsertKey, col.FullName)
		if col.Key == "PRI" {
			t.Primarykey = col.Name
		} else if col.Key == "UNI" {
			t.UniqueIndex = append(t.UniqueIndex, col.Name)
		}
	}

	//保存预备SQL语句。
	t.sqlSelect = fmt.Sprintf("SELECT %s FROM %s ", strings.Join(itemsSelect, ", "), t.fullName)
	t.sqlInsert = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", t.fullName, strings.Join(itemsInsertKey, ", "), strings.Join(itemsInsertArgs, ", "))
	t.sqlDelete = fmt.Sprintf("DELETE FROM %s WHERE %s = ?", t.fullName, t.Primarykey)
	t.sqlUpdate = fmt.Sprintf("UPDATE %s ", t.fullName)

	return t, nil
}
