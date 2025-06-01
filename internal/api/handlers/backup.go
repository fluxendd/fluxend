package handlers

import (
	"fluxend/internal/api/dto"
	"fluxend/internal/api/mapper"
	"fluxend/internal/api/response"
	"fluxend/internal/domain/backup"
	"fluxend/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type BackupHandler struct {
	backupService backup.Service
}

func NewBackupHandler(injector *do.Injector) (*BackupHandler, error) {
	backupService := do.MustInvoke[backup.Service](injector)

	return &BackupHandler{backupService: backupService}, nil
}

// List retrieves all backups for a project
//
// @Summary List backups
// @Description Retrieve a list of all backups for the specified project
// @Tags Backups
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Success 200 {array} response.Response{content=[]backup.Response} "List of backups"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /backups [get]
func (bh *BackupHandler) List(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}
	authUser, _ := auth.NewAuth(c).User()

	backups, err := bh.backupService.List(request.ProjectUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToBackupResourceCollection(backups))
}

// Show retrieves details of a specific backup
//
// @Summary Retrieve backup
// @Description Get details of a specific backup
// @Tags Backups
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param backupUUID path string true "Backup UUID"
//
// @Success 200 {object} response.Response{content=backup.Response} "Backup details"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /backups/{backupUUID} [get]
func (bh *BackupHandler) Show(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	backupUUID, err := request.GetUUIDPathParam(c, "backupUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	fetchedBackup, err := bh.backupService.GetByUUID(backupUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToBackupResource(&fetchedBackup))
}

// Store creates a new backup
//
// @Summary Create backup
// @Description Create a new backup
// @Tags Backups
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param backup body dto.DefaultRequestWithProjectHeader true "Project UUID"
//
// @Success 201 {object} response.Response{content=backup.Response} "Backup created"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /backups [post]
func (bh *BackupHandler) Store(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	storedBackup, err := bh.backupService.Create(request.ProjectUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, mapper.ToBackupResource(&storedBackup))
}

// Delete removes a backup
//
// @Summary Delete backup
// @Description Remove a backup from the project
// @Tags Backups
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param backupUUID path string true "Backup UUID"
//
// @Success 204 "Backup deleted"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /backups/{backupUUID} [delete]
func (bh *BackupHandler) Delete(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	backupUUID, err := request.GetUUIDPathParam(c, "backupUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	if _, err := bh.backupService.Delete(backupUUID, authUser); err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.DeletedResponse(c, nil)
}
