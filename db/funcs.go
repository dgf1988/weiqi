package db

import (
	"fmt"
	"strings"
)

func copyMap(src map[string]interface{}) (dest map[string]interface{}) {
	dest = make(map[string]interface{})
	for k, v := range src {
		dest[k] = v
	}
	return
}

func formatMapToInsert(kvs map[string]interface{}) (string, []interface{}) {
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

func formatMapToSet(kvs map[string]interface{}) (string, []interface{}) {
	sqlitems := make([]string, 0)
	values := make([]interface{}, 0)
	for k, v := range kvs {
		sqlitems = append(sqlitems, fmt.Sprint(k, "=?"))
		values = append(values, v)
	}
	return strings.Join(sqlitems, ","), values
}

func formatMapToWhere(kvs map[string]interface{}) (string, []interface{}) {
	sqlitems := make([]string, 0)
	values := make([]interface{}, 0)
	for k, v := range kvs {
		if v == nil {
			sqlitems = append(sqlitems, fmt.Sprint(k, " is null"))
			continue
		}
		sqlitems = append(sqlitems, fmt.Sprint(k, "=?"))
		values = append(values, v)
	}
	return strings.Join(sqlitems, " AND "), values
}
