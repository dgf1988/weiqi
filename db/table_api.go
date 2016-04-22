package db

import (
	"fmt"
	"strings"
)


func GetTable(databasename, tablename string) (ITable, error) {
	cols, err := GetColumns(databasename, tablename)
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


// Count 统计
func (t typeTable) Count(query string, args ...interface{}) (int64, error) {
	var num int64
	err := dbQueryRow(fmt.Sprintf("%s %s", t.sqlSelectCount, query), args...).Scan(&num)
	if err != nil {
		return -1, err
	}
	return num, nil
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
		listwhere = append(listwhere, t.Columns[i].FullName + "=?")
		listparam = append(listparam, args[i])
	}

	res, err := dbExec(fmt.Sprintf("%s WHERE %s limit 1", t.sqlDelete, strings.Join(listwhere, " AND ")), listparam...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (t *typeTable) Get(args ...interface{}) IRow {
	listwhere := make([]string, 0)
	listparam := make([]interface{}, 0)
	for i := range args {
		if args[i] == nil {
			continue
		}
		listwhere = append(listwhere, t.Columns[i].FullName + "=?")
		listparam = append(listparam, args[i])
	}
	strSql := fmt.Sprintf("%s WHERE %s limit 1", t.sqlSelect, strings.Join(listwhere, " AND "))
	return &typeRow{
		dbQueryRow(strSql, listparam...),
		t,
	}
}

func (t *typeTable) Find(args ...interface{}) ISet {
	listwhere := make([]string, 0)
	listparam := make([]interface{}, 0)
	for i := range args {
		if args[i] == nil {
			continue
		}
		listwhere = append(listwhere, t.Columns[i].FullName + "=?")
		listparam = append(listparam, args[i])
	}
	query := fmt.Sprintf("WHERE %s limit 1", strings.Join(listwhere, " AND "))
	return &Setter{
		t, query, listparam,
	}
}

func (t *typeTable) List(take, skip int) (IRows, error) {
	rows, err := dbQuery(fmt.Sprintf("%s limit ?, ?", t.sqlSelect), skip, take)
	if err != nil {
		return nil, err
	}
	return &typeRows{
		rows, t, t.makeNullableScans(),
	}, nil
}

func (t *typeTable) Query(query string, args ...interface{}) (IRows, error) {
	strSql := fmt.Sprintf("%s %s", t.sqlSelect, query)
	rows, err := dbQuery(strSql, args...)
	if err != nil {
		return nil, err
	}
	return &typeRows{
		rows, t, t.makeNullableScans(),
	}, nil
}