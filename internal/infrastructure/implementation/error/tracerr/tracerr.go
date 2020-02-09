package tracerr

import (
	"github.com/ztrue/tracerr"
)

// Tracerr is implementation of ErrorTrace interface(package infrastructure)
// via ztrue/tracerr
type Tracerr struct{}

// New instance of Tracerr
func New() *Tracerr {
	return new(Tracerr)
}

func (*Tracerr) Wrap(err error) error {
	if err == nil {
		return err
	}
	return tracerr.Wrap(err)
}

func (tr *Tracerr) WrapWithText(
	err error,
	text string,
) error {
	if err == nil {
		return err
	}
	return tracerr.Wrap(tr.New(text + " " + err.Error()))
}

func (*Tracerr) New(message string) error {
	return tracerr.New(message)
}

func (*Tracerr) Errorf(
	message string,
	args ...interface{},
) error {
	return tracerr.Errorf(message, args...)
}
