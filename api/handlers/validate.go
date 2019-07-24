package api

import (
	"strings"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
)

func validateUser(user *models.UserPrivateInfo) error {
	name := strings.TrimSpace(user.Name)
	if name == "" || len(name) < 3 {
		return re.ErrorInvalidName()
	}
	user.Name = name

	password := strings.TrimSpace(user.Password)
	if len(password) < 3 {
		return re.ErrorInvalidPassword()
	}
	return nil
}
