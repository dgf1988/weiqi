package logger

import (
	"fmt"
	"log"
	"os"
)

//Logger 日志
type Logger struct {
	Name   string
	logger *log.Logger
}

//New 创建新日志
func New(name string) *Logger {
	var l = &Logger{name, &log.Logger{}}
	l.logger.SetFlags(log.LstdFlags)
	return l
}

//SetFlags 设置
func (l *Logger) SetFlags(flag int) {
	l.logger.SetFlags(flag)
}

//SetPrefix 设置前缀
func (l *Logger) SetPrefix(prefix string) {
	l.logger.SetPrefix(prefix)
}

//Output 输出
func (l Logger) Output(s string) error {
	var err error
	var f *os.File
	if f, err = os.OpenFile(fmt.Sprint(l.Name, ".log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666); err == nil {
		defer f.Close()
		l.logger.SetOutput(f)
		return l.logger.Output(2, s)
	}
	log.Panic(err.Error())
	return err
}

//Print 打印
func (l Logger) Print(v ...interface{}) {
	l.Output(fmt.Sprint(v...))
}

//Printf 格式化
func (l Logger) Printf(format string, v ...interface{}) {
	l.Output(fmt.Sprintf(format, v...))
}

//Println 打印行
func (l Logger) Println(v ...interface{}) {
	l.Output(fmt.Sprintln(v...))
}

//Fatal 打印并退出
func (l Logger) Fatal(v ...interface{}) {
	l.Output(fmt.Sprint(v...))
	os.Exit(1)
}

//Fatalf 格式化打印并退出
func (l Logger) Fatalf(format string, v ...interface{}) {
	l.Output(fmt.Sprintf(format, v...))
	os.Exit(1)
}

//Fatalln 打印行并退出
func (l Logger) Fatalln(v ...interface{}) {
	l.Output(fmt.Sprintln(v...))
	os.Exit(1)
}

//Panic 打印并引发异常
func (l Logger) Panic(v ...interface{}) {
	var s = fmt.Sprint(v...)
	l.Output(s)
	panic(s)
}

//Panicf 格式化打印并引发异常
func (l Logger) Panicf(format string, v ...interface{}) {
	var s = fmt.Sprintf(format, v...)
	l.Output(s)
	panic(s)
}

//Panicln 打印行并引发异常
func (l Logger) Panicln(v ...interface{}) {
	var s = fmt.Sprintln(v...)
	l.Output(s)
	panic(s)
}
