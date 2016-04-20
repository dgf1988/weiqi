package db

import (
	"database/sql"
	"fmt"
	"strings"
	"reflect"
	"errors"

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

func (t Table) makeScans() []interface{} {
	return make([]interface{}, t.ColumnNumbers)
}

func (t Table) makeNullableScans() []interface{} {
	scans := make([]interface{}, t.ColumnNumbers)
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
		return nil, errors.New("db: the object must be a pointer which point to a struct 1")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return nil, errors.New("db: the object must be a pointer which point to a struct 2")
	}
	if rv.NumField() != t.ColumnNumbers {
		return nil, errors.New("db: the object field numbers not eq table column numbers !")
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

func (t Table) GetBy(obj interface{}, key interface{}) error {
	scans, err := t.makeStructScans(obj)
	if err != nil {
		return err
	}
	row := queryRow(fmt.Sprintf("%s where %s = ? limit 1", t.sqlSelect, t.Primarykey), key)
	return row.Scan(scans...)
}

// GetArray 按主键取数据
func (t Table) GetSlice(key interface{}) ([]interface{}, error) {
	scans := t.makeNullableScans()
	row := queryRow(fmt.Sprintf("%s where %s = ? limit 1", t.sqlSelect, t.Primarykey), key)
	err := row.Scan(scans...)
	if err != nil {
		return nil, err
	}
	return t.parseSlice(scans)
}

// GetMap	按主键取数据，输出字典
func (t Table) GetMap(key interface{}) (map[string]interface{}, error) {
	scans := t.makeNullableScans()
	row := queryRow(fmt.Sprintf("%s where %s = ? limit 1", t.sqlSelect, t.Primarykey), key)
	err := row.Scan(scans...)
	if err != nil {
		return nil, err
	}
	return t.parseMap(scans)
}

func (t Table) FindBy(obj interface{}, where string, args ...interface{}) error {
	scans, err := t.makeStructScans(obj)
	if err != nil {
		return err
	}
	row := queryRow(fmt.Sprintf("%s where %s limit 1", t.sqlSelect, where), args...)
	return row.Scan(scans...)
}

// FindArray 查找数据
func (t Table) FindSlice(where string, args ...interface{}) ([]interface{}, error) {
	row := queryRow(fmt.Sprintf("%s where %s limit 1", t.sqlSelect, where), args...)
	scans := t.makeNullableScans()
	err := row.Scan(scans...)
	if err != nil {
		return nil, err
	}
	return t.parseSlice(scans)
}

// FindMap	查找数据，输出字典
func (t Table) FindMap(where string, args ...interface{}) (map[string]interface{}, error) {
	row := queryRow(fmt.Sprintf("%s where %s limit 1", t.sqlSelect, where), args...)
	scans := t.makeNullableScans()
	err := row.Scan(scans...)
	if err != nil {
		return nil, err
	}
	return t.parseMap(scans)
}

func (t Table) ListBy(objs interface{}, skip int) error {
	rv_objs := reflect.Indirect(reflect.ValueOf(objs))
	if rv_objs.Kind() != reflect.Slice && rv_objs.Kind() != reflect.Array {
		return errors.New(fmt.Sprintf("db: the object must be a slice or array %v", rv_objs.Kind()))
	}

	take := rv_objs.Len()
	rows, err := query(fmt.Sprintf("%s order by %s desc limit ?,?", t.sqlSelect, t.Primarykey), skip, take)
	if err != nil {
		return err
	}
	defer rows.Close()

	scans := t.makeScans()
	for i:=0;i<take; i++{
		for j := range scans {
			scans[j] = rv_objs.Index(i).Field(j).Addr().Interface()
		}
		if !rows.Next() {
			break
		}
		err = rows.Scan(scans...)
		if err != nil {
			return err
		}
	}
	return rows.Err()
}

// ListArray 列出数据
func (t Table) ListArray(take, skip int) ([][]interface{}, error) {
	rows, err := query(fmt.Sprintf("%s order by %s desc limit ?,?", t.sqlSelect, t.Primarykey), skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scans := t.makeNullableScans()
	results := make([][]interface{}, 0)
	for rows.Next() {
		err = rows.Scan(scans...)
		if err != nil {
			return nil, err
		}
		a, err := t.parseSlice(scans)
		if err != nil {
			return results, err
		}
		results = append(results, a)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// List 列出数据，输出字典
func (t Table) ListMap(take, skip int) ([]map[string]interface{}, error) {
	rows, err := query(fmt.Sprintf("%s order by %s desc limit ?,?", t.sqlSelect, t.Primarykey), skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dest := t.makeNullableScans()
	results := make([]map[string]interface{}, 0)
	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			return nil, err
		}
		a, err := t.parseMap(dest)
		if err != nil {
			return results, err
		}
		results = append(results, a)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// Query 查询数据
func (t Table) QueryArray(take, skip int, where string, args ...interface{}) ([][]interface{}, error) {
	rows, err := query(fmt.Sprintf("%s WHERE %s ORDER BY %s DESC LIMIT ?,?", t.sqlSelect, where, t.Primarykey), append(args, skip, take)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dest := t.makeNullableScans()
	results := make([][]interface{}, 0)
	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			return nil, err
		}
		a, err := t.parseSlice(dest)
		if err != nil {
			return results, err
		}
		results = append(results, a)
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
	dest := t.makeNullableScans()
	results := make([]map[string]interface{}, 0)
	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			return nil, err
		}
		a, err := t.parseMap(dest)
		if err != nil {
			return results, err
		}
		results = append(results, a)
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
