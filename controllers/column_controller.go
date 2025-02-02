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

type ColumnController struct {
	columnService services.ColumnService
}

func NewColumnController(injector *do.Injector) (*ColumnController, error) {
	columnService := do.MustInvoke[services.ColumnService](injector)

	return &ColumnController{columnService: columnService}, nil
}

func (pc *ColumnController) Store(c echo.Context) error {
	var request requests.ColumnCreateRequest
	authenticatedUser, _ := utils.NewAuth(c).User()

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	projectID, err := utils.GetUintPathParam(c, "projectID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	tableID, err := utils.GetUintPathParam(c, "tableID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	table, err := pc.columnService.Create(projectID, tableID, &request, authenticatedUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.TableResource(&table))
}

func (pc *ColumnController) Alter(c echo.Context) error {
	var request requests.ColumnUpdateRequest
	authenticatedUser, _ := utils.NewAuth(c).User()

	projectID, err := utils.GetUintPathParam(c, "projectID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	tableID, err := utils.GetUintPathParam(c, "tableID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	renamedTable, err := pc.columnService.Alter(tableID, projectID, &request, authenticatedUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.TableResource(renamedTable))
}

func (pc *ColumnController) Delete(c echo.Context) error {
	var request requests.ColumnDeleteRequest
	authenticatedUser, _ := utils.NewAuth(c).User()

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	projectID, err := utils.GetUintPathParam(c, "projectID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	tableID, err := utils.GetUintPathParam(c, "tableID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if _, err := pc.columnService.Delete(tableID, request.OrganizationID, projectID, &request, authenticatedUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}
