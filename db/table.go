package db

import (
	"fmt"
	"strings"
	"database/sql"
)

// Table 保存表信息
type Table struct {
	DatabaseName string
	Name         string
	Columns      []Column
	Length 		 int
	Primarykey   string
	UniqueIndex  []string
}

// ToSql 输出表结构Sql语句。
func (t Table) ToSql() string {
	stritems := make([]string, 0)
	stritems = append(stritems, fmt.Sprintf("CREATE TABLE `%s` (", t.Name))
	colitems :=  make([]string, 0)
	for i := range t.Columns {
		colitems = append(colitems, "\t" + t.Columns[i].ToSql())
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

func (t Table) getScans() []interface{} {
	scans := make([]interface{}, t.Length)
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

// Get 按主键取数据
func (t Table) Get(key interface{}) ([]interface{}, error) {
	row := queryRow(fmt.Sprintf("select * from %s.%s where %s = ? limit 1", t.DatabaseName, t.Name, t.Primarykey), key)
	dest := t.getScans()
	err := row.Scan(dest...)
	if err != nil {
		return nil, err
	}
	return t.parseScansArray(dest), nil
}

// GetMap	按主键取数据，输出字典
func (t Table) GetMap(key interface{}) (map[string]interface{}, error) {
	row := queryRow(fmt.Sprintf("select * from %s.%s where %s = ? limit 1", t.DatabaseName, t.Name, t.Primarykey), key)
	dest := t.getScans()
	err := row.Scan(dest...)
	if err != nil {
		return nil, err
	}
	return t.parseScansMap(dest), nil
}

// Find 查找数据
func (t Table) Find(where string, args ...interface{}) ([]interface{}, error) {
	row := queryRow(fmt.Sprintf("select * from %s.%s where %s limit 1", t.DatabaseName, t.Name, where), args...)
	scans := t.getScans()
	err := row.Scan(scans...)
	if err != nil {
		return nil, err
	}
	return t.parseScansArray(scans), nil
}

// FindMap	查找数据，输出字典
func (t Table) FindMap(where string, args ...interface{}) (map[string]interface{}, error) {
	row := queryRow(fmt.Sprintf("select * from %s.%s where %s limit 1", t.DatabaseName, t.Name, where), args...)
	scans := t.getScans()
	err := row.Scan(scans...)
	if err != nil {
		return nil, err
	}
	return t.parseScansMap(scans), nil
}

// List 列出数据
func (t Table) List(take, skip int) ([][]interface{}, error) {
	rows, err := query(fmt.Sprintf("select * from %s.%s order by %s limit ?,?", t.DatabaseName, t.Name, t.Primarykey), skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dest := t.getScans()
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
	rows, err := query(fmt.Sprintf("select * from %s.%s order by %s limit ?,?", t.DatabaseName, t.Name, t.Primarykey), skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dest := t.getScans()
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
func (t Table) Query(take, skip int, where string, v ...interface{}) ([][]interface{}, error) {
	rows, err := query(fmt.Sprintf("select * from %s.%s where %s order by %s limit %d,%d", t.DatabaseName, t.Name, where, t.Primarykey, skip, take), v...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dest := t.getScans()
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
func (t Table) QueryMap(take, skip int, where string, v ...interface{}) ([]map[string]interface{}, error) {
	rows, err := query(fmt.Sprintf("select * from %s.%s where %s order by %s limit %d,%d", t.DatabaseName, t.Name, where, t.Primarykey, skip, take), v...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dest := t.getScans()
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
	sets, values := formatMapToSets(kvs)
	res, err := exec(fmt.Sprintf("update %s.%s set %s where %s='%v' limit 1", t.DatabaseName, t.Name, sets, t.Primarykey, key), values...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

// SetMany 批量修改数据
func (t Table) SetMany(where string, items []interface{}, kvs map[string]interface{}) (int64, error) {
	sets, values := formatMapToSets(kvs)
	for i := range items {
		values = append(values, items[i])
	}
	res, err := exec(fmt.Sprintf("update %s.%s set %s where %s", t.DatabaseName, t.Name, sets, where), values...)
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

func formatMapToInsert(kvs  map[string]interface{}) (string, []interface{}) {
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

func formatMapToSets(kvs map[string]interface{}) (string, []interface{}) {
	sqlitems := make([]string, 0)
	values := make([]interface{}, 0)
	for k, v := range kvs {
		sqlitems = append(sqlitems, fmt.Sprint(k, "=?"))
		values = append(values, v)
	}
	return strings.Join(sqlitems, ","), values
}

func newTable() *Table {
	return &Table{
		Columns:make([]Column, 0), UniqueIndex:make([]string, 0),
	}
}

func GetTable(databasename, tablename string) (*Table, error) {
	cols, err := GetColumns(databasename, tablename)
	if err != nil {
		return nil, err
	}
	var table = newTable()
	table.DatabaseName = databasename
	table.Name = tablename
	table.Columns = cols
	table.Length = len(cols)
	for _, col := range cols {
		if col.Key == "PRI" {
			table.Primarykey = col.Name
		} else if col.Key == "UNI" {
			table.UniqueIndex = append(table.UniqueIndex, col.Name)
		}
	}
	return table, nil
}

