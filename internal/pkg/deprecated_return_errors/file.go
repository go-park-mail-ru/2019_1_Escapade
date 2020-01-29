package rerrors

import (
	"strings"
)

// ErrorInvalidFile call it, if client give you
// 	invalid file as a request parameter
func ErrorInvalidFile() error {
	return New("Invalid file")
}

// ErrorInvalidFileFormat call it if file format isnt allowed
func ErrorInvalidFileFormat(allowedTypes []string) error {
	message := []string{"Invalid file format. You can upload an image in one of the following formats:"}

	errorText := strings.Join(append(message, strings.Join(allowedTypes, ",")), "")
	return New(errorText)
}

// ErrorInvalidFileSize call it if file size is above acceptable
func ErrorInvalidFileSize(fileSize, maxSize int64) error {
	return Errorf("Invalid file size:%v. Maximum allowed size %v", fileSize, maxSize)
}

func FileWrapper(err error) error {
	if err == nil {
		return New("Invalid file")
	}
	return New("Invalid file. More: " + err.Error())
}
