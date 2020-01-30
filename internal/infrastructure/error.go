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

type ErrorTraceDefault struct{}

func (*ErrorTraceDefault) Wrap(err error) error {
	return err
}
func (*ErrorTraceDefault) WrapWithText(
	err error,
	text string,
) error {
	return err
}
func (*ErrorTraceDefault) New(message string) error {
	return errors.New(message)
}
func (*ErrorTraceDefault) Errorf(
	message string,
	args ...interface{},
) error {
	return fmt.Errorf(message, args)
}
