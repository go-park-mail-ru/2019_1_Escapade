package api

import (
	"io/ioutil"
	"mime/multipart"
	"os"
	"strconv"

	//"escapade/internal/misc"
	//"escapade/internal/models"
	"escapade/internal/models"
	re "escapade/internal/return_errors"
)

func (h *Handler) setfiles(users []*models.UserPublicInfo) (err error) {

	for _, user := range users {
		if user.FileName == "" {
			return re.ErrorAvatarNotFound()
		}
		var filepath string
		if user.FileName == "default" {
			filepath = h.PlayersAvatarsStorage + user.FileName + "/1.png"
			user.Photo, err = ioutil.ReadFile(filepath)
		} else {
			filepath = h.PlayersAvatarsStorage + "users/" +
				strconv.Itoa(user.ID) + "/" + user.FileName
			user.Photo, err = ioutil.ReadFile(filepath)
		}
		//fmt.Println("user.Photo:" + string(user.Photo))
		if err != nil {
			return re.ErrorAvatarNotFound()
		}
	}
	return nil
}

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
