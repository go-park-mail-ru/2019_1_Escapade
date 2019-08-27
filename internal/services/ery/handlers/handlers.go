package eryhandlers

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/database"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	// erydb "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/database"

	"net/http"
)

// Handler is struct
type Handler struct {
	DB         *database.DB
	AuthClient config.AuthClient
}

// Init получить проинициализированный объект Handler
func Init(DB *database.DB, c *config.Configuration) (handler *Handler) {
	handler = &Handler{
		DB:         DB,
		AuthClient: c.AuthClient,
	}
	return
}

// Close закрыть соединие с БД
func (h *Handler) Close() error {
	return h.DB.Close()
}

// секретное слово
var salt = "NO"

func (h *Handler) HandleUser(rw http.ResponseWriter, r *http.Request) {
	var result api.Result

	utils.Debug(false, "HandleUser")

	switch r.Method {
	case http.MethodPost: //✔ зарегистрироваться
		result = h.CreateUser(rw, r)
	case http.MethodGet: // получить информацию о своем аккаунте
		result = h.GetUser(rw, r)
	case http.MethodDelete:
		// удалить свой аккаует
	case http.MethodPut: // обновить информацию о своем аккаунте
		result = h.UpdateUser(rw, r)
	case http.MethodOptions:
		return
	default:
		utils.Debug(true, "/api/user wrong request:", r.Method)
	}
	api.SendResult(rw, result)
	return
}

func (h *Handler) HandleUsers(rw http.ResponseWriter, r *http.Request) {
	var result api.Result

	switch r.Method {
	case http.MethodGet: // Получить список всех пользователей с возможностью поиска по имени
		result = h.GetUsers(rw, r)
	case http.MethodOptions:
		return
	default:
		utils.Debug(true, "/api/user wrong request:", r.Method)
	}
	api.SendResult(rw, result)
	return
}

/*
Реализовать поиск по картинкам одного пользователя
*/

// поиск по проектам...
func (h *Handler) HandleProjectsSearch(rw http.ResponseWriter, r *http.Request) {
	var result api.Result

	switch r.Method {
	case http.MethodGet: // Получить список всех проектов с возможностью поиска по имени
		result = h.projectsSearch(rw, r)
	case http.MethodOptions:
		return
	default:
		utils.Debug(true, "/api/user wrong request:", r.Method)
	}
	api.SendResult(rw, result)
	return
}

func (h *Handler) HandleUserID(rw http.ResponseWriter, r *http.Request) {
	var result api.Result

	if r.Method == http.MethodOptions {
		return
	}

	userID, err := api.IDFromPath(r, "user_id")
	if err != nil {
		api.SendResult(rw,
			api.NewResult(http.StatusBadRequest,
				"HandleProjectID", nil, err))
		return
	}

	switch r.Method {
	case http.MethodGet:
		// получить информацию о пользователе
		result = h.GetUserByID(rw, r, userID)
	default:
		utils.Debug(true, "/api/user wrong request:", r.Method)
	}
	api.SendResult(rw, result)
	return
}

func (h *Handler) HandleSession(rw http.ResponseWriter, r *http.Request) {
	var result api.Result

	switch r.Method {
	case http.MethodPost: // войти в аккаунт
		result = h.Login(rw, r)
	case http.MethodDelete: // выйти из аккаунта
		result = h.Logout(rw, r)
	case http.MethodPut: // обновить пароль
		result = h.UpdatePrivate(rw, r)
	case http.MethodOptions:
		return
	default:
		utils.Debug(true, "/api/session wrong request:", r.Method)
	}
	api.SendResult(rw, result)
	return
}

func (h *Handler) HandleProjects(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	var result api.Result

	userID, err := api.GetUserIDFromAuthRequest(r)
	if err != nil {
		api.SendResult(rw,
			api.NewResult(http.StatusUnauthorized,
				"HandleProjects", nil, re.AuthWrapper(err)))
		return
	}

	switch r.Method {
	case http.MethodPost: // создать проект
		result = h.projectsCreate(rw, r, userID)
	case http.MethodGet: // получить список проектов
		result = h.projectsGet(rw, r, userID)
	default:
		utils.Debug(true, "/api/session wrong request:", r.Method)
	}
	api.SendResult(rw, result)
	return
}

func (h *Handler) HandleProjectID(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	var result api.Result
	const place = "HandleProjectID"

	userID, err := api.GetUserIDFromAuthRequest(r)
	if err != nil {
		api.SendResult(rw,
			api.NewResult(http.StatusUnauthorized, place, nil, re.AuthWrapper(err)))
		return
	}

	projectID, err := api.IDFromPath(r, "project_id")
	if err != nil {
		api.SendResult(rw,
			api.NewResult(http.StatusBadRequest, place, nil, err))
		return
	}

	switch r.Method {
	case http.MethodPost: /* Подать заявку или принять приглашение вступить
		в проект */
		result = h.projectEnter(rw, r, projectID, userID)
	case http.MethodDelete: /* отменить заявку, отказаться от приглашения
		или покинуть проект */
		result = h.projectExit(rw, r, projectID, userID)
	case http.MethodGet: // получить информацию о проекте
		result = h.projectGet(rw, r, projectID, userID)
	case http.MethodPut: // Обновить информацию о проекте
		result = h.projectUpdate(rw, r, projectID, userID)
	case http.MethodOptions:
		return
	default:
		utils.Debug(true, "/api/session wrong request:", r.Method)
	}
	api.SendResult(rw, result)
	return
}

func (h *Handler) HandleProjectIDMembers(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	var result api.Result
	const place = "HandleProjectIDMembers"

	memberID, projectID, goalID, err := h.sceneParamsFromRequest(r)
	if err != nil {
		api.SendResult(rw, api.NewResult(http.StatusBadRequest, place, nil, err))
		return
	}

	switch r.Method {
	case http.MethodPost: // пригласить/принять заявку о вступлении в проект
		result = h.projectAcceptUser(rw, r, projectID, goalID, memberID)
	case http.MethodDelete: // исключить из проекта(в том числе отменить заявку или приглашение)
		result = h.projectKickUser(rw, r, projectID, goalID, memberID)
	case http.MethodPut: // Обновить информацию об участнике проекта
		result = h.projectUpdateUser(rw, r, projectID, goalID, memberID)
	default:
		utils.Debug(true, "/api/session wrong request:", r.Method)
	}
	api.SendResult(rw, result)
	return
}

func (h *Handler) HandleProjectIDMembersToken(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	var result api.Result
	const place = "HandleProjectIDMembersToken"

	memberID, projectID, goalID, err := h.sceneParamsFromRequest(r)
	if err != nil {
		api.SendResult(rw, api.NewResult(http.StatusBadRequest, place, nil, err))
		return
	}

	switch r.Method {
	case http.MethodPut: // Обновить права доступа участника проекта
		result = h.projectUpdateUserToken(rw, r, projectID, goalID, memberID)
	default:
		utils.Debug(true, "/api/session wrong request:", r.Method)
	}
	api.SendResult(rw, result)
	return
}

func (h *Handler) HandleScene(rw http.ResponseWriter, r *http.Request) {
	var result api.Result

	switch r.Method {
	case http.MethodPost: // Создать новую сцену
		result = h.sceneCreate(rw, r)
	case http.MethodOptions:
		return
	default:
		utils.Debug(true, "/api/session wrong request:", r.Method)
	}
	api.SendResult(rw, result)
	return
}

func (h *Handler) HandleSceneID(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	var result api.Result
	const place = "HandleSceneID"

	userID, projectID, sceneID, err := h.sceneParamsFromRequest(r)
	if err != nil {
		api.SendResult(rw, api.NewResult(http.StatusBadRequest, place, nil, err))
		return
	}

	switch r.Method {
	case http.MethodGet: // Получить сцену со всеми объектами
		result = h.sceneObjectsGet(rw, r, userID, projectID, sceneID)
	case http.MethodPut: // Обновить сцену
		result = h.sceneUpdate(rw, r, userID, projectID, sceneID)
	case http.MethodDelete: // Удалить сцену
		result = h.sceneDelete(rw, r, userID, projectID, sceneID)
	default:
		utils.Debug(true, "/api/session wrong request:", r.Method)
	}
	api.SendResult(rw, result)
	return
}

func (h *Handler) HandleSceneErythrocyte(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	var result api.Result
	const place = "HandleSceneID"

	userID, projectID, sceneID, err := h.sceneParamsFromRequest(r)
	if err != nil {
		api.SendResult(rw, api.NewResult(http.StatusBadRequest, place, nil, err))
		return
	}

	switch r.Method {
	case http.MethodPost: // Создать новый эритроцит
		result = h.erythrocyteCreate(rw, r, userID, projectID, sceneID)
	case http.MethodPut: // Обновить эритроцит
		result = h.erythrocyteUpdate(rw, r, userID, projectID, sceneID)
	case http.MethodDelete: // Удалить эритроцит
		result = h.erythrocyteDelete(rw, r, userID, projectID, sceneID)
	default:
		utils.Debug(true, "/api/session wrong request:", r.Method)
	}
	api.SendResult(rw, result)
	return
}

func (h *Handler) HandleSceneErythrocyteObject(rw http.ResponseWriter, r *http.Request) {
	var result api.Result

	const place = "HandleSceneErythrocyteObject"

	userID, projectID, sceneID, err := h.sceneParamsFromRequest(r)
	if err != nil {
		api.SendResult(rw, api.NewResult(http.StatusBadRequest, place, nil, err))
		return
	}

	switch r.Method {
	case http.MethodPost: // Создать новую модель/текстуру/снимок
		result = h.eryobjectCreate(rw, r, userID, projectID, sceneID)
	case http.MethodPut: // обновить модель/текстуру/снимок
		result = h.eryobjectUpdate(rw, r, userID, projectID, sceneID)
	case http.MethodDelete: // удалить модель/текстуру/снимок
		result = h.eryobjectDelete(rw, r, userID, projectID, sceneID)
	case http.MethodOptions:
		return
	default:
		utils.Debug(true, "/api/session wrong request:", r.Method)
	}
	api.SendResult(rw, result)
	return
}

func (h *Handler) HandleSceneDisease(rw http.ResponseWriter, r *http.Request) {
	var result api.Result

	const place = "HandleSceneDisease"

	userID, projectID, sceneID, err := h.sceneParamsFromRequest(r)
	if err != nil {
		api.SendResult(rw, api.NewResult(http.StatusBadRequest, place, nil, err))
		return
	}

	switch r.Method {
	case http.MethodPost: // Создать новую болезнь
		result = h.diseaseCreate(rw, r, userID, projectID, sceneID)
	case http.MethodPut: // Обновить болезнь
		result = h.diseaseUpdate(rw, r, userID, projectID, sceneID)
	case http.MethodDelete: // Удалить болезнь
		result = h.diseaseDelete(rw, r, userID, projectID, sceneID)
	case http.MethodOptions:
		return
	default:
		utils.Debug(true, "/api/session wrong request:", r.Method)
	}
	api.SendResult(rw, result)
	return
}

func (h *Handler) sceneParamsFromRequest(r *http.Request) (int32, int32, int32, error) {
	var (
		userID, projectID, sceneID int32
		err                        error
	)
	userID, err = api.GetUserIDFromAuthRequest(r)
	if err != nil {
		return userID, projectID, sceneID, re.AuthWrapper(err)
	}

	projectID, err = api.IDFromPath(r, "project_id")
	if err != nil {
		return userID, projectID, sceneID, err
	}

	sceneID, err = api.IDFromPath(r, "scene_id")

	return userID, projectID, sceneID, err
}

func (h *Handler) userProjectParamsFromRequest(r *http.Request) (int32, int32, int32, error) {
	var (
		memberID, projectID, goalID int32
		err                         error
	)
	memberID, err = api.GetUserIDFromAuthRequest(r)
	if err != nil {
		return memberID, projectID, goalID, re.AuthWrapper(err)
	}

	projectID, err = api.IDFromPath(r, "project_id")
	if err != nil {
		return memberID, projectID, goalID, err
	}

	goalID, err = api.IDFromPath(r, "user_id")

	return memberID, projectID, goalID, err
}
