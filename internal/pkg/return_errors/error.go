package rerrors

import (
	"github.com/ztrue/tracerr"
)

type Closer interface {
	Close() error
}

func Close(closers ...Closer) error {
	var err error
	for _, closer := range closers {
		Do(closer.Close, &err)
	}
	return err
}

// func CatchErrors(actions ...func() error) error {
// 	var err error
// 	for _, action := range actions {
// 		Do(action, &err)
// 	}
// 	return err
// }

func Do(action func() error, updateErr *error) {
	err := action()
	if err != nil {
		if *updateErr != nil {
			*updateErr = tracerr.Wrap(err)
		} else {
			*updateErr = err
		}
	}
}

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
