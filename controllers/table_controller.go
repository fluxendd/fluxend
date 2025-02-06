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

type TableController struct {
	tableService services.TableService
}

func NewTableController(injector *do.Injector) (*TableController, error) {
	tableService := do.MustInvoke[services.TableService](injector)

	return &TableController{tableService: tableService}, nil
}

func (pc *TableController) List(c echo.Context) error {
	var request requests.DefaultRequest
	authUserId, _ := utils.NewAuth(c).Id()

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	projectID, err := utils.GetUUIDPathParam(c, "projectID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	paginationParams := utils.ExtractPaginationParams(c)
	tables, err := pc.tableService.List(paginationParams, request.OrganizationID, projectID, authUserId)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.TableResourceCollection(tables))
}

func (pc *TableController) Show(c echo.Context) error {
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

func (pc *TableController) Store(c echo.Context) error {
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

func (pc *TableController) Rename(c echo.Context) error {
	var request requests.TableRenameRequest
	authUser, _ := utils.NewAuth(c).User()

	projectID, err := utils.GetUUIDPathParam(c, "projectID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	tableID, err := utils.GetUUIDPathParam(c, "tableID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	renamedTable, err := pc.tableService.Rename(tableID, projectID, authUser, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.TableResource(&renamedTable))
}

func (pc *TableController) Delete(c echo.Context) error {
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
