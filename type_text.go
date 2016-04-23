package weiqi

import "time"

type Text struct {
	Id     int64
	Text   string
	Status int64
	Create time.Time
	Update time.Time
}
