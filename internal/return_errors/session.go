package rerrors

import "errors"

// ErrorSessionQueryNotOK redis cant execute query
func ErrorSessionQueryNotOK(msg string) error {
	return errors.New("Redis query result not ok, it is:"+msg)
}
