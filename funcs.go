package weiqi

import (
	"html/template"
	"time"
)

var (
	defFuncMap = template.FuncMap{
		"HasLogin": func(u *U) bool {
			return u != nil && u.Name != ""
		},
	}
)

//ParseDate 解析日期字符串
func ParseDate(dateStr string) (time.Time, error) {
	var (
		date time.Time
		err  error
	)
	date, err = time.Parse(ConstStdDate, dateStr)
	if err != nil {
		date, err = time.Parse(ConstLongDate, dateStr)
		if err != nil {
			date, err = time.Parse(ConstShortDate, dateStr)
			if err != nil {
				return time.Time{}, err
			}
		}
	}
	return date, err
}
