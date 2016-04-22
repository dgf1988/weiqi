package db

import (
	"reflect"
	"time"
)

func parseValue(src interface{}) interface{} {
	if s, ok := src.(Nullable); ok {
		src, _ = s.Value()
	}
	if src == nil {
		return nil
	}
	return reflect.Indirect(reflect.ValueOf(src)).Interface()
}

func copyValue(dest interface{}, src interface{}) error {
	if s, ok := src.(Nullable); ok {
		src, _ = s.Value()
	}
	switch s := src.(type) {
	//int64
	case *int64 :
		switch d := dest.(type) {
		case *int64:
			if d == nil {
				return ErrNilPtr
			}
			*d = *s
			return nil
		}
	case int64:
		switch d := dest.(type) {
		case *int64:
			if d == nil {
				return ErrNilPtr
			}
			*d = s
			return nil
		}
	case *float64:
		switch d := dest.(type) {
		case *float64 :
			if d == nil {
				return ErrNilPtr
			}
			*d = *s
			return nil
		}
	case float64:
		switch d := dest.(type) {
		case *float64:
			if d == nil {
				return ErrNilPtr
			}
			*d = s
			return nil
		}
	case *bool :
		switch d := dest.(type) {
		case *bool :
			if d == nil {
				return ErrNilPtr
			}
			*d = *s
			return nil
		}
	case bool :
		switch d := dest.(type) {
		case *bool :
			if d == nil {
				return ErrNilPtr
			}
			*d = s
			return nil
		}
	case *string:
		switch d := dest.(type) {
		case *string:
			if d == nil {
				return ErrNilPtr
			}
			*d = *s
			return nil
		}
	case string:
		switch d := dest.(type) {
		case *string:
			if d == nil {
				return ErrNilPtr
			}
			*d = s
			return nil
		}
	case *[]byte :
		switch d := dest.(type) {
		case *[]byte :
			if d == nil {
				return ErrNilPtr
			}
			*d = *s
			return nil
		}
	case []byte :
		switch d := dest.(type) {
		case *[]byte:
			if d == nil {
				return ErrNilPtr
			}
			*d = s
			return nil
		}
	case *time.Time:
		switch d := dest.(type) {
		case *time.Time:
			if d == nil {
				return ErrNilPtr
			}
			*d = *s
			return nil
		}
	case time.Time:
		switch d := dest.(type) {
		case *time.Time:
			if d == nil {
				return ErrNilPtr
			}
			*d = s
			return nil
		}
	}
	return NewErrorf("db: type error %s => %s", reflect.TypeOf(src), reflect.TypeOf(dest))
}
