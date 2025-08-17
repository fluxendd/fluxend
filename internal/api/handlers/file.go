package handlers

import (
	"fluxend/internal/api/dto"
	fileDto "fluxend/internal/api/dto/storage/file"
	"fluxend/internal/api/mapper"
	"fluxend/internal/api/response"
	"fluxend/internal/domain/storage/file"
	"fluxend/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type FileHandler struct {
	fileService file.Service
}

func NewFileHandler(injector *do.Injector) (*FileHandler, error) {
	fileService := do.MustInvoke[file.Service](injector)

	return &FileHandler{fileService: fileService}, nil
}

// List retrieves all files in a container
//
// @Summary List files
// @Description Retrieve a list of all files in a specific container
// @Tags Files
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param containerUUID path string true "Container UUID"
//
// @Param page query string false "Page number for pagination"
// @Param limit query string false "Number of items per page"
// @Param sort query string false "Field to sort by"
// @Param order query string false "Sort order (asc or desc)"
//
// @Success 200 {array} response.Response{content=[]file.Response} "List of files"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /containers/{containerUUID}/files [get]
func (fh *FileHandler) List(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	containerUUID, err := request.GetUUIDPathParam(c, "containerUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, "Invalid container UUID")
	}

	paginationParams := request.ExtractPaginationParams(c)
	files, err := fh.fileService.List(paginationParams, containerUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToFileResourceCollection(files))
}

// Show retrieves details of a specific file.
//
// @Summary Retrieve file
// @Description Get details of a specific file
// @Tags Files
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param containerUUID path string true "Container UUID"
// @Param fileUUID path string true "File UUID"
//
// @Success 200 {object} response.Response{content=file.Response} "File details"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /containers/{containerUUID}/files/{fileUUID} [get]
func (fh *FileHandler) Show(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
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

	fetchedFile, err := fh.fileService.GetByUUID(fileUUID, containerUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToFileResource(&fetchedFile))
}

// Store creates a new file in a container
//
// @Summary Create file
// @Description Create a new file in a specific container
// @Tags Files
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param containerUUID path string true "Container UUID"
// @Param file body file.CreateRequest true "File details"
//
// @Success 201 {object} response.Response{content=file.Response} "File details"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /containers/{containerUUID}/files [post]
func (fh *FileHandler) Store(c echo.Context) error {
	var request fileDto.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	containerUUID, err := request.GetUUIDPathParam(c, "containerUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	createdFile, err := fh.fileService.Create(containerUUID, fileDto.ToCreateFileInput(&request), authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, mapper.ToFileResource(&createdFile))
}

// Rename updates the name of a file
//
// @Summary Rename file
// @Description Update the name of a specific file. In other words, moves the file to a new location.
// @Tags Files
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param containerUUID path string true "Container UUID"
// @Param fileUUID path string true "File UUID"
// @Param file body file.RenameRequest true "New file name"
//
// @Success 200 {object} response.Response{content=file.Response} "File details"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /containers/{containerUUID}/files/{fileUUID} [put]
func (fh *FileHandler) Rename(c echo.Context) error {
	var request fileDto.RenameRequest
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

	updatedFile, err := fh.fileService.Rename(fileUUID, containerUUID, authUser, fileDto.ToRenameFileInput(&request))
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToFileResource(updatedFile))
}

// Download Retrieves a presigned URL for downloading a file
//
// @Summary Download file
// @Description Get a presigned URL to download a specific file
// @Tags Files
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param containerUUID path string true "Container UUID"
// @Param fileUUID path string true "File UUID"
//
// @Success 200 {object} response.Response{content=file.DownloadResponse} "File details"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /containers/{containerUUID}/files/{fileUUID}/download [get]
func (fh *FileHandler) Download(c echo.Context) error {
	var request dto.DefaultRequest
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

	url, err := fh.fileService.CreatePresignedURL(fileUUID, containerUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToDownloadResource(url, 3600))
}

// Delete removes a file from a container
//
// @Summary Delete file
// @Description Permanently remove a specific file from a container
// @Tags Files
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param containerUUID path string true "Container UUID"
// @Param fileUUID path string true "File UUID"
//
// @Success 204 "File deleted"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /containers/{containerUUID}/files/{fileUUID} [delete]
func (fh *FileHandler) Delete(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
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

	if _, err := fh.fileService.Delete(fileUUID, containerUUID, authUser); err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.DeletedResponse(c, nil)
}
