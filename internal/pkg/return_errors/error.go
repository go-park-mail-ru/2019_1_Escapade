package rerrors

import (
	"github.com/ztrue/tracerr"
)

func Wrap(err error) error {
	if err == nil {
		return err
	}
	return tracerr.Wrap(err)
}

func New(message string) error {
	return tracerr.New(message)
}

func Errorf(message string, args ...interface{}) error {
	return tracerr.Errorf(message, args...)
}

func NoNil(args ...interface{}) error {
	for _, arg := range args {
		if arg == nil {
			return InterfaceIsNil()
		}
	}
	return nil
}
