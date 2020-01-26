package handlers

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"

	uuid "github.com/satori/go.uuid"

	ih "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/photo"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api"
	delivery "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/delivery/http"
)

// ImageHandler handle requests associated with images
type ImageHandler struct {
	image api.ImageUseCaseI
	rep   delivery.RepositoryI
}

func NewImageHandler(image api.ImageUseCaseI, rep delivery.RepositoryI) *ImageHandler {
	handler := &ImageHandler{image, rep}
	return handler
}

// TODO add deleting
// HandleAvatar process any operation associated with user
// avatar: load and get
func (h *ImageHandler) Handle(rw http.ResponseWriter, r *http.Request) {
	ih.Route(rw, r, ih.MethodHandlers{
		http.MethodPost:    h.PostImage,
		http.MethodGet:     h.GetImage,
		http.MethodOptions: nil})
}

// delete it
func (h *ImageHandler) GetImage(rw http.ResponseWriter, r *http.Request) ih.Result {
	var (
		err     error
		fileKey string
		url     models.Avatar
	)

	if name, _ := h.rep.GetName(r); name == "" {
		id, err := ih.GetUserIDFromAuthRequest(r)
		if err != nil {
			return ih.NewResult(http.StatusBadRequest, nil, re.AuthWrapper(err))
		}
		fileKey, err = h.image.FetchByID(r.Context(), id)
	} else {
		fileKey, err = h.image.FetchByName(r.Context(), name)
	}

	if err != nil {
		return ih.NewResult(http.StatusNotFound, nil, re.NoAvatarWrapper(err))
	}

	url.URL, err = photo.GetImageFromS3(fileKey)
	if err != nil {
		return ih.NewResult(http.StatusNotFound, nil, re.NoAvatarWrapper(err))
	}

	return ih.NewResult(http.StatusOK, &url, nil)
}

// PostImage create avatar
// @Summary Create user avatar
// @Description Load new avatar to the current user. The current one is the one whose token is provided.
// @ID PostImage
// @Security OAuth2Application[write]
// @Tags account
// @Accept  multipart/form-data
// @Produce  json
// @Param file formData file true "account image"
// @Success 201 {object} models.Result "Avatar created successfully"
// @Failure 401 {object} models.Result "Required authorization"
// @Failure 500 {object} models.Result "Avatar not found"
// @Router /avatar [POST]
func (h *ImageHandler) PostImage(rw http.ResponseWriter, r *http.Request) ih.Result {
	var (
		err    error
		file   multipart.File
		userID int32
		url    models.Avatar
	)

	userID, err = ih.GetUserIDFromAuthRequest(r)
	if err != nil {
		return ih.NewResult(http.StatusUnauthorized, nil, re.AuthWrapper(err))
	}

	file, err = h.getFileFromRequst(rw, r)
	if err != nil {
		return ih.NewResult(http.StatusBadRequest, nil, err)
	}
	defer file.Close()

	url, err = h.saveFile(r.Context(), file, userID)
	if err != nil {
		return ih.NewResult(http.StatusInternalServerError, nil, err)
	}

	return ih.NewResult(http.StatusCreated, &url, nil)
}

func (h *ImageHandler) getFileFromRequst(rw http.ResponseWriter, r *http.Request) (multipart.File, error) {
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

	file, handle, err = r.FormFile("file")
	if err != nil {
		return nil, re.FileWrapper(err)
	}
	if err = re.NoNil(file, handle); err != nil {
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

func (h *ImageHandler) checkFileType(fileType string) error {
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

func (h *ImageHandler) saveFile(ct context.Context, file multipart.File, userID int32) (models.Avatar, error) {
	var (
		fileKey = uuid.NewV4().String()
		url     models.Avatar
	)

	err := photo.SaveImageInS3(fileKey, file)
	if err != nil {
		return url, re.ServerWrapper(err)
	}

	err = h.image.Update(ct, fileKey, userID)
	if err != nil {
		photo.DeleteImageFromS3(fileKey)
		return url, re.DatabaseWrapper(err)
	}

	url.URL, err = photo.GetImageFromS3(fileKey)
	if err != nil {
		return url, re.NoAvatarWrapper(err)
	}

	return url, nil
}
