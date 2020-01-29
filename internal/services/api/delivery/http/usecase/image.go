package handlers

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/domens/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	ih "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api"
	delivery "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/delivery/http"
)

// ImageHandler handle requests associated with images
type ImageHandler struct {
	image   api.ImageUseCaseI
	service infrastructure.PhotoServiceI
	rep     delivery.RepositoryI

	trace infrastructure.ErrorTrace
}

func NewImageHandler(
	image api.ImageUseCaseI,
	rep delivery.RepositoryI,
	service infrastructure.PhotoServiceI,
	trace infrastructure.ErrorTrace,
) *ImageHandler {
	return &ImageHandler{
		image:   image,
		service: service,
		rep:     rep,
		trace:   trace,
	}
}

// TODO add deleting

// ! never used delete it
func (h *ImageHandler) GetImage(
	rw http.ResponseWriter,
	r *http.Request,
) ih.Result {
	var (
		err     error
		fileKey string
		url     models.Avatar
	)

	if name, _ := h.rep.GetName(r); name == "" {
		id, err := ih.GetUserIDFromAuthRequest(r)
		if err != nil {
			return ih.NewResult(
				http.StatusUnauthorized,
				nil,
				h.trace.WrapWithText(err, ErrAuth),
			)
		}
		fileKey, err = h.image.FetchByID(r.Context(), id)
	} else {
		fileKey, err = h.image.FetchByName(r.Context(), name)
	}

	if err != nil {
		return ih.NewResult(
			http.StatusNotFound,
			nil,
			h.trace.WrapWithText(err, ErrAvatarNotFound),
		)
	}

	url.URL, err = h.service.GetImage(fileKey)
	if err != nil {
		return ih.NewResult(
			http.StatusNotFound,
			nil,
			h.trace.WrapWithText(err, ErrAvatarNotFound),
		)
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
func (h *ImageHandler) PostImage(
	rw http.ResponseWriter,
	r *http.Request,
) ih.Result {
	var (
		err    error
		file   multipart.File
		userID int32
		url    models.Avatar
	)

	userID, err = ih.GetUserIDFromAuthRequest(r)
	if err != nil {
		return ih.NewResult(
			http.StatusUnauthorized,
			nil,
			h.trace.WrapWithText(err, ErrAuth),
		)
	}

	file, err = h.getFileFromRequst(rw, r)
	if err != nil {
		return ih.NewResult(
			http.StatusBadRequest,
			nil,
			h.trace.Wrap(err),
		)
	}
	defer file.Close()

	url, err = h.saveFile(r.Context(), file, userID)
	if err != nil {
		return ih.NewResult(
			http.StatusInternalServerError,
			nil,
			h.trace.Wrap(err),
		)
	}

	return ih.NewResult(http.StatusCreated, &url, nil)
}

func (h *ImageHandler) getFileFromRequst(
	rw http.ResponseWriter,
	r *http.Request,
) (multipart.File, error) {
	maxFileSize := h.service.MaxFileSize()

	r.Body = http.MaxBytesReader(rw, r.Body, maxFileSize)

	if err := r.ParseMultipartForm(maxFileSize); err != nil {
		return nil, h.trace.WrapWithText(err, ErrInvalidFile)
	}

	var (
		err    error
		file   multipart.File
		handle *multipart.FileHeader
	)

	file, handle, err = r.FormFile(FormFileName)
	if err != nil {
		return nil, h.trace.WrapWithText(err, ErrInvalidFile)
	}
	if file == nil || handle == nil {
		return nil, h.trace.New(ErrInvalidFile)
	}

	var (
		buff     bytes.Buffer
		fileSize int64
	)
	if fileSize, err = buff.ReadFrom(file); err != nil {
		file.Close()
		return nil, h.trace.WrapWithText(err, ErrInvalidFile)
	}

	if fileSize > maxFileSize {
		file.Close()
		return nil, h.trace.Errorf(
			ErrInvalidFileSize,
			fileSize,
			maxFileSize,
		)
	}

	if _, err := file.Seek(0, 0); err != nil {
		file.Close()
		return nil, h.trace.Errorf(
			ErrInvalidFileSize,
			fileSize,
			maxFileSize,
		)
	}

	err = h.checkFileType(handle.Header.Get(ContentTypeHeader))
	if err != nil {
		file.Close()
		return nil, err
	}

	return file, nil

}

func (h *ImageHandler) checkFileType(fileType string) error {
	var (
		found            = false
		allowedFileTypes = h.service.AllowedFileTypes()
	)
	for _, allowed := range allowedFileTypes {
		if fileType == allowed {
			found = true
			break
		}
	}
	if !found {
		message := []string{ErrInvalidFileFormat}
		errorText := strings.Join(append(
			message,
			strings.Join(allowedFileTypes, ","),
		), "")
		return h.trace.Errorf(errorText)
	}
	return nil
}

func (h *ImageHandler) saveFile(
	ct context.Context,
	file multipart.File,
	userID int32,
) (models.Avatar, error) {
	var (
		fileKey = uuid.NewV4().String()
		url     models.Avatar
	)

	err := h.service.SaveImage(fileKey, file)
	if err != nil {
		return url, h.trace.WrapWithText(
			err,
			ErrFailedImageSaveInService,
		)
	}

	err = h.image.Update(ct, fileKey, userID)
	if err != nil {
		h.service.DeleteImage(fileKey)
		return url, h.trace.WrapWithText(
			err,
			ErrFailedImageSaveInDatabase,
		)
	}

	url.URL, err = h.service.GetImage(fileKey)
	if err != nil {
		return url, h.trace.WrapWithText(
			err,
			ErrFailedImageSaveInService,
		)
	}

	return url, nil
}
