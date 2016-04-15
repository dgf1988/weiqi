package main

import (
	"log"
	"os"
)

const (
	Log_File_Error  = "err.log"
	Log_File_Access = "access.log"
	Log_File_Debug  = "debug.log"
)

func logError(h *Http, msg string) {
	ferr, err := os.OpenFile(Log_File_Error, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	defer ferr.Close()
	log.New(ferr, "", log.LstdFlags).Println(h.R.RemoteAddr, h.R.Method, h.R.Host, h.R.URL.String(), msg)
}

func logAccess(h *Http) {
	fout, err := os.OpenFile(Log_File_Access, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	defer fout.Close()
	log.New(fout, "", log.LstdFlags).Println(
		h.R.RemoteAddr, h.R.Method, h.R.Host, h.R.URL.String(), h.R.UserAgent())
}

func logDebug(a ...interface{}) {
	fout, err := os.OpenFile(Log_File_Debug, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	defer fout.Close()
	log.New(fout, "", log.LstdFlags).Println(a...)
}
