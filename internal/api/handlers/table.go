package handlers

import (
	"fluxton/internal/api/dto"
	"fluxton/internal/api/dto/database"
	databaseMapper "fluxton/internal/api/mapper/database"
	"fluxton/internal/api/response"
	databaseDomain "fluxton/internal/domain/database"
	"fluxton/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type TableHandler struct {
	tableService databaseDomain.TableService
}

func NewTableHandler(injector *do.Injector) (*TableHandler, error) {
	tableService := do.MustInvoke[databaseDomain.TableService](injector)

	return &TableHandler{tableService: tableService}, nil
}

// List retrieves all tables within a project.
//
// @Summary List all tables
// @Description Retrieve a list of tables in a specified project.
// @Tags Tables
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @param Header X-Project header string true "Project UUID"
//
// @Success 200 {object} responses.Response{content=[]resources.TableResponse} "List of tables"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /tables [get]
func (tc *TableHandler) List(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	tables, err := tc.tableService.List(request.ProjectUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, databaseMapper.ToTableResourceCollection(tables))
}

// Show retrieves details of a specific table.
//
// @Summary Get table details
// @Description Retrieve details of a specific table within a project.
// @Tags Tables
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @param Header X-Project header string true "Project UUID"
//
// @Param tableUUID path string true "Table UUID"
//
// @Success 200 {object} responses.Response{content=resources.TableResponse} "Table details"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 404 "Table not found"
// @Failure 500 "Internal server error"
//
// @Router /tables/{tableUUID} [get]
func (tc *TableHandler) Show(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	fullTableName := c.Param("fullTableName")
	if fullTableName == "" {
		return response.BadRequestResponse(c, "Table name is required")
	}

	table, err := tc.tableService.GetByName(fullTableName, request.ProjectUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, databaseMapper.ToTableResource(&table))
}

// Store creates a new table within a project.
//
// @Summary Create a new table
// @Description Define and create a new table within a specified project.
// @Tags Tables
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @param Header X-Project header string true "Project UUID"
//
// @Param table body table_requests.CreateRequest true "Table definition JSON"
//
// @Success 201 {object} responses.Response{content=resources.TableResponse} "Table created"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 422 "Unprocessable entity"
// @Failure 500 "Internal server error"
//
// @Router /tables [post]
func (tc *TableHandler) Store(c echo.Context) error {
	var request database.CreateTableRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	table, err := tc.tableService.Create(database.ToCreateTableInput(request), authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, databaseMapper.ToTableResource(&table))
}

// Upload creates a new table within a project using uploaded file
//
// @Summary Create a new table
// @Description Define and create a new table within a specified project.
// @Tags Tables
//
// @Accept Multipart/form-data
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @param Header X-Project header string true "Project UUID"
//
// @Param table body table_requests.UploadRequest true "Table definition multipart/form-data"
//
// @Success 201 {object} responses.Response{content=resources.TableResponse} "Table created"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 422 "Unprocessable entity"
// @Failure 500 "Internal server error"
//
// @Router /tables/upload [post]
func (tc *TableHandler) Upload(c echo.Context) error {
	var request database.UploadTableRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	table, err := tc.tableService.Upload(database.ToUploadTableInput(request), authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, databaseMapper.ToTableResource(&table))
}

// Duplicate creates a duplicate of an existing table.
//
// @Summary Duplicate a table
// @Description Create a copy of a specified table within a project.
// @Tags Tables
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @param Header X-Project header string true "Project UUID"
//
// @Param tableUUID path string true "Table UUID"
// @Param new_name body table_requests.RenameRequest true "Duplicate table name JSON"
//
// @Success 201 {object} responses.Response{content=resources.TableResponse} "Table duplicated"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 422 "Unprocessable entity"
// @Failure 500 "Internal server error"
//
// @Router /tables/{tableUUID}/duplicate [put]
func (tc *TableHandler) Duplicate(c echo.Context) error {
	var request database.RenameTableRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	fullTableName := c.Param("fullTableName")
	if fullTableName == "" {
		return response.BadRequestResponse(c, "Table name is required")
	}

	duplicatedTable, err := tc.tableService.Duplicate(fullTableName, authUser, database.ToRenameTableInput(request))
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, databaseMapper.ToTableResource(duplicatedTable))
}

// Rename updates the name of an existing table.
//
// @Summary Rename a table
// @Description Change the name of a specific table within a project.
// @Tags Tables
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @param Header X-Project header string true "Project UUID"
//
// @Param tableUUID path string true "Table UUID"
// @Param new_name body table_requests.RenameRequest true "New table name JSON"
//
// @Success 200 {object} responses.Response{content=resources.TableResponse} "Table renamed"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 422 "Unprocessable entity"
// @Failure 500 "Internal server error"
//
// @Router /tables/{tableUUID}/rename [put]
func (tc *TableHandler) Rename(c echo.Context) error {
	var request database.RenameTableRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	fullTableName := c.Param("fullTableName")
	if fullTableName == "" {
		return response.BadRequestResponse(c, "Table name is required")
	}

	renamedTable, err := tc.tableService.Rename(fullTableName, authUser, database.ToRenameTableInput(request))
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, databaseMapper.ToTableResource(&renamedTable))
}

// Delete removes a table permanently from a project.
//
// @Summary Delete a table
// @Description Permanently delete a specific table from a given project.
// @Tags Tables
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @param Header X-Project header string true "Project UUID"
//
// @Param tableUUID path string true "Table UUID"
//
// @Success 204 "Table deleted successfully"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 404 "Table not found"
// @Failure 500 "Internal server error"
//
// @Router /tables/{tableUUID} [delete]
func (tc *TableHandler) Delete(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	fullTableName := c.Param("fullTableName")
	if fullTableName == "" {
		return response.BadRequestResponse(c, "Table name is required")
	}

	if _, err := tc.tableService.Delete(fullTableName, request.ProjectUUID, authUser); err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.DeletedResponse(c, nil)
}
