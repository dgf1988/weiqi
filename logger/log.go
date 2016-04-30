package logger

import (
	"os"
	"log"
	"fmt"
)

type Logger struct {
	Name string
	logger *log.Logger
}

func New(name string) *Logger {
	return &Logger{name, &log.Logger{}}
}

func (l *Logger) SetFlags(flag int) {
	l.logger.SetFlags(flag)
}

func (l *Logger) SetPrefix(prefix string) {
	l.logger.SetPrefix(prefix)
}

func (l Logger) Output(s string) error {
	if f, err := os.OpenFile(fmt.Sprint(l.Name, ".log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666); err == nil {
		defer f.Close()
		l.logger.SetOutput(f)
		return l.logger.Output(2, s)
	} else {
		panic(err.Error())
	}
}

func (l Logger) Print(v ...interface{}) {
	l.Output(fmt.Sprint(v...))
}

func (l Logger) Printf(format string, v ...interface{}) {
	l.Output(fmt.Sprintf(format, v...))
}

func (l Logger) Println(v ...interface{}) {
	l.Output(fmt.Sprintln(v...))
}

func (l Logger) Fatal(v ...interface{}) {
	l.Output(fmt.Sprint(v...))
	os.Exit(1)
}

func (l Logger) Fatalf(format string, v ...interface{}) {
	l.Output(fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (l Logger) Fatalln(v ...interface{}) {
	l.Output(fmt.Sprintln(v...))
	os.Exit(1)
}

func (l Logger) Panic(v ...interface{}) {
	var s = fmt.Sprint(v...)
	l.Output(s)
	panic(s)
}

func (l Logger) Panicf(format string, v ...interface{}) {
	var s = fmt.Sprintf(format, v...)
	l.Output(s)
	panic(s)
}

func (l Logger) Panicln(v ...interface{}) {
	var s = fmt.Sprintln(v...)
	l.Output(s)
	panic(s)
}