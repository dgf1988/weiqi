package weiqi

import (
	"errors"
	"fmt"
)

var (
	ErrInputEmpty = errors.New("输入为空")
)

type WeiqiError struct {
	Msg string
}

func NewWeiqiError(msg string) *WeiqiError {
	return &WeiqiError{msg}
}

func (e WeiqiError) Error() string {
	return fmt.Sprint("weiqi: ", e.Msg)
}
