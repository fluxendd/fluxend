package controllers

import (
	"fluxton/requests/bucket_requests"
	"fluxton/resources"
	"fluxton/responses"
	"fluxton/services"
	"fluxton/utils"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type FileController struct {
	fileService services.FileService
}

func NewFileController(injector *do.Injector) (*FileController, error) {
	fileService := do.MustInvoke[services.FileService](injector)

	return &FileController{fileService: fileService}, nil
}

func (fc *FileController) List(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	bucketUUID, err := utils.GetUUIDPathParam(c, "bucketUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, "Invalid bucket UUID")
	}

	paginationParams := utils.ExtractPaginationParams(c)
	files, err := fc.fileService.List(paginationParams, bucketUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FileResourceCollection(files))
}

func (fc *FileController) Show(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	bucketUUID, err := utils.GetUUIDPathParam(c, "bucketUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	fileUUID, err := utils.GetUUIDPathParam(c, "fileUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	file, err := fc.fileService.GetByUUID(fileUUID, bucketUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FileResource(&file))
}

func (fc *FileController) Store(c echo.Context) error {
	var request bucket_requests.CreateFileRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	bucketUUID, err := utils.GetUUIDPathParam(c, "bucketUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	file, err := fc.fileService.Create(bucketUUID, &request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.FileResource(&file))
}

func (fc *FileController) Rename(c echo.Context) error {
	var request bucket_requests.RenameFileRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	fileUUID, err := utils.GetUUIDPathParam(c, "fileUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	bucketUUID, err := utils.GetUUIDPathParam(c, "bucketUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	updatedFile, err := fc.fileService.Rename(fileUUID, bucketUUID, authUser, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FileResource(updatedFile))
}

func (fc *FileController) Delete(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	bucketUUID, err := utils.GetUUIDPathParam(c, "bucketUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	fileUUID, err := utils.GetUUIDPathParam(c, "fileUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if _, err := fc.fileService.Delete(fileUUID, bucketUUID, authUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}
