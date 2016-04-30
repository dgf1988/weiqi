package weiqi

import (
	"log"
	"github.com/dgf1988/weiqi/logger"
)
var (
	errlogger = logger.New("weiqierror")
)

func init() {
	errlogger.SetFlags(log.LstdFlags)
	errlogger.SetPrefix("[WeiqiError]")
}

func logError(format string, args ...interface{}) {
	errlogger.Printf(format, args...)
}
