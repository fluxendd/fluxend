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

type IndexController struct {
	tableService services.TableService
}

func NewIndexController(injector *do.Injector) (*IndexController, error) {
	tableService := do.MustInvoke[services.TableService](injector)

	return &IndexController{tableService: tableService}, nil
}

func (pc *IndexController) List(c echo.Context) error {
	var request requests.DefaultRequest
	authUser, _ := utils.NewAuth(c).User()

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	projectID, err := utils.GetUUIDPathParam(c, "projectID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	paginationParams := utils.ExtractPaginationParams(c)
	tables, err := pc.tableService.List(paginationParams, request.OrganizationID, projectID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.TableResourceCollection(tables))
}

func (pc *IndexController) Show(c echo.Context) error {
	var request requests.DefaultRequest
	authUser, _ := utils.NewAuth(c).User()

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	tableID, err := utils.GetUUIDPathParam(c, "tableID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	table, err := pc.tableService.GetByID(tableID, request.OrganizationID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.TableResource(&table))
}

func (pc *IndexController) Store(c echo.Context) error {
	var request requests.TableCreateRequest
	authUser, _ := utils.NewAuth(c).User()

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	projectID, err := utils.GetUUIDPathParam(c, "projectID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	table, err := pc.tableService.Create(&request, projectID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.TableResource(&table))
}

func (pc *IndexController) Delete(c echo.Context) error {
	var request requests.DefaultRequest
	authUser, _ := utils.NewAuth(c).User()

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	projectID, err := utils.GetUUIDPathParam(c, "projectID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	tableID, err := utils.GetUUIDPathParam(c, "tableID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if _, err := pc.tableService.Delete(tableID, request.OrganizationID, projectID, authUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}
