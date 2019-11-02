package handlers

import (
	"bytes"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/photo"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	ih "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
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
func (h *Handler) GetImage(rw http.ResponseWriter, r *http.Request) ih.Result {
	const place = "GetImage"
	var (
		err     error
		fileKey string
		url     models.Avatar
	)

	if name, _ := h.getName(r); name == "" {
		id, err := ih.GetUserIDFromAuthRequest(r)
		if err != nil {
			return ih.NewResult(http.StatusBadRequest, place, nil, re.AuthWrapper(err))
		}
		fileKey, err = h.DB.GetImageByID(id)
	} else {
		fileKey, err = h.DB.GetImageByName(name)
	}

	if err != nil {
		return ih.NewResult(http.StatusNotFound, place, nil, re.NoAvatarWrapper(err))
	}

	url.URL, err = photo.GetImageFromS3(fileKey)
	if err != nil {
		return ih.NewResult(http.StatusNotFound, place, nil, re.NoAvatarWrapper(err))
	}

	return ih.NewResult(http.StatusOK, place, &url, nil)
}

// PostImage create avatar поделать курл
// @Summary Create user avatar
// @Description Create user avatar
// @ID PostImage
// @Success 201 {object} models.Result "Avatar created successfully"
// @Failure 401 {object} models.Result "Required authorization"
// @Failure 500 {object} models.Result "Avatar not found"
// @Router /avatar [POST]
func (h *Handler) PostImage(rw http.ResponseWriter, r *http.Request) ih.Result {
	const place = "PostImage"

	var (
		err    error
		file   multipart.File
		userID int32
		handle *multipart.FileHeader
		url    models.Avatar
	)

	if userID, err = ih.GetUserIDFromAuthRequest(r); err != nil {
		return ih.NewResult(http.StatusUnauthorized, place, nil, re.AuthWrapper(err))
	}

	maxFileSize := photo.MaxFileSize()

	r.Body = http.MaxBytesReader(rw, r.Body, maxFileSize)

	utils.Debug(false, "ParseMultipartForm")
	if err := r.ParseMultipartForm(maxFileSize); err != nil {
		return ih.NewResult(http.StatusBadRequest, place, nil, re.FileWrapper(err))
	}

	utils.Debug(false, "FormFile")
	if file, handle, err = r.FormFile("file"); err != nil || file == nil || handle == nil {
		return ih.NewResult(http.StatusBadRequest, place, nil, re.FileWrapper(err))
	}

	defer file.Close()

	var (
		buff     bytes.Buffer
		fileSize int64
	)
	utils.Debug(false, "ReadFrom")
	if fileSize, err = buff.ReadFrom(file); err != nil {
		return ih.NewResult(http.StatusBadRequest, place, nil, re.FileWrapper(err))
	}

	if fileSize > maxFileSize {
		return ih.NewResult(http.StatusBadRequest, place, nil,
			re.ErrorInvalidFileSize(fileSize, maxFileSize))
	}

	utils.Debug(false, "Seek")
	if _, err := file.Seek(0, 0); err != nil {
		return ih.NewResult(http.StatusBadRequest, place, nil,
			re.ErrorInvalidFileSize(fileSize, maxFileSize))
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
		return ih.NewResult(http.StatusBadRequest, place, nil, re.ErrorInvalidFileFormat(allowedFileTypes))
	}

	fileKey := uuid.NewV4().String()

	err = photo.SaveImageInS3(fileKey, file)
	if err != nil {
		return ih.NewResult(http.StatusInternalServerError, place, nil, re.ServerWrapper(err))
	}

	if err = h.DB.PostImage(fileKey, userID); err != nil {
		photo.DeleteImageFromS3(fileKey)
		return ih.NewResult(http.StatusInternalServerError, place, nil, re.DatabaseWrapper(err))
	}

	if url.URL, err = photo.GetImageFromS3(fileKey); err != nil {
		return ih.NewResult(http.StatusInternalServerError, place, nil, re.NoAvatarWrapper(err))
	}

	return ih.NewResult(http.StatusCreated, place, &url, nil)
}

// 192 -> 141
