package controllers

import (
	"fluxton/errs"
	"fluxton/requests"
	"fluxton/requests/column_requests"
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

// List retrieves all columns within a project.
//
// @Summary List all columns
// @Description Retrieve a list of columns in a specified table.
// @Tags Columns
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @param Header X-Project header string true "Project UUID"
//
// @Success 200 {object} responses.Response{content=[]resources.ColumnResponse} "List of columns"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /tables/{fullTableName}/columns [get]
func (cc *ColumnController) List(c echo.Context) error {
	var request requests.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	fullTableName := c.Param("fullTableName")
	if fullTableName == "" {
		return responses.BadRequestResponse(c, "Table name is required")
	}

	columns, err := cc.columnService.List(fullTableName, request.ProjectUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.ColumnResourceCollection(columns))
}

// Store adds new columns to a table.
//
// @Summary Add new columns to a table
// @Description Create new columns in a specified table within a project.
// @Tags Columns
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param fullTableName path string true "Full table name"
// @Param columns body column_requests.CreateRequest true "Columns JSON"
//
// @Success 201 {object} responses.Response{content=resources.TableResponse} "Columns created"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 422 "Unprocessable entity"
// @Failure 500 "Internal server error"
//
// @Router /tables/{fullTableName}/columns [post]
func (cc *ColumnController) Store(c echo.Context) error {
	var request column_requests.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	fullTableName := c.Param("fullTableName")
	if fullTableName == "" {
		return responses.BadRequestResponse(c, "Table name is required")
	}

	columns, err := cc.columnService.CreateMany(fullTableName, &request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.ColumnResourceCollection(columns))
}

// Alter modifies column types in a table.
//
// @Summary Modify column types in a table
// @Description Alter the data type of existing columns in a specified table.
// @Tags Columns
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param fullTableName path string true "Full table name"
// @Param columns body column_requests.CreateRequest true "Updated column definitions"
//
// @Success 200 {object} responses.Response{content=resources.TableResponse} "Columns altered"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 422 "Unprocessable entity"
// @Failure 500 "Internal server error"
//
// @Router /tables/{fullTableName}/columns [put]
func (cc *ColumnController) Alter(c echo.Context) error {
	var request column_requests.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	fullTableName := c.Param("fullTableName")
	if fullTableName == "" {
		return responses.BadRequestResponse(c, "Table name is required")
	}

	columns, err := cc.columnService.AlterMany(fullTableName, &request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.ColumnResourceCollection(columns))
}

// Rename updates the name of an existing column.
//
// @Summary Rename a column in a table
// @Description Change the name of a specific column in a given table.
// @Tags Columns
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param fullTableName path string true "Full table name"
// @Param column_name path string true "Existing Column Name"
// @Param new_name body column_requests.RenameRequest true "New column name JSON"
//
// @Success 200 {object} responses.Response "Column renamed"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 422 "Unprocessable entity"
// @Failure 500 "Internal server error"
//
// @Router /tables/{fullTableName}/columns/{columnName} [put]
func (cc *ColumnController) Rename(c echo.Context) error {
	var request column_requests.RenameRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	fullTableName, columnName, err := cc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	columns, err := cc.columnService.Rename(columnName, fullTableName, &request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.ColumnResourceCollection(columns))
}

// Delete removes a column from a table.
//
// @Summary Delete a column from a table
// @Description Permanently delete a specific column from a given table.
// @Tags Columns
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param fullTableName path string true "Full table name"
// @Param column_name path string true "Column Name"
//
// @Success 204 "Column deleted successfully"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 404 "Column not found"
// @Failure 500 "Internal server error"
//
// @Router /tables/{fullTableName}/columns/{columnName} [delete]
func (cc *ColumnController) Delete(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	fullTableName, columnName, err := cc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	projectUUID, err := utils.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if _, err := cc.columnService.Delete(columnName, fullTableName, projectUUID, authUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}

func (cc *ColumnController) parseRequest(c echo.Context) (string, string, error) {
	fullTableName := c.Param("fullTableName")
	if fullTableName == "" {
		return "", "", errs.NewBadRequestError("Table name is required")
	}

	columnName := c.Param("columnName")

	return fullTableName, columnName, nil
}
