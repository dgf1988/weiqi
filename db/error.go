package db

import "fmt"

var (
	errNilPtr = newErrorf("db: destination pointer is nil")
)

type typeError struct {
	msg string
}

func newErrorf(format string, args ...interface{}) typeError {
	return typeError{msg: fmt.Sprintf(format, args...)}
}

func (e typeError) Error() string {
	return e.msg
}
