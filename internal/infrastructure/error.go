package infrastructure

type ErrorTrace interface {
	Wrap(err error) error
	WrapWithText(err error, text string) error
	New(message string) error
	Errorf(message string, args ...interface{}) error
}
