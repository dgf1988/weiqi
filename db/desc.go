package db

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"reflect"
	"time"
)

// Type 保存数据库里的类型信息
type Type struct {
	Name   string
	Type	reflect.Type
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
	switch t.Name {
	case "char", "varchar", "text", "mediumtext":
		t.Type = reflect.TypeOf(t.Name)
	case "date", "datetime", "year", "timestamp":
		t.Type = reflect.TypeOf(time.Time{})
	case "int":
		t.Type = reflect.TypeOf(t.Length)
	default:
		panic("db: unknow column's type")
	}
	return nil
}

func (t Type) Value() (driver.Value, error) {
	return t.ToSql(), nil
}

// Default 保存数据库里的默认值信息
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

// Column 保存行结构信息
type Column struct {
	DatabaseName string
	TableName    string
	Name         string
	FullName     string

	Order        int

	Default Default

	Nullable bool

	Type 	Type

	Key     string
	Extra   string
	Comment string

	IsAuto  bool
}

// ToSql 输出行结构
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

// GetColumns 取某表的行结构
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
		col.FullName = fmt.Sprint(col.TableName, ".", col.Name)

		col.Nullable = scannullable == "YES"
		col.IsAuto = col.Default.V == "CURRENT_TIMESTAMP" || col.Extra == "auto_increment"

		cols = append(cols, col)
	}
	return cols, rows.Err()
}
