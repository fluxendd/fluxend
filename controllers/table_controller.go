package controllers

import (
	"fluxton/requests"
	"fluxton/requests/table_requests"
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

func (tc *TableController) List(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	projectID, err := utils.GetUUIDPathParam(c, "projectID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	paginationParams := utils.ExtractPaginationParams(c)
	tables, err := tc.tableService.List(paginationParams, projectID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.TableResourceCollection(tables))
}

func (tc *TableController) Show(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	projectID, err := utils.GetUUIDPathParam(c, "projectID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	tableID, err := utils.GetUUIDPathParam(c, "tableID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	table, err := tc.tableService.GetByID(tableID, projectID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.TableResource(&table))
}

func (tc *TableController) Store(c echo.Context) error {
	var request table_requests.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	projectID, err := utils.GetUUIDPathParam(c, "projectID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	table, err := tc.tableService.Create(&request, projectID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.TableResource(&table))
}

func (tc *TableController) Duplicate(c echo.Context) error {
	var request table_requests.RenameRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	projectID, err := utils.GetUUIDPathParam(c, "projectID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	tableID, err := utils.GetUUIDPathParam(c, "tableID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	duplicatedTable, err := tc.tableService.Duplicate(tableID, projectID, authUser, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.TableResource(duplicatedTable))
}

func (tc *TableController) Rename(c echo.Context) error {
	var request table_requests.RenameRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	projectID, err := utils.GetUUIDPathParam(c, "projectID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	tableID, err := utils.GetUUIDPathParam(c, "tableID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	renamedTable, err := tc.tableService.Rename(tableID, projectID, authUser, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.TableResource(&renamedTable))
}

func (tc *TableController) Delete(c echo.Context) error {
	var request requests.DefaultRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	projectID, err := utils.GetUUIDPathParam(c, "projectID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	tableID, err := utils.GetUUIDPathParam(c, "tableID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if _, err := tc.tableService.Delete(tableID, projectID, authUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}
