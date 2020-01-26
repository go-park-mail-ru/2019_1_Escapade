package synced

import (
	"github.com/ztrue/tracerr"
)

func Close(closers ...func() error) error {
	var err error
	for _, callClose := range closers {
		Do(callClose, &err)
	}
	return err
}

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
