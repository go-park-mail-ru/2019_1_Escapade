package log

import "log"

type Log struct{}

func New() *Log {
	return new(Log)
}

func (*Log) Fatal(v ...interface{}) {
	log.Fatal(v...)
}
func (*Log) Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}
func (*Log) Print(v ...interface{}) {
	log.Print(v...)
}
func (*Log) Println(v ...interface{}) {
	log.Println(v...)
}
func (*Log) Printf(format string, v ...interface{}) {
	log.Println(v...)
}
