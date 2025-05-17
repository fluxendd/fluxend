package controllers

import (
	"fluxton/internal/api/response"
	"fluxton/pkg/auth"
	"fluxton/requests"
	"fluxton/requests/container_requests"
	"fluxton/resources"
	"fluxton/services"
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

// List retrieves all files in a container
//
// @Summary List all files in a container
// @Description Retrieve a list of all files in a specific container
// @Tags Files
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param containerUUID path string true "Container UUID"
//
// @Param page query string false "Page number for pagination"
// @Param limit query string false "Number of items per page"
// @Param sort query string false "Field to sort by"
// @Param order query string false "Sort order (asc or desc)"
//
// @Success 200 {array} responses.Response{content=[]resources.FileResponse} "List of files"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /containers/{containerUUID}/files [get]
func (fc *FileController) List(c echo.Context) error {
	var request requests.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	containerUUID, err := request.GetUUIDPathParam(c, "containerUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, "Invalid container UUID")
	}

	paginationParams := request.ExtractPaginationParams(c)
	files, err := fc.fileService.List(paginationParams, containerUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, resources.FileResourceCollection(files))
}

// Show retrieves details of a specific file.
//
// @Summary Show details of a single file
// @Description Get details of a specific file
// @Tags Files
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param containerUUID path string true "Container UUID"
// @Param fileUUID path string true "File UUID"
//
// @Success 200 {object} responses.Response{content=resources.FileResponse} "File details"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /containers/{containerUUID}/files/{fileUUID} [get]
func (fc *FileController) Show(c echo.Context) error {
	var request requests.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	containerUUID, err := request.GetUUIDPathParam(c, "containerUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	fileUUID, err := request.GetUUIDPathParam(c, "fileUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	file, err := fc.fileService.GetByUUID(fileUUID, containerUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, resources.FileResource(&file))
}

// Store creates a new file in a container
//
// @Summary Create a new file
// @Description Create a new file in a specific container
// @Tags Files
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param containerUUID path string true "Container UUID"
// @Param file body container_requests.CreateFileRequest true "File details"
//
// @Success 201 {object} responses.Response{content=resources.FileResponse} "File details"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /containers/{containerUUID}/files [post]
func (fc *FileController) Store(c echo.Context) error {
	var request container_requests.CreateFileRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	containerUUID, err := request.GetUUIDPathParam(c, "containerUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	file, err := fc.fileService.Create(containerUUID, &request, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, resources.FileResource(&file))
}

// Rename updates the name of a file
//
// @Summary Rename a file
// @Description Update the name of a specific file. In other words, moves the file to a new location.
// @Tags Files
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param containerUUID path string true "Container UUID"
// @Param fileUUID path string true "File UUID"
// @Param file body container_requests.RenameFileRequest true "New file name"
//
// @Success 200 {object} responses.Response{content=resources.FileResponse} "File details"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /containers/{containerUUID}/files/{fileUUID}/rename [put]
func (fc *FileController) Rename(c echo.Context) error {
	var request container_requests.RenameFileRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	fileUUID, err := request.GetUUIDPathParam(c, "fileUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	containerUUID, err := request.GetUUIDPathParam(c, "containerUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	updatedFile, err := fc.fileService.Rename(fileUUID, containerUUID, authUser, &request)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, resources.FileResource(updatedFile))
}

// Delete removes a file from a container
//
// @Summary Delete a file
// @Description Permanently remove a specific file from a container
// @Tags Files
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param containerUUID path string true "Container UUID"
// @Param fileUUID path string true "File UUID"
//
// @Success 204 "File deleted"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /containers/{containerUUID}/files/{fileUUID} [delete]
func (fc *FileController) Delete(c echo.Context) error {
	var request requests.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	containerUUID, err := request.GetUUIDPathParam(c, "containerUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	fileUUID, err := request.GetUUIDPathParam(c, "fileUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	if _, err := fc.fileService.Delete(fileUUID, containerUUID, authUser, request); err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.DeletedResponse(c, nil)
}
