package rerrors

import (
	"errors"
)

// UserInProjectNotFound godoc
func UserInProjectNotFound() error {
	return errors.New("Вы не состоите в данном проекте")
}

// UserInProjectNotFoundWrapper godoc
func UserInProjectNotFoundWrapper(err error) error {
	return errors.New("Вы не состоите в данном проекте. Ошибка:" + err.Error())
}

// ProjectNotAllowed godoc
func ProjectNotAllowed() error {
	return errors.New("У вас недостаточно прав для выполнения для данного действия")
}

// NotProjectOwner godoc
func NotProjectOwner() error {
	return errors.New("Вы не владелец проекта. Только владелец может выполнить данное действие")
}

// AlreadyInProject godoc
func AlreadyInProject() error {
	return errors.New("Вы итак являетесь участником")
}

func InternalError(msg string) error {
	return errors.New("Внутренняя ошибка. Подробнее:" + msg)
}

func NoSuchUserInProject() error {
	return errors.New("В проекте нет такого пользователя")
}

func AlreadyApplied() error {
	return errors.New("Вы уже подавали заявку в этот проект. Дождитесь, пока ее рассмотрят")
}

func AlreadyInvited() error {
	return errors.New("Вы уже пригласили данного пользователя. Дождитесь его решения")
}

func InvalidObjectType() error {
	return errors.New("тип объекта указан более 1 раза или отсутсвует(текстура, форма, изображение)")
}
