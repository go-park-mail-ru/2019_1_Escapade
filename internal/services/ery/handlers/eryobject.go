package eryhandlers

import (
	"bytes"
	"net/http"
	"mime/multipart"

	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/photo"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/models"

)

func (h *Handler) eryobjectCreate(rw http.ResponseWriter, r *http.Request,
	userID, projectID, sceneID int32) api.Result {
	const place = "eryobjectCreate"

	var (
		file   multipart.File
		handle *multipart.FileHeader
		err    error
	)

	maxFileSize := int64(6000000000)

	if err := r.ParseMultipartForm(maxFileSize); err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, re.FileWrapper(err))
	}

	r.Body = http.MaxBytesReader(rw, r.Body, maxFileSize)

	if file, handle, err = r.FormFile("File"); err != nil || file == nil || handle == nil {
		return api.NewResult(http.StatusBadRequest, place, nil, re.FileWrapper(err))
	}

	defer file.Close()

	var (
		buff     bytes.Buffer
		fileSize int64
	)

	if fileSize, err = buff.ReadFrom(file); err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, re.FileWrapper(err))
	}

	if fileSize > maxFileSize {
		return api.NewResult(http.StatusBadRequest, place, nil,
			re.ErrorInvalidFileSize(fileSize, maxFileSize))
	}

	if _, err := file.Seek(0, 0); err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil,
			re.ErrorInvalidFileSize(fileSize, maxFileSize))
	}
	/*
		var (
			fileType         = handle.Header.Get("Content-Type")
			found            = false
			allowedFileTypes = make([]string, 0)
		)
		utils.Debug(false, "type:", fileType)

		allowedFileTypes = append(allowedFileTypes, "image/jpg", "image/jpeg", "image/png", "image/gif", "obj", "")
		for _, allowed := range allowedFileTypes {
			if fileType == allowed {
				found = true
				break
			}
		}
		if !found {
			return api.NewResult(http.StatusBadRequest, place, nil, re.ErrorInvalidFileFormat(allowedFileTypes))
		}*/

	utils.Debug(false, "Set ")
	//fileKey := uuid.NewV4().String()

	var eryOBJ models.EryObject
	err = eryOBJ.Set(r, handle.Filename, utils.String32(int32(fileSize)), "artyom/"+handle.Filename)
	if err != nil {
		return api.NewResult(http.StatusBadRequest, place, nil, err)
	}

	utils.Debug(false, "SaveImageInS3")
	eryOBJ.Path = "artyom/" + handle.Filename
	err = photo.SaveImageInS3(eryOBJ.Path, file)
	if err != nil {
		return api.NewResult(http.StatusInternalServerError, place, nil, re.ServerWrapper(err))
	}

	utils.Debug(false, "CreateEryObject")
	err = h.DB.CreateEryObject(userID, projectID, sceneID, &eryOBJ)
	if err != nil {
		photo.DeleteImageFromS3(eryOBJ.Path)
		return api.NewResult(http.StatusInternalServerError, place, nil, re.DatabaseWrapper(err))
	}
	eryOBJ.Path, _ = photo.GetImageFromS3(eryOBJ.Path)
	return api.NewResult(http.StatusCreated, place, &eryOBJ, nil)
}

func (h *Handler) eryobjectUpdate(rw http.ResponseWriter, r *http.Request,
	userID, projectID, sceneID, objectID int32) api.Result {
	return api.UpdateModel(r, &models.EryObjectUpdate{}, "eryobjectUpdate", false,
		func(userID int32) (api.JSONtype, error) {
			return h.DB.GetEryObject(objectID)
		},
		func(interf api.JSONtype) error {
			eryObject, ok := interf.(*models.EryObject)
			if !ok {
				return re.NoUpdate()
			}
			return h.DB.UpdateEryObject(userID, projectID, *eryObject)
		})
}

func (h *Handler) eryobjectDelete(rw http.ResponseWriter, r *http.Request,
	userID, projectID, sceneID, objectID int32) api.Result {
	const place = "erythrocyteDelete"

	err := h.DB.DeleteEryObject(userID, projectID, objectID)
	if err != nil {
		return api.NewResult(http.StatusInternalServerError, place, nil, err)
	}

	return api.NewResult(http.StatusCreated, place, nil, err)
}
