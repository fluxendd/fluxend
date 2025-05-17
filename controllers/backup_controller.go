package controllers

import (
	"fluxton/internal/api/dto"
	"fluxton/internal/api/response"
	"fluxton/pkg/auth"
	"fluxton/resources"
	"fluxton/services"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type BackupController struct {
	backupService services.BackupService
}

func NewBackupController(injector *do.Injector) (*BackupController, error) {
	backupService := do.MustInvoke[services.BackupService](injector)

	return &BackupController{backupService: backupService}, nil
}

// List retrieves all backups for a project
//
// @Summary List all backups
// @Description Retrieve a list of all backups for the specified project
// @Tags Backups
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Success 200 {array} responses.Response{content=[]resources.BackupResponse} "List of backups"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /backups [get]
func (bc *BackupController) List(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}
	authUser, _ := auth.NewAuth(c).User()

	backups, err := bc.backupService.List(request.ProjectUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, resources.BackupResourceCollection(backups))
}

// Show retrieves details of a specific backup
//
// @Summary Show details of a single backup
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
// @Success 200 {object} responses.Response{content=resources.BackupResponse} "Backup details"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /backups/{backupUUID} [get]
func (bc *BackupController) Show(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	backupUUID, err := request.GetUUIDPathParam(c, "backupUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	backup, err := bc.backupService.GetByUUID(backupUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, resources.BackupResource(&backup))
}

// Store creates a new backup
//
// @Summary Create a new backup
// @Description Add a new backup
// @Tags Backups
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param backup body requests.DefaultRequestWithProjectHeader true "Project UUID"
//
// @Success 201 {object} responses.Response{content=resources.BackupResponse} "Backup created"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /backups [post]
func (bc *BackupController) Store(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	backup, err := bc.backupService.Create(&request, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, resources.BackupResource(&backup))
}

// Delete removes a backup
//
// @Summary Delete a backup
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
func (bc *BackupController) Delete(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	backupUUID, err := request.GetUUIDPathParam(c, "backupUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	if _, err := bc.backupService.Delete(request, backupUUID, authUser); err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.DeletedResponse(c, nil)
}
