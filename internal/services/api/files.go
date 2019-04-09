package api

import (
	"io/ioutil"
	"mime/multipart"
	"os"
)

func saveFile(path string, name string, file multipart.File, mode os.FileMode) (err error) {
	var (
		data []byte
	)

	os.MkdirAll(path, mode)

	if data, err = ioutil.ReadAll(file); err != nil {
		return
	}

	if err = ioutil.WriteFile(path+"/"+name, data, mode); err != nil {
		return
	}

	return
}
