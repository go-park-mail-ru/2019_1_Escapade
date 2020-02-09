package infrastructure

import (
	"errors"
	"fmt"
)

type ErrorTrace interface {
	Wrap(err error) error
	WrapWithText(err error, text string) error
	New(message string) error
	Errorf(message string, args ...interface{}) error
}

type ErrorTraceNil struct{}

func (*ErrorTraceNil) Wrap(err error) error {
	return err
}
func (*ErrorTraceNil) WrapWithText(
	err error,
	text string,
) error {
	return err
}
func (*ErrorTraceNil) New(message string) error {
	return errors.New(message)
}
func (*ErrorTraceNil) Errorf(
	message string,
	args ...interface{},
) error {
	return fmt.Errorf(message, args)
}
