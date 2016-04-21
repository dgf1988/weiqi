package db

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

// Type 保存数据库里的类型信息
type Type struct {
	Name   string
	Value  int64
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
	typebuf, ok := v.([]byte)
	if !ok {
		typebuf, ok = v.([]uint8)
		if !ok {
			return errors.New(fmt.Sprintf("db: Type Scan Error: not accept type in %v", v))
		}
	}
	if len(typebuf) == 0 {
		return errors.New("db: no type name from scanner")
	}
	t.Name = string(typebuf)
	switch t.Name {
	//int64
	case "int":
		t.Value = typeInt

	//string
	case "char":
		t.Value = typeChar
	case "varchar":
		t.Value = typeVarchar
	case "text":
		t.Value = typeText
	case "mediumtext":
		t.Value = typeMediumTtext
	case "longtext":
		t.Value = typeLongtext

	//time.Time
	case "date":
		t.Value = typeDate
	case "datetime":
		t.Value = typeDatetime
	case "year":
		t.Value = typeYear
	case "timestamp":
		t.Value = typeTimestamp
	case "time":
		t.Value = typeTime

	case "float":
		t.Value = typeFloat

	//Error
	default:
		return errors.New(fmt.Sprintf("db: not supported type %v", t.Name))
	}
	return nil
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

	Order int

	Default Default

	Nullable bool

	Type Type

	Key     string
	Extra   string
	Comment string

	IsAuto bool
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
		DATA_TYPE,
		COLUMN_KEY,	EXTRA, COLUMN_COMMENT
	FROM
		information_schema.COLUMNS
	WHERE
		TABLE_SCHEMA = ? AND TABLE_NAME = ?
	ORDER BY
		 ORDINAL_POSITION
	`
	rows, err := dbQuery(sqlquery, databasename, tablename)
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
		col.FullName = fmt.Sprintf("%s.`%s`", col.TableName, col.Name)

		col.Nullable = scannullable == "YES"
		col.IsAuto = col.Default.V == "CURRENT_TIMESTAMP" || col.Extra == "auto_increment"

		cols = append(cols, col)
	}
	return cols, rows.Err()
}
