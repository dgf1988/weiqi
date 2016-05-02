package db

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func parseValue(src interface{}) interface{} {
	if s, ok := src.(driver.Valuer); ok {
		src, _ = s.Value()
	}
	if src == nil {
		return nil
	}
	return reflect.Indirect(reflect.ValueOf(src)).Interface()
}

func convertValue(dest interface{}, src interface{}) error {
	if s, ok := src.(driver.Valuer); ok {
		src, _ = s.Value()
	}
	if d, ok := dest.(sql.Scanner); ok {
		return d.Scan(src)
	}
	switch s := src.(type) {
	case *int64:
		return convertValue(dest, *s)
	case *bool:
		return convertValue(dest, *s)
	case *float64:
		return convertValue(dest, *s)
	case *string:
		return convertValue(dest, *s)
	case *time.Time:
		return convertValue(dest, *s)
	case *[]byte:
		return convertValue(dest, *s)
	//int64
	case int64:
		switch d := dest.(type) {
		case *int64:
			if d == nil {
				return errNilPtr
			}
			*d = s
			return nil
		case *string:
			if d == nil {
				return errNilPtr
			}
			*d = fmt.Sprint(s)
			return nil
		case *bool:
			if d == nil {
				return errNilPtr
			}
			if s == 0 {
				*d = false
				return nil
			} else if s == 1 {
				*d = true
				return nil
			}
			return errors.New(fmt.Sprintf("db: the int64(%v) can't convert value to bool.", s))
		case *float64:
			if d == nil {
				return errNilPtr
			}
			*d = float64(s)
			return nil
		}
	case float64:
		switch d := dest.(type) {
		case *float64:
			if d == nil {
				return errNilPtr
			}
			*d = s
			return nil
		case *string:
			if d == nil {
				return errNilPtr
			}
			*d = fmt.Sprint(s)
			return nil
		case *bool:
			if d == nil {
				return errNilPtr
			}
			if s == 0.0 {
				*d = false
				return nil
			} else if s == 1.0 {
				*d = true
				return nil
			}
			return errors.New(fmt.Sprintf("db: the float64(%v) can't convert value to bool.", s))
		}
	case bool:
		switch d := dest.(type) {
		case *bool:
			if d == nil {
				return errNilPtr
			}
			*d = s
			return nil
		case *string:
			if d == nil {
				return errNilPtr
			}
			*d = fmt.Sprint(s)
			return nil
		case *float64:
			if d == nil {
				return errNilPtr
			}
			if s {
				*d = 1.0
			} else {
				*d = 0.0
			}
			return nil
		case *int64:
			if d == nil {
				return errNilPtr
			}
			if s {
				*d = 1
			} else {
				*d = 0
			}
			return nil
		}
	case string:
		switch d := dest.(type) {
		case *string:
			if d == nil {
				return errNilPtr
			}
			*d = s
			return nil
		case *int64:
			if d == nil {
				return errNilPtr
			}
			value, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return err
			} else {
				*d = value
				return nil
			}
		case *float64:
			if d == nil {
				return errNilPtr
			}
			value, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return err
			} else {
				*d = value
				return nil
			}
		case *bool:
			if d == nil {
				return errNilPtr
			}
			value, err := strconv.ParseBool(s)
			if err != nil {
				return err
			} else {
				*d = value
				return nil
			}
		case *time.Time:
			if d == nil {
				return errNilPtr
			}
			value, err := time.Parse("2006-01-02 15:04:05", s)
			if err != nil {
				value, err = time.Parse("2006-01-02", s)
				if err != nil {
					return err
				}
			}
			*d = value
			return nil
		}
	case []byte:
		switch d := dest.(type) {
		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = s
			return nil
		default:
			return convertValue(dest, string(s))
		}
	case time.Time:
		switch d := dest.(type) {
		case *string:
			if d == nil {
				return errNilPtr
			}
			*d = s.Format("2006-01-02 15:04:05")
			return nil
		case *time.Time:
			if d == nil {
				return errNilPtr
			}
			*d = s
			return nil
		}
	}
	return newErrorf("db: convertValue: type error: %T(%v) => %T", src, src, dest)
}
