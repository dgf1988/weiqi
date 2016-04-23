package db

import (
	"database/sql/driver"
	"time"
)

type iNullable interface {
	Scan(value interface{}) error
	Value() (driver.Value, error)
}

// NullTime 可空时间结构体
type nullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *nullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt nullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

type nullBytes struct {
	Bytes []byte
	Valid bool
}

func (nb *nullBytes) Scan(value interface{}) error {
	nb.Bytes, nb.Valid = value.([]byte)
	return nil
}

func (nb nullBytes) Value() (driver.Value, error) {
	if !nb.Valid {
		return nil, nil
	}
	return nb.Bytes, nil
}
