package h

import (
	"github.com/dgf1988/weiqi/logger"
	"log"
)

var (
	errorlogger    = logger.New("error")
	accesslogger   = logger.New("access")
	spiderlogger   = logger.New("spider")
	notfoundlogger = logger.New("notfound")
)

func init() {
	errorlogger.SetPrefix("[Error: 500]")
	errorlogger.SetFlags(log.LstdFlags)

	accesslogger.SetPrefix("[Access]")
	accesslogger.SetFlags(log.LstdFlags)

	spiderlogger.SetPrefix("[Spider]")
	spiderlogger.SetFlags(log.LstdFlags)

	notfoundlogger.SetPrefix("[Notfound: 404]")
	notfoundlogger.SetFlags(log.LstdFlags)
}
