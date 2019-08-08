package api

import (
	"bytes"
	"fmt"

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
		name    string
		fileKey string
		url     models.Avatar
	)

	if name, err = h.getName(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, re.ErrorInvalidName(), place)
		utils.PrintResult(err, http.StatusNotFound, place)
	}

	if fileKey, err = h.DB.GetImage(name); err != nil {
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
		file   multipart.File
		userID int32
		handle *multipart.FileHeader
		url    models.Avatar
	)

	if userID, err = h.getUserIDFromAuthRequest(r); err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		utils.SendErrorJSON(rw, re.ErrorAuthorization(), place)
		utils.PrintResult(err, http.StatusUnauthorized, place)
		return
	}

	utils.Debug(false, "r.FormFile next")
	maxFileSize := photo.MaxFileSize()

	if err := r.ParseMultipartForm(maxFileSize); err != nil {
		utils.Debug(false, "ParseMultipartForm err", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, re.ErrorInvalidFile(), place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	r.Body = http.MaxBytesReader(rw, r.Body, maxFileSize)

	if file, handle, err = r.FormFile("file"); err != nil || file == nil || handle == nil {
		utils.Debug(false, "r.FormFile err", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, re.ErrorInvalidFile(), place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	defer file.Close()

	var buff bytes.Buffer
	fileSize, err := buff.ReadFrom(file)
	fmt.Println("fileSize:", fileSize)

	if err != nil {
		utils.Debug(false, "ReadFrom err", err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, re.ErrorInvalidFile(), place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	utils.Debug(false, "compare sizes", fileSize, maxFileSize)
	if fileSize > maxFileSize {
		err = re.ErrorInvalidFileSize(fileSize, maxFileSize)
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, err, place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}

	if _, err := file.Seek(0, 0); err != nil {
		utils.Debug(false, "file.Seek err", err.Error())
		return
	}

	var (
		fileType         = handle.Header.Get("Content-Type")
		found            = false
		allowedFileTypes = photo.AllowedFileTypes()
	)
	for _, allowed := range allowedFileTypes {
		if fileType == allowed {
			found = true
			break
		}
	}
	if !found {
		utils.Debug(false, "found type", fileType)
		rw.WriteHeader(http.StatusBadRequest)
		utils.SendErrorJSON(rw, re.ErrorInvalidFileFormat(allowedFileTypes), place)
		utils.PrintResult(err, http.StatusBadRequest, place)
		return
	}
	//Генерация уник.ключа для хранения картинки
	fileKey := uuid.NewV4()

	utils.Debug(false, "save next", fileType)

	err = photo.SaveImage(fileKey.String(), file)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.SendErrorJSON(rw, re.ErrorServer(), place)
		utils.PrintResult(err, http.StatusInternalServerError, place)
		return
	}

	utils.Debug(false, "h.DB.PostImage next")

	if err = h.DB.PostImage(fileKey.String(), userID); err != nil {
		photo.DeleteImage(fileKey.String())
		rw.WriteHeader(http.StatusInternalServerError)
		utils.SendErrorJSON(rw, re.ErrorDataBase(), place)
		utils.PrintResult(err, http.StatusInternalServerError, place)
		return
	}

	utils.Debug(false, "photo.GetImage next")

	if url.URL, err = photo.GetImage(fileKey.String()); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		utils.SendErrorJSON(rw, re.ErrorAvatarNotFound(), place)
		utils.PrintResult(err, http.StatusNotFound, place)
		return
	}

	utils.SendSuccessJSON(rw, url, place)
	utils.PrintResult(err, http.StatusCreated, place)
}
