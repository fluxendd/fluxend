package controllers

import (
	"fluxton/requests"
	"fluxton/resources"
	"fluxton/responses"
	"fluxton/services"
	"fluxton/utils"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type RowController struct {
	rowService services.RowService
}

func NewRowController(injector *do.Injector) (*RowController, error) {
	rowService := do.MustInvoke[services.RowService](injector)

	return &RowController{rowService: rowService}, nil
}

func (pc *RowController) Show(c echo.Context) error {
	var request requests.DefaultRequest
	authenticatedUser, _ := utils.NewAuth(c).User()

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	projectID, err := utils.GetUUIDPathParam(c, "projectID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	rowID, err := utils.GetUintPathParam(c, "rowID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	row, err := pc.rowService.GetByID(c.Param("tableName"), uint64(rowID), request.OrganizationID, projectID, authenticatedUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.RowResource(row))
}

func (pc *RowController) List(c echo.Context) error {
	var request requests.DefaultRequest
	authenticatedUserId, _ := utils.NewAuth(c).Id()

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	projectID, err := utils.GetUUIDPathParam(c, "projectID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	paginationParams := utils.ExtractPaginationParams(c)
	rows, err := pc.rowService.List(paginationParams, c.Param("tableName"), request.OrganizationID, projectID, authenticatedUserId)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.RowResourceCollection(rows))
}

func (rc *RowController) Store(c echo.Context) error {
	var request requests.RowCreateRequest
	authenticatedUser, _ := utils.NewAuth(c).User()

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	projectID, err := utils.GetUUIDPathParam(c, "projectID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	row, err := rc.rowService.Create(&request, projectID, c.Param("tableName"), authenticatedUser)
	if err != nil {
		return responses.UnprocessableResponse(c, []string{err.Error()})
	}

	return responses.CreatedResponse(c, resources.RowResource(row))
}
