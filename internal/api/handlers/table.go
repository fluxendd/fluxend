package handlers

import (
	"fluxend/internal/api/dto"
	databaseDto "fluxend/internal/api/dto/database"
	"fluxend/internal/api/mapper"
	"fluxend/internal/api/response"
	"fluxend/internal/domain/database"
	"fluxend/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type TableHandler struct {
	tableService database.TableService
}

func NewTableHandler(injector *do.Injector) (*TableHandler, error) {
	tableService := do.MustInvoke[database.TableService](injector)

	return &TableHandler{tableService: tableService}, nil
}

// List retrieves all tables within a project.
//
// @Summary List
// @Description Retrieve a list of tables in a specified project.
// @Tags Tables
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @param Header X-Project header string true "Project UUID"
//
// @Success 200 {object} response.Response{content=[]database.TableResponse} "List of tables"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /tables [get]
func (th *TableHandler) List(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	tables, err := th.tableService.List(request.ProjectUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToTableResourceCollection(tables))
}

// Show retrieves details of a specific table.
//
// @Summary Retrieve
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
// @Success 200 {object} response.Response{content=database.TableResponse} "Table details"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 404 "Table not found"
// @Failure 500 "Internal server error"
//
// @Router /tables/{tableUUID} [get]
func (th *TableHandler) Show(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	fullTableName := c.Param("fullTableName")
	if fullTableName == "" {
		return response.BadRequestResponse(c, "Table name is required")
	}

	table, err := th.tableService.GetByName(fullTableName, request.ProjectUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToTableResource(&table))
}

// Store creates a new table within a project.
//
// @Summary Create
// @Description Define and create a new table within a specified project.
// @Tags Tables
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @param Header X-Project header string true "Project UUID"
//
// @Param table body database.CreateTableRequest true "Table definition JSON"
//
// @Success 201 {object} response.Response{content=database.TableResponse} "Table created"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 422 "Unprocessable entity"
// @Failure 500 "Internal server error"
//
// @Router /tables [post]
func (th *TableHandler) Store(c echo.Context) error {
	var request databaseDto.CreateTableRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	table, err := th.tableService.Create(databaseDto.ToCreateTableInput(request), authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, mapper.ToTableResource(&table))
}

// Upload creates a new table within a project using uploaded file
//
// @Summary Upload
// @Description Define and create a new table within a specified project.
// @Tags Tables
//
// @Accept Multipart/form-data
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @param Header X-Project header string true "Project UUID"
//
// @Param table body database.UploadTableRequest true "Table definition multipart/form-data"
//
// @Success 201 {object} response.Response{content=database.TableResponse} "Table created"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 422 "Unprocessable entity"
// @Failure 500 "Internal server error"
//
// @Router /tables/upload [post]
func (th *TableHandler) Upload(c echo.Context) error {
	var request databaseDto.UploadTableRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	table, err := th.tableService.Upload(databaseDto.ToUploadTableInput(request), authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, mapper.ToTableResource(&table))
}

// Duplicate creates a duplicate of an existing table.
//
// @Summary Duplicate
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
// @Param new_name body database.RenameTableRequest true "Duplicate table name JSON"
//
// @Success 201 {object} response.Response{content=database.TableResponse} "Table duplicated"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 422 "Unprocessable entity"
// @Failure 500 "Internal server error"
//
// @Router /tables/{tableUUID}/duplicate [put]
func (th *TableHandler) Duplicate(c echo.Context) error {
	var request databaseDto.RenameTableRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	fullTableName := c.Param("fullTableName")
	if fullTableName == "" {
		return response.BadRequestResponse(c, "Table name is required")
	}

	duplicatedTable, err := th.tableService.Duplicate(fullTableName, authUser, databaseDto.ToRenameTableInput(request))
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToTableResource(duplicatedTable))
}

// Rename updates the name of an existing table.
//
// @Summary Rename
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
// @Param new_name body database.RenameTableRequest true "New table name JSON"
//
// @Success 200 {object} response.Response{content=database.TableResponse} "Table renamed"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 422 "Unprocessable entity"
// @Failure 500 "Internal server error"
//
// @Router /tables/{tableUUID}/rename [put]
func (th *TableHandler) Rename(c echo.Context) error {
	var request databaseDto.RenameTableRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	fullTableName := c.Param("fullTableName")
	if fullTableName == "" {
		return response.BadRequestResponse(c, "Table name is required")
	}

	renamedTable, err := th.tableService.Rename(fullTableName, authUser, databaseDto.ToRenameTableInput(request))
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToTableResource(&renamedTable))
}

// Delete removes a table permanently from a project.
//
// @Summary Delete
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
func (th *TableHandler) Delete(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	fullTableName := c.Param("fullTableName")
	if fullTableName == "" {
		return response.BadRequestResponse(c, "Table name is required")
	}

	if _, err := th.tableService.Delete(fullTableName, request.ProjectUUID, authUser); err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.DeletedResponse(c, nil)
}
