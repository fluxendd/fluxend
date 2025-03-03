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

// Store adds new columns to a table.
//
// @Summary Add new columns to a table
// @Description Create new columns in a specified table within a project.
// @Tags columns
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param project_id path string true "Project ID"
// @Param table_id path string true "Table ID"
// @Param columns body column_requests.CreateRequest true "Columns JSON"
//
// @Success 201 {object} responses.Response{content=resources.TableResponse} "Columns created"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 422 "Unprocessable entity"
// @Failure 500 "Internal server error"
//
// @Router /projects/{project_id}/tables/{table_id}/columns [post]
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

// Alter modifies column types in a table.
//
// @Summary Modify column types in a table
// @Description Alter the data type of existing columns in a specified table.
// @Tags columns
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param project_id path string true "Project ID"
// @Param table_id path string true "Table ID"
// @Param columns body column_requests.CreateRequest true "Updated column definitions"
//
// @Success 200 {object} responses.Response{content=resources.TableResponse} "Columns altered"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 422 "Unprocessable entity"
// @Failure 500 "Internal server error"
//
// @Router /projects/{project_id}/tables/{table_id}/columns [put]
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

// Rename updates the name of an existing column.
//
// @Summary Rename a column in a table
// @Description Change the name of a specific column in a given table.
// @Tags columns
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param project_id path string true "Project ID"
// @Param table_id path string true "Table ID"
// @Param column_name path string true "Existing Column Name"
// @Param new_name body column_requests.RenameRequest true "New column name JSON"
//
// @Success 200 {object} responses.Response "Column renamed"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 422 "Unprocessable entity"
// @Failure 500 "Internal server error"
//
// @Router /projects/{project_id}/tables/{table_id}/columns/{column_name} [put]
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

// Delete removes a column from a table.
//
// @Summary Delete a column from a table
// @Description Permanently delete a specific column from a given table.
// @Tags columns
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param project_id path string true "Project ID"
// @Param table_id path string true "Table ID"
// @Param column_name path string true "Column Name"
//
// @Success 204 "Column deleted successfully"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 404 "Column not found"
// @Failure 500 "Internal server error"
//
// @Router /projects/{project_id}/tables/{table_id}/columns/{column_name} [delete]
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
