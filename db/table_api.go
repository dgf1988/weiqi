package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

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
	keys := make([]string, 0)
	for _, col := range cols {
		keys = append(keys, col.FullName)
		if col.Key == "PRI" {
			t.Primarykey = col.Name
		} else if col.Key == "UNI" {
			t.UniqueIndex = append(t.UniqueIndex, col.Name)
		}
	}

	strKeys := strings.Join(keys, ", ")

	//保存预备SQL语句。
	t.sqlInsert = fmt.Sprintf("INSERT INTO %s", t.fullName)
	t.sqlDeleteByPrimarykey = fmt.Sprintf("DELETE FROM %s WHERE %s = ? LIMIT 1", t.fullName, t.Primarykey)

	t.sqlUpdate = fmt.Sprintf("UPDATE %s", t.fullName)

	t.sqlSelect = fmt.Sprintf("SELECT %s FROM %s ", strKeys, t.fullName)
	t.sqlSelectByPrimarykey = fmt.Sprintf("SELECT %s FROM %s WHERE %s = ? limit 1", strKeys, t.fullName, t.Primarykey)

	t.sqlOrderByPrimarykey = fmt.Sprintf("ORDER BY %s", t.Primarykey)
	return t, nil
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

func (t Table) GetBy(key interface{}, obj interface{}) error {
	scans, err := t.makeStructScans(obj)
	if err != nil {
		return err
	}
	return t.get(key, scans...)
}

// GetArray 按主键取数据
func (t Table) GetSlice(key interface{}) ([]interface{}, error) {
	scans := t.makeNullableScans()

	err := t.get(key, scans...)
	if err != nil {
		return nil, err
	}

	return t.parseSlice(scans)
}

// GetMap	按主键取数据，输出字典
func (t Table) GetMap(key interface{}) (map[string]interface{}, error) {
	scans := t.makeNullableScans()

	err := t.get(key, scans...)
	if err != nil {
		return nil, err
	}

	return t.parseMap(scans)
}

func (t Table) FindBy(obj interface{}, args ...interface{}) error {
	scans, err := t.makeStructScans(obj)
	if err != nil {
		return err
	}
	return t.find(args...).Scan(scans...)
}

// FindArray 查找数据
func (t Table) FindSlice(args ...interface{}) ([]interface{}, error) {
	scans := t.makeNullableScans()

	err := t.find(args...).Scan(scans...)
	if err != nil {
		return nil, err
	}
	return t.parseSlice(scans)
}

// FindMap	查找数据，输出字典
func (t Table) FindMap(args ...interface{}) (map[string]interface{}, error) {
	scans := t.makeNullableScans()

	err := t.find(args...).Scan(scans...)
	if err != nil {
		return nil, err
	}
	return t.parseMap(scans)
}

func (t Table) ListBy(objects interface{}, skip int) (int, error) {
	rv_object := reflect.Indirect(reflect.ValueOf(objects))
	if rv_object.Kind() != reflect.Slice && rv_object.Kind() != reflect.Array {
		return -1, NewErrorf("db: the objects (%s) is not an interface one of slice or array", rv_object.Kind())
	}
	take := rv_object.Len()

	rows, err := t.query(t.sqlOrderByPrimarykey+" desc limit ?,?", skip, take)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	scans := t.makeScans()
	var i int = 0
	for ; rows.Next(); i++ {
		if rv_object.Index(i).Kind() != reflect.Struct {
			return i, NewErrorf("db: the elem (%s) of array (%d) is not a struct object", rv_object.Index(i).Kind(), i)
		}
		if rv_object.Index(i).NumField() != t.ColumnNumbers {
			return i, NewErrorf("db: the elem of array (%d), numbers of field (%d) not equals numbers of column (%d)", i, rv_object.Index(i).NumField(), t.ColumnNumbers)
		}

		for j := range scans {
			scans[j] = rv_object.Index(i).Field(j).Addr().Interface()
		}

		err = rows.Scan(scans...)
		if err != nil {
			return i, err
		}
	}

	return i, rows.Err()
}

// ListArray 列出数据
func (t Table) ListSlice(take, skip int) ([][]interface{}, error) {
	rows, err := t.query(t.sqlOrderByPrimarykey+" desc limit ?,?", skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scans := t.makeNullableScans()
	datas := make([][]interface{}, 0)

	for rows.Next() {
		err = rows.Scan(scans...)
		if err != nil {
			return nil, err
		}

		a, err := t.parseSlice(scans)
		if err != nil {
			return datas, err
		}

		datas = append(datas, a)
	}

	return datas, rows.Err()
}

// List 列出数据，输出字典
func (t Table) ListMap(take, skip int) ([]map[string]interface{}, error) {
	rows, err := t.query(t.sqlOrderByPrimarykey+" desc limit ?,?", skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scans := t.makeNullableScans()
	datas := make([]map[string]interface{}, 0)

	for rows.Next() {
		err = rows.Scan(scans...)
		if err != nil {
			return nil, err
		}

		a, err := t.parseMap(scans)
		if err != nil {
			return datas, err
		}

		datas = append(datas, a)
	}

	return datas, rows.Err()
}

func (t Table) QueryBy(object interface{}, query string, args ...interface{}) (int, error) {
	rv_object := reflect.Indirect(reflect.ValueOf(object))
	if rv_object.Kind() != reflect.Slice && rv_object.Kind() != reflect.Array {
		return -1, NewErrorf("db: the objects (%s) is not an interface one of slice or array", rv_object.Kind())
	}
	take := rv_object.Len()

	rows, err := t.query(query, args...)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	scans := t.makeScans()
	var i int = 0
	for ; rows.Next() && i < take; i++ {
		if rv_object.Index(i).Kind() != reflect.Struct {
			return i, NewErrorf("db: the elem (%s) of array (%d) is not a struct object", rv_object.Index(i).Kind(), i)
		}
		if rv_object.Index(i).NumField() != t.ColumnNumbers {
			return i, NewErrorf("db: the elem of array (%d), numbers of field (%d) not equals numbers of column (%d)", i, rv_object.Index(i).NumField(), t.ColumnNumbers)
		}

		for j := range scans {
			scans[j] = rv_object.Index(i).Field(j).Addr().Interface()
		}

		err = rows.Scan(scans...)
		if err != nil {
			return i, err
		}
	}

	return i, rows.Err()
}

// Query 查询数据
func (t Table) QuerySlice(query string, args ...interface{}) ([][]interface{}, error) {
	rows, err := t.query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scans := t.makeNullableScans()
	datas := make([][]interface{}, 0)

	for rows.Next() {
		err = rows.Scan(scans...)
		if err != nil {
			return nil, err
		}

		a, err := t.parseSlice(scans)
		if err != nil {
			return datas, err
		}

		datas = append(datas, a)
	}

	return datas, rows.Err()
}

// QueryMap 查询数据，输出字典
func (t Table) QueryMap(query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := t.query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scans := t.makeNullableScans()
	datas := make([]map[string]interface{}, 0)

	for rows.Next() {
		err = rows.Scan(scans...)
		if err != nil {
			return nil, err
		}

		a, err := t.parseMap(scans)
		if err != nil {
			return datas, err
		}

		datas = append(datas, a)
	}

	return datas, rows.Err()
}

// Add 添加数据
func (t Table) Add(args ...interface{}) (int64, error) {
	return t.add(args...)
}

func (t Table) Del(key interface{}) (int64, error) {
	return t.del(key)
}

func (t Table) Set(key interface{}, args ...interface{}) (int64, error) {
	return t.set(key, args...)
}

func (t Table) Get(key interface{}, scans ...interface{}) error {
	return t.get(key, scans)
}

func (t Table) Find(args ...interface{}) *sql.Row {
	return t.find(args...)
}

func (t Table) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return t.query(query, args...)
}

func (t Table) Update(key interface{}, datas map[string]interface{}) (int64, error) {
	return t.update(key, datas)
}

func (t Table) UpdateMany(datas map[string]interface{}, query string, args ...interface{}) (int64, error) {
	return t.udpateMany(datas, query, args...)
}

// Count 统计
func (t Table) Count() int64 {
	return dbCount(t.fullName)
}

// Count 条件统计
func (t Table) CountBy(where string, args ...interface{}) int64 {
	return dbCountBy(t.fullName, where, args...)
}
