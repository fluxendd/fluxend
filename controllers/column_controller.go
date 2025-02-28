package controllers

import (
	"fluxton/requests/column_requests"
	"fluxton/resources"
	"fluxton/responses"
	"fluxton/services"
	"fluxton/utils"
	"github.com/google/uuid"
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

func (cc *ColumnController) Store(c echo.Context) error {
	var request column_requests.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	projectID, tableID, _, err := cc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	table, err := cc.columnService.CreateMany(projectID, tableID, &request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.TableResource(&table))
}

func (cc *ColumnController) Alter(c echo.Context) error {
	var request column_requests.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	projectID, tableID, _, err := cc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	alteredTable, err := cc.columnService.AlterMany(tableID, projectID, &request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.TableResource(alteredTable))
}

func (cc *ColumnController) Rename(c echo.Context) error {
	var request column_requests.RenameRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	projectID, tableID, columnName, err := cc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	renamedTable, err := cc.columnService.Rename(columnName, tableID, projectID, &request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.TableResource(renamedTable))
}

func (cc *ColumnController) Delete(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	projectID, tableID, columnName, err := cc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if _, err := cc.columnService.Delete(columnName, tableID, projectID, authUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}

func (cc *ColumnController) parseRequest(c echo.Context) (uuid.UUID, uuid.UUID, string, error) {
	projectID, err := utils.GetUUIDPathParam(c, "projectID", true)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, "", err
	}

	tableID, err := utils.GetUUIDPathParam(c, "tableID", true)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, "", err
	}

	columnName := c.Param("columnName")

	return projectID, tableID, columnName, nil
}
