package infrastructure

type Logger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

type LoggerNil struct{}

func (*LoggerNil) Fatal(v ...interface{})                 {}
func (*LoggerNil) Fatalf(format string, v ...interface{}) {}
func (*LoggerNil) Print(v ...interface{})                 {}
func (*LoggerNil) Println(v ...interface{})               {}
func (*LoggerNil) Printf(format string, v ...interface{}) {}
