package api

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/photo"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"

	"mime/multipart"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

// GetImage returns user avatar
// @Summary Get user avatar
// @Description Get user avatar
// @ID GetImage
// @Success 200 {object} models.Result "Avatar found successfully"
// @Failure 401 {object} models.Result "Required authorization"
// @Failure 404 {object} models.Result "Avatar not found"
// @Router /avatar [GET]
func (h *Handler) GetImage(rw http.ResponseWriter, r *http.Request) {
	const place = "GetImage"
	var (
		err     error
		userID  int
		fileKey string
		url     models.Avatar
	)

	if userID, err = h.getUserIDFromCookie(r, h.Session); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		utils.SendErrorJSON(rw, re.ErrorAuthorization(), place)
		utils.PrintResult(err, http.StatusUnauthorized, place)
		return
	}

	if fileKey, err = h.DB.GetImage(userID); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		utils.SendErrorJSON(rw, re.ErrorAvatarNotFound(), place)
		utils.PrintResult(err, http.StatusNotFound, place)
		return
	}

	url.URL, err = photo.GetImage(fileKey)
	if err != nil {
		utils.Debug(false, "Failed to sign request", err)
		rw.WriteHeader(http.StatusNotFound)
		utils.SendErrorJSON(rw, re.ErrorAvatarNotFound(), place)
		utils.PrintResult(err, http.StatusNotFound, place)
		return
	}

	utils.SendSuccessJSON(rw, url, place)
	utils.PrintResult(err, http.StatusOK, place)
}

// PostImage create avatar
// @Summary Create user avatar
// @Description Create user avatar
// @ID PostImage
// @Success 201 {object} models.Result "Avatar created successfully"
// @Failure 401 {object} models.Result "Required authorization"
// @Failure 500 {object} models.Result "Avatar not found"
// @Router /avatar [POST]
func (h *Handler) PostImage(rw http.ResponseWriter, r *http.Request) {
	const place = "PostImage"

	var (
		err    error
		input  multipart.File
		userID int
		handle *multipart.FileHeader
		url    models.Avatar
	)

	if userID, err = h.getUserIDFromCookie(r, h.Session); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		utils.SendErrorJSON(rw, re.ErrorAuthorization(), place)
		utils.PrintResult(err, http.StatusUnauthorized, place)
		return
	}

	if input, handle, err = r.FormFile("file"); err != nil || input == nil || handle == nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.SendErrorJSON(rw, re.ErrorInvalidFile(), place)
		utils.PrintResult(err, http.StatusInternalServerError, place)
		return
	}

	defer input.Close()

	fileType := handle.Header.Get("Content-Type")
	//Генерация уник.ключа для хранения картинки
	fileKey := uuid.NewV4()

	switch fileType {
	case "image/jpeg":
		err = photo.SaveImage(fileKey.String(), input)
	case "image/png":
		err = photo.SaveImage(fileKey.String(), input)
	default:
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, re.ErrorInvalidFileFormat(), place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.SendErrorJSON(rw, re.ErrorServer(), place)
		utils.PrintResult(err, http.StatusInternalServerError, place)
		return
	}

	if err = h.DB.PostImage(fileKey.String(), userID); err != nil {
		photo.DeleteImage(fileKey.String())
		rw.WriteHeader(http.StatusInternalServerError)
		utils.SendErrorJSON(rw, re.ErrorDataBase(), place)
		utils.PrintResult(err, http.StatusInternalServerError, place)
		return
	}

	if url.URL, err = photo.GetImage(fileKey.String()); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		utils.SendErrorJSON(rw, re.ErrorAvatarNotFound(), place)
		utils.PrintResult(err, http.StatusNotFound, place)
		return
	}

	utils.SendSuccessJSON(rw, url, place)
	utils.PrintResult(err, http.StatusCreated, place)
}
