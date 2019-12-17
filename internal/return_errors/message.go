package rerrors

import "errors"

func ErrorMessageInvalidID() error {
	return errors.New("Invalid message id")
}
