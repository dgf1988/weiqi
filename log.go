package weiqi

import (
	"os"
	"log"
)

const (
	c_errfilename = "err.log"
)


func logError(format string, args ...interface{}) {
	if f, err := os.OpenFile(c_errfilename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666); err == nil {
		defer f.Close()
		log.New(f, "[WeiqiError]", log.LstdFlags).Printf(format, args...)
	} else {
		panic(err.Error())
	}
}
