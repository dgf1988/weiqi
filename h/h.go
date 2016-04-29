package h

import (
	"github.com/dgf1988/weiqi/logger"
	"log"
)


var (
	errorlogger = logger.Logger{"error", &log.Logger{}}
	accesslogger = logger.Logger{"access", &log.Logger{}}
	spiderlogger = logger.Logger{"spider", &log.Logger{}}
)

func init() {
	errorlogger.SetPrefix("[Error: 500]")
	errorlogger.SetFlags(log.LstdFlags)

	accesslogger.SetPrefix("[Access]")
	accesslogger.SetFlags(log.LstdFlags)

	spiderlogger.SetPrefix("[Spider]")
	spiderlogger.SetFlags(log.LstdFlags)
}