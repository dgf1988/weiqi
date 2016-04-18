package db

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"database/sql"
	"time"
)

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}
// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}
// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

type Type struct {
	Name   string
	Length int
}

func (t Type) ToSql() string {
	switch t.Name {
	case "text", "mediumtext", "date", "datetime", "timestamp":
		return t.Name
	}
	return fmt.Sprintf("%s(%d)", t.Name, t.Length)
}

func (t *Type) Scan(v interface{}) error {
	str := string(v.([]uint8))
	var (
		strtype string
		strnum  string
		intnum  int
		ii      int
	)
	for i, ch := range str {
		if ch == '(' {
			strtype = str[:i]
			ii = i + 1
		}
		if ch == ')' {
			strnum = str[ii:i]
		}
	}
	intnum, _ = strconv.Atoi(strnum)
	if strtype == "" {
		t.Name, t.Length = str, 0
	} else {
		t.Name, t.Length = strtype, intnum
	}
	return nil
}

func (t Type) Value() (driver.Value, error) {
	return t.ToSql(), nil
}

type Default struct {
	Null             bool
	V                string
	CurrentTimestamp bool
}

func (d Default) ToSql() string {
	if !d.Null {
		if d.CurrentTimestamp {
			return "DEFAULT CURRENT_TIMESTAMP"
		} else {
			return fmt.Sprintf("DEFAULT '%s'", d.V)
		}
	} else {
		return "DEFAULT NULL"
	}
}

func (d *Default) Scan(v interface{}) error {
	if v == nil {
		d.Null, d.V, d.CurrentTimestamp = true, "", false
	} else {
		d.Null = false
		d.V = string(v.([]uint8))
		if d.V == "CURRENT_TIMESTAMP" {
			d.CurrentTimestamp = true
		}
	}
	return nil
}

func (d Default) Value() (driver.Value, error) {
	if d.Null {
		return d.V, nil
	} else {
		return nil, nil
	}
}

type Column struct {
	DatabaseName string
	TableName    string
	Name         string
	Order        int

	Default Default

	Nullable bool

	Type Type

	Key     string
	Extra   string
	Comment string
}

func (c Column) ToSql() string {
	stritems := make([]string, 0)
	stritems = append(stritems, fmt.Sprintf("`%s`", c.Name), c.Type.ToSql())
	if c.Nullable {
		stritems = append(stritems, "NULL", c.Default.ToSql())
	} else {
		stritems = append(stritems, "NOT NULL")
		if !c.Default.Null {
			stritems = append(stritems, c.Default.ToSql())
		}
	}
	stritems = append(stritems, c.Extra)
	return strings.Join(stritems, " ")
}

func GetColumns(databasename, tablename string) ([]Column, error) {
	sqlquery := `
	SELECT
		TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ORDINAL_POSITION,
		COLUMN_DEFAULT, IS_NULLABLE,
		COLUMN_TYPE,
		COLUMN_KEY,	EXTRA, COLUMN_COMMENT
	FROM
		information_schema.COLUMNS
	WHERE
		TABLE_SCHEMA = ? AND TABLE_NAME = ?
	ORDER BY
		 ORDINAL_POSITION
	`
	rows, err := query(sqlquery, databasename, tablename)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols := make([]Column, 0)
	for rows.Next() {
		var (
			col          Column
			scannullable string
		)
		err = rows.Scan(&col.DatabaseName, &col.TableName, &col.Name, &col.Order,
			&col.Default, &scannullable, &col.Type,
			&col.Key, &col.Extra, &col.Comment,
		)
		if err != nil {
			return nil, err
		}
		col.Nullable = scannullable == "YES"
		cols = append(cols, col)
	}
	return cols, rows.Err()
}

type Table struct {
	DatabaseName string
	Name         string
	Columns      []Column
	Length 		 int
	Primarykey   string
	UniqueIndex  []string
}

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

func (t Table) Get(key interface{}) ([]interface{}, error) {
	row := queryRow(fmt.Sprintf("select * from %s.%s where %s = ? limit 1", t.DatabaseName, t.Name, t.Primarykey), key)
	dest := make([]interface{}, t.Length)
	for i := range t.Columns {
		switch t.Columns[i].Type.Name {
		case "int":
			dest[i] = new(sql.NullInt64)
		case "date", "year", "datetime", "time", "timestamp":
			dest[i] = new(NullTime)
		default:
			dest[i] = new(sql.NullString)
		}
	}
	err := row.Scan(dest...)
	if err != nil {
		return nil, err
	}
	data := make([]interface{}, 0)
	for i := range dest {
		switch dest[i].(type) {
		case *sql.NullString:
			nullstr := dest[i].(*sql.NullString)
			if nullstr.Valid {
				data = append(data, nullstr.String)
			} else {
				data = append(data, nil)
			}
		case *sql.NullInt64:
			nullint := dest[i].(*sql.NullInt64)
			if nullint.Valid {
				data = append(data, nullint.Int64)
			} else {
				data = append(data, nil)
			}
		case *NullTime:
			nulltime := dest[i].(*NullTime)
			if nulltime.Valid {
				data = append(data, nulltime.Time)
			} else {
				data = append(data, nil)
			}
		}
	}
	return data, nil
}

func (t Table) List(take, skip int) ([][]interface{}, error) {
	rows, err := query(fmt.Sprintf("select * from %s.%s order by %s limit ?,?", t.DatabaseName, t.Name, t.Primarykey), skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dest := make([]interface{}, t.Length)
	for i := range t.Columns {
		switch t.Columns[i].Type.Name {
		case "int":
			dest[i] = new(sql.NullInt64)
		case "date", "year", "datetime", "time", "timestamp":
			dest[i] = new(NullTime)
		default:
			dest[i] = new(sql.NullString)
		}
	}
	results := make([][]interface{}, 0)
	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			return nil, err
		}
		data := make([]interface{}, 0)
		for i := range dest {
			switch dest[i].(type) {
			case *sql.NullString:
				nullstr := dest[i].(*sql.NullString)
				if nullstr.Valid {
					data = append(data, nullstr.String)
				} else {
					data = append(data, nil)
				}
			case *sql.NullInt64:
				nullint := dest[i].(*sql.NullInt64)
				if nullint.Valid {
					data = append(data, nullint.Int64)
				} else {
					data = append(data, nil)
				}
			case *NullTime:
				nulltime := dest[i].(*NullTime)
				if nulltime.Valid {
					data = append(data, nulltime.Time)
				} else {
					data = append(data, nil)
				}
			}
		}
		results = append(results, data)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
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
