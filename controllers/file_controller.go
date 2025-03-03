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

// List retrieves all files in a bucket
//
// @Summary List all files in a bucket
// @Description Retrieve a list of all files in a specific bucket
// @Tags files
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param bucketUUID path string true "Bucket UUID"
//
// @Success 200 {array} responses.Response{content=[]resources.FileResponse} "List of files"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /buckets/{bucketUUID}/files [get]
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

// Show retrieves details of a specific file.
//
// @Summary Show details of a single file
// @Description Get details of a specific file
// @Tags files
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param bucketUUID path string true "Bucket UUID"
// @Param fileUUID path string true "File UUID"
//
// @Success 200 {object} responses.Response{content=resources.FileResponse} "File details"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /buckets/{bucketUUID}/files/{fileUUID} [get]
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

// Store creates a new file in a bucket
//
// @Summary Create a new file
// @Description Create a new file in a specific bucket
// @Tags files
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param bucketUUID path string true "Bucket UUID"
// @Param file body bucket_requests.CreateFileRequest true "File details"
//
// @Success 201 {object} responses.Response{content=resources.FileResponse} "File details"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /buckets/{bucketUUID}/files [post]
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

// Rename updates the name of a file
//
// @Summary Rename a file
// @Description Update the name of a specific file. In other words, moves the file to a new location.
// @Tags files
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param bucketUUID path string true "Bucket UUID"
// @Param fileUUID path string true "File UUID"
// @Param file body bucket_requests.RenameFileRequest true "New file name"
//
// @Success 200 {object} responses.Response{content=resources.FileResponse} "File details"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /buckets/{bucketUUID}/files/{fileUUID}/rename [put]
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

// Delete removes a file from a bucket
//
// @Summary Delete a file
// @Description Permanently remove a specific file from a bucket
// @Tags files
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param bucketUUID path string true "Bucket UUID"
// @Param fileUUID path string true "File UUID"
//
// @Success 204 "File deleted"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /buckets/{bucketUUID}/files/{fileUUID} [delete]
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
