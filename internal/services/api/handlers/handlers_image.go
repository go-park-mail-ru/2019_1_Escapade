package handlers

import (
	"bytes"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/photo"

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
		fileKey, err = h.Db.image.FetchByID(id)
	} else {
		fileKey, err = h.Db.image.FetchByName(name)
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
		url    models.Avatar
	)

	if userID, err = ih.GetUserIDFromAuthRequest(r); err != nil {
		return ih.NewResult(http.StatusUnauthorized, place, nil, re.AuthWrapper(err))
	}

	file, err = h.getFileFromRequst(rw, r)
	if err != nil {
		return ih.NewResult(http.StatusBadRequest, place, nil, err)
	}
	defer file.Close()

	url, err = h.saveFile(file, userID)
	if err != nil {
		return ih.NewResult(http.StatusInternalServerError, place, nil, err)
	}

	return ih.NewResult(http.StatusCreated, place, &url, nil)
}

func (h *Handler) getFileFromRequst(rw http.ResponseWriter, r *http.Request) (multipart.File, error) {
	maxFileSize := photo.MaxFileSize()

	r.Body = http.MaxBytesReader(rw, r.Body, maxFileSize)

	if err := r.ParseMultipartForm(maxFileSize); err != nil {
		return nil, re.FileWrapper(err)
	}

	var (
		err    error
		file   multipart.File
		handle *multipart.FileHeader
	)

	if file, handle, err = r.FormFile("file"); err != nil || file == nil || handle == nil {
		return nil, re.FileWrapper(err)
	}

	var (
		buff     bytes.Buffer
		fileSize int64
	)
	if fileSize, err = buff.ReadFrom(file); err != nil {
		file.Close()
		return nil, re.FileWrapper(err)
	}

	if fileSize > maxFileSize {
		file.Close()
		return nil, re.ErrorInvalidFileSize(fileSize, maxFileSize)
	}

	if _, err := file.Seek(0, 0); err != nil {
		file.Close()
		return nil, re.ErrorInvalidFileSize(fileSize, maxFileSize)
	}

	err = h.checkFileType(handle.Header.Get("Content-Type"))
	if err != nil {
		file.Close()
		return nil, err
	}

	return file, nil

}

func (h *Handler) checkFileType(fileType string) error {
	var (
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
		return re.ErrorInvalidFileFormat(allowedFileTypes)
	}
	return nil
}

func (h *Handler) saveFile(file multipart.File, userID int32) (models.Avatar, error) {
	var (
		fileKey = uuid.NewV4().String()
		url     models.Avatar
	)

	err := photo.SaveImageInS3(fileKey, file)
	if err != nil {
		return url, re.ServerWrapper(err)
	}

	if err = h.Db.image.Update(fileKey, userID); err != nil {
		photo.DeleteImageFromS3(fileKey)
		return url, re.DatabaseWrapper(err)
	}

	if url.URL, err = photo.GetImageFromS3(fileKey); err != nil {
		return url, re.NoAvatarWrapper(err)
	}

	return url, nil
}

// 192 -> 141 -> 180
