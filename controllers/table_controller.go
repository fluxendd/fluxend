package controllers

import (
	"fluxton/internal/api/response"
	"fluxton/pkg/auth"
	"fluxton/requests"
	"fluxton/requests/table_requests"
	"fluxton/resources"
	"fluxton/services"
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
func (tc *TableController) List(c echo.Context) error {
	var request requests.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	tables, err := tc.tableService.List(request.ProjectUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, resources.TableResourceCollection(tables))
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
func (tc *TableController) Show(c echo.Context) error {
	var request requests.DefaultRequestWithProjectHeader
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

	return response.SuccessResponse(c, resources.TableResource(&table))
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
func (tc *TableController) Store(c echo.Context) error {
	var request table_requests.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	table, err := tc.tableService.Create(&request, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, resources.TableResource(&table))
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
func (tc *TableController) Upload(c echo.Context) error {
	var request table_requests.UploadRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	table, err := tc.tableService.Upload(&request, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, resources.TableResource(&table))
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
func (tc *TableController) Duplicate(c echo.Context) error {
	var request table_requests.RenameRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	fullTableName := c.Param("fullTableName")
	if fullTableName == "" {
		return response.BadRequestResponse(c, "Table name is required")
	}

	duplicatedTable, err := tc.tableService.Duplicate(fullTableName, authUser, &request)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, resources.TableResource(duplicatedTable))
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
func (tc *TableController) Rename(c echo.Context) error {
	var request table_requests.RenameRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	fullTableName := c.Param("fullTableName")
	if fullTableName == "" {
		return response.BadRequestResponse(c, "Table name is required")
	}

	renamedTable, err := tc.tableService.Rename(fullTableName, authUser, &request)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, resources.TableResource(&renamedTable))
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
func (tc *TableController) Delete(c echo.Context) error {
	var request requests.DefaultRequestWithProjectHeader
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
