package rerrors

import (
	"errors"
	"fmt"
	"strings"
)

// ErrorInvalidFile call it, if client give you
// 	invalid file as a request parameter
func ErrorInvalidFile() error {
	return errors.New("Invalid file")
}

// ErrorInvalidFileFormat call it if file format isnt allowed
func ErrorInvalidFileFormat(allowedTypes []string) error {
	message := []string{"Invalid file format. You can upload an image in one of the following formats:"}

	errorText := strings.Join(append(message, strings.Join(allowedTypes, ",")), "")
	return errors.New(errorText)
}

// ErrorInvalidFileSize call it if file size is above acceptable
func ErrorInvalidFileSize(fileSize, maxSize int64) error {
	message := fmt.Sprintf("Invalid file size:%v. Maximum allowed size %v", fileSize, maxSize)

	return errors.New(message)
}
