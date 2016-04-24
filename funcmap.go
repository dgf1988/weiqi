package weiqi

import (
	"html/template"
)

var (
	defFuncMap = template.FuncMap{
		"HasLogin": func(u *User) bool {
			return u != nil && u.Name != ""
		},
	}
)
