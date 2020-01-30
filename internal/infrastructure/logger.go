package infrastructure

type LoggerI interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

type LoggerEmpty struct{}
func (*LoggerEmpty) Fatal(v ...interface{}) {}
func (*LoggerEmpty) 	Fatalf(format string, v ...interface{}) {}
func (*LoggerEmpty) 	Print(v ...interface{}) {}
func (*LoggerEmpty) 	Println(v ...interface{}) {}
func (*LoggerEmpty) 	Printf(format string, v ...interface{}) {}