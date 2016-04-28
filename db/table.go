package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
)

func GetTable(databasename, tablename string) (Table, error) {
	cols, err := getColumns(databasename, tablename)
	if err != nil {
		return nil, err
	}
	var t = newTable()
	//设置表名
	t.DatabaseName = databasename
	t.Name = tablename
	t.Fullname = fmt.Sprint(t.DatabaseName, ".", t.Name)

	//保存字段
	t.Columns = cols
	t.ColumnNumbers = len(cols)
	if t.ColumnNumbers <= 0 {
		return nil, newErrorf("db: table not found")
	}

	//遍历字段，读取其它信息。
	keys := make([]string, 0)
	for _, col := range cols {
		keys = append(keys, col.FullName)
		t.sqlArgMark = append(t.sqlArgMark, "?")
		if col.Key == "PRI" {
			t.Primarykey = col.Name
		} else if col.Key == "UNI" {
			t.UniqueIndex = append(t.UniqueIndex, col.Name)
		}
	}

	strKeys := strings.Join(keys, ", ")

	//保存预备SQL语句。
	t.sqlInsert = fmt.Sprintf("INSERT INTO %s", t.Fullname)
	t.sqlDelete = fmt.Sprintf("DELETE FROM %s", t.Fullname)
	t.sqlUpdate = fmt.Sprintf("UPDATE %s", t.Fullname)
	t.sqlSelect = fmt.Sprintf("SELECT %s FROM %s ", strKeys, t.Fullname)
	t.sqlSelectCount = fmt.Sprintf("SELECT COUNT(%s) FROM %s", t.Primarykey, t.Fullname)

	return t, nil
}

// ToSql 输出表结构Sql语句。
func (t typeTable) ToSql() string {
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

// Add 添加数据
func (t typeTable) Add(values ...interface{}) (int64, error) {
	listcolname := make([]string, 0)
	listParam := make([]interface{}, 0)
	for i := range values {
		if values[i] == nil {
			continue
		}
		listcolname = append(listcolname, t.Columns[i].FullName)
		listParam = append(listParam, values[i])
	}
	res, err := dbExec(fmt.Sprintf("%s (%s) VALUES (%s)", t.sqlInsert, strings.Join(listcolname, ", "), strings.Join(t.sqlArgMark[:len(listParam)], ", ")), listParam...)
	if err != nil {
		return -1, err
	}
	return res.LastInsertId()
}

func (t typeTable) Del(args ...interface{}) (int64, error) {
	listwhere := make([]string, 0)
	listparam := make([]interface{}, 0)
	for i := range args {
		if args[i] == nil {
			continue
		}
		listwhere = append(listwhere, t.Columns[i].FullName+"=?")
		listparam = append(listparam, args[i])
	}

	res, err := dbExec(fmt.Sprintf("%s WHERE %s limit 1", t.sqlDelete, strings.Join(listwhere, " AND ")), listparam...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (t *typeTable) Update(args ...interface{}) Set {
	listwhere := make([]string, 0)
	listparam := make([]interface{}, 0)
	for i := range args {
		if args[i] == nil {
			continue
		}
		listwhere = append(listwhere, t.Columns[i].FullName+"=?")
		listparam = append(listparam, args[i])
	}
	query := fmt.Sprintf("WHERE %s limit 1", strings.Join(listwhere, " AND "))
	return &typeSetter{
		t, query, listparam,
	}
}

func (t *typeTable) Get(args ...interface{}) Row {
	listwhere := make([]string, 0)
	listparam := make([]interface{}, 0)
	for i := range args {
		if args[i] == nil {
			continue
		}
		listwhere = append(listwhere, t.Columns[i].FullName+"=?")
		listparam = append(listparam, args[i])
	}
	strSql := fmt.Sprintf("%s WHERE %s limit 1", t.sqlSelect, strings.Join(listwhere, " AND "))
	return &typeRow{
		dbQueryRow(strSql, listparam...),
		t,
	}
}

func (t *typeTable) Find(args ...interface{}) (Rows, error) {
	listwhere := make([]string, 0)
	listparam := make([]interface{}, 0)
	for i := range args {
		if args[i] == nil {
			continue
		}
		listwhere = append(listwhere, t.Columns[i].FullName+"=?")
		listparam = append(listparam, args[i])
	}
	strSql := fmt.Sprintf("%s WHERE %s", t.sqlSelect, strings.Join(listwhere, " AND "))
	rows, err := dbQuery(strSql, listparam...)
	if err != nil {
		return nil, err
	}
	return &typeRows{
		rows, t, t.makeNullableScans(),
	}, nil
}

func (t *typeTable) Any(args ...interface{}) (Rows, error) {
	listwhere := make([]string, 0)
	listparam := make([]interface{}, 0)
	for i := range args {
		if args[i] == nil {
			continue
		}
		listwhere = append(listwhere, t.Columns[i].FullName+"=?")
		listparam = append(listparam, args[i])
	}
	strSql := fmt.Sprintf("%s WHERE %s", t.sqlSelect, strings.Join(listwhere, " OR "))
	rows, err := dbQuery(strSql, listparam...)
	if err != nil {
		return nil, err
	}
	return &typeRows{
		rows, t, t.makeNullableScans(),
	}, nil
}

func (t *typeTable) List(take, skip int) (Rows, error) {
	rows, err := dbQuery(fmt.Sprintf("%s limit ?, ?", t.sqlSelect), skip, take)
	if err != nil {
		return nil, err
	}
	return &typeRows{
		rows, t, t.makeNullableScans(),
	}, nil
}

func (t *typeTable) ListDesc(take, skip int) (Rows, error) {
	rows, err := dbQuery(fmt.Sprintf("%s ORDER BY %s DESC limit ?, ?", t.sqlSelect, t.Primarykey), skip, take)
	if err != nil {
		return nil, err
	}
	return &typeRows{
		rows, t, t.makeNullableScans(),
	}, nil
}

func (t *typeTable) Query(query string, args ...interface{}) (Rows, error) {
	strSql := fmt.Sprintf("%s %s", t.sqlSelect, query)
	rows, err := dbQuery(strSql, args...)
	if err != nil {
		return nil, err
	}
	return &typeRows{
		rows, t, t.makeNullableScans(),
	}, nil
}

// Count 统计
func (t typeTable) Count(query string, args ...interface{}) (int64, error) {
	var num int64
	err := dbQueryRow(fmt.Sprintf("%s %s", t.sqlSelectCount, query), args...).Scan(&num)
	if err != nil {
		return -1, err
	}
	return num, nil
}

// Table 保存表信息
type typeTable struct {
	//数据库名
	DatabaseName string
	//表名
	Name string
	//字段结构信息
	Columns []typeColumn
	//字段数量
	ColumnNumbers int
	//主键
	Primarykey string
	//唯一键
	UniqueIndex []string

	Fullname string
	// 预备Sql执行语句
	sqlInsert string

	sqlSelect string

	sqlDelete string
	sqlUpdate string

	sqlSelectCount string
	sqlArgMark     []string
}

func newTable() *typeTable {
	return &typeTable{
		Columns: make([]typeColumn, 0), UniqueIndex: make([]string, 0), sqlArgMark: make([]string, 0),
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
			scans[i] = new(nullTime)
		case typeChar, typeVarchar, typeText, typeMediumTtext, typeLongtext:
			scans[i] = new(sql.NullString)
		case typeFloat, typeDouble, typeDecimal:
			scans[i] = new(sql.NullFloat64)

		default:
			scans[i] = new(nullBytes)
		}
	}
	return scans
}

func (t typeTable) makeStructScans(object interface{}) ([]interface{}, error) {
	scans := make([]interface{}, t.ColumnNumbers)
	rv := reflect.ValueOf(object)
	if rv.Kind() != reflect.Ptr {
		return nil, newErrorf("db: the object (%s) is not a pointer", rv.Kind())
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return nil, newErrorf("db: the pointer (%s) can't point to a struct object", rv.Kind())
	}
	if rv.NumField() != t.ColumnNumbers {
		return nil, newErrorf("db: the object field numbers (%d) not equals table column numbers (%d)", rv.NumField(), t.ColumnNumbers)
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

/*

func (typeTable) sqlQuery(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(query, args...)
}

func (typeTable) sqlQueryRow(query string, args ...interface{}) *sql.Row {
	return db.QueryRow(query, args...)
}

func (typeTable) sqlExec(query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
}

*/
