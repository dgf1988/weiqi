package db

import (
	"fmt"
	"strings"
)

type IRow interface {
	Scan(dest ...interface{}) error
}

type IRows interface {
	Close() error
	Err() error
	Next() bool
	Scan(dest ...interface{}) error
}

type ITable interface {
	Add(args ...interface{}) (int64, error)
	Del(key interface{}) (int64, error)

	Set(key interface{}, args ...interface{}) error
	Get(key interface{}, dest ...interface{}) error

	Query(query string, args ...interface{}) (IRows, error)
	Update(datas map[string]interface{}, query string, args ...interface{}) (int64, error)
}

type Row struct {
	t *Table
	query string
	args []interface{}
}

func (r Row) Scan(scans ...interface{}) error {
	keys := make([]string, 0)
	dest := make([]interface{}, 0)
	for i := range scans {
		if scans[i] == nil {
			continue
		}
		keys = append(keys, r.t.Columns[i].FullName)
		dest = append(dest, scans[i])
	}
	return dbQueryRow(fmt.Sprintf("SELECT %s FROM %s %s", strings.Join(keys, ", "), r.t.Fullname, r.query), r.args...).Scan(dest...)
}