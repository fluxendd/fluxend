package handlers

import (
	"fluxton/internal/api/dto"
	"fluxton/internal/api/dto/database"
	databaseMapper "fluxton/internal/api/mapper/database"
	"fluxton/internal/api/response"
	databaseDomain "fluxton/internal/domain/database"
	"fluxton/pkg/auth"
	"fluxton/pkg/errors"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type ColumnHandler struct {
	columnService databaseDomain.ColumnService
}

func NewColumnHandler(injector *do.Injector) (*ColumnHandler, error) {
	columnService := do.MustInvoke[databaseDomain.ColumnService](injector)

	return &ColumnHandler{columnService: columnService}, nil
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
// @Success 200 {object} response.Response{content=[]database.ColumnResponse} "List of columns"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /tables/{fullTableName}/columns [get]
func (ch *ColumnHandler) List(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	fullTableName := c.Param("fullTableName")
	if fullTableName == "" {
		return response.BadRequestResponse(c, "Table name is required")
	}

	columns, err := ch.columnService.List(fullTableName, request.ProjectUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, databaseMapper.ToColumnResourceCollection(columns))
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
// @Param columns body database.CreateColumnRequest true "Columns JSON"
//
// @Success 201 {object} response.Response{content=database.ColumnResponse} "Columns created"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 422 "Unprocessable entity"
// @Failure 500 "Internal server error"
//
// @Router /tables/{fullTableName}/columns [post]
func (ch *ColumnHandler) Store(c echo.Context) error {
	var request database.CreateColumnRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	fullTableName := c.Param("fullTableName")
	if fullTableName == "" {
		return response.BadRequestResponse(c, "Table name is required")
	}

	columns, err := ch.columnService.CreateMany(fullTableName, database.ToCreateColumnInput(request), authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, databaseMapper.ToColumnResourceCollection(columns))
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
// @Param columns body database.CreateColumnRequest true "Updated column definitions"
//
// @Success 200 {object} response.Response{content=database.ColumnResponse} "Columns altered"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 422 "Unprocessable entity"
// @Failure 500 "Internal server error"
//
// @Router /tables/{fullTableName}/columns [put]
func (ch *ColumnHandler) Alter(c echo.Context) error {
	var request database.CreateColumnRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	fullTableName := c.Param("fullTableName")
	if fullTableName == "" {
		return response.BadRequestResponse(c, "Table name is required")
	}

	columns, err := ch.columnService.AlterMany(fullTableName, database.ToCreateColumnInput(request), authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, databaseMapper.ToColumnResourceCollection(columns))
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
// @Param new_name body database.RenameColumnRequest true "New column name JSON"
//
// @Success 200 {object} response.Response "Column renamed"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 422 "Unprocessable entity"
// @Failure 500 "Internal server error"
//
// @Router /tables/{fullTableName}/columns/{columnName} [put]
func (ch *ColumnHandler) Rename(c echo.Context) error {
	var request database.RenameColumnRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	fullTableName, columnName, err := ch.parseRequest(c)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	columns, err := ch.columnService.Rename(columnName, fullTableName, database.ToRenameColumnInput(request), authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, databaseMapper.ToColumnResourceCollection(columns))
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
func (ch *ColumnHandler) Delete(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	fullTableName, columnName, err := ch.parseRequest(c)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	if _, err := ch.columnService.Delete(columnName, fullTableName, request.ProjectUUID, authUser); err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.DeletedResponse(c, nil)
}

func (ch *ColumnHandler) parseRequest(c echo.Context) (string, string, error) {
	fullTableName := c.Param("fullTableName")
	if fullTableName == "" {
		return "", "", errors.NewBadRequestError("Table name is required")
	}

	columnName := c.Param("columnName")

	return fullTableName, columnName, nil
}
