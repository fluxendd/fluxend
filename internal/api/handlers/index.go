package handlers

import (
	"fluxend/internal/api/dto"
	databaseDto "fluxend/internal/api/dto/database"
	"fluxend/internal/api/response"
	databaseDomain "fluxend/internal/domain/database"
	"fluxend/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type IndexHandler struct {
	indexService databaseDomain.IndexService
}

func NewIndexHandler(injector *do.Injector) (*IndexHandler, error) {
	indexService := do.MustInvoke[databaseDomain.IndexService](injector)

	return &IndexHandler{indexService: indexService}, nil
}

// List Indexes
//
// @Summary List indexes
// @Description Retrieve a list of indexes for a given table.
// @Tags Indexes
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param tableUUID path string true "Table UUID"
//
// @Success 200 {object} response.Response{content=[]dto.GenericResponse} "List of indexes"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /tables/{tableUUID}/indexes [get]
func (ih *IndexHandler) List(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	fullTableName := c.Param("fullTableName")
	if fullTableName == "" {
		return response.BadRequestResponse(c, "Table name is required")
	}

	indexes, err := ih.indexService.List(fullTableName, request.ProjectUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, dto.GenericResourceCollection(indexes))
}

// Show Index
//
// @Summary Retrieve index
// @Description Retrieve details for a specific index in a table.
// @Tags Indexes
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param tableUUID path string true "Table UUID"
// @Param index_name path string true "Index Name"
//
// @Success 200 {object} response.Response{content=dto.GenericResponse} "Index details"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 404 "Index not found"
// @Failure 500 "Internal server error"
//
// @Router /tables/{tableUUID}/indexes/{indexName} [get]
func (ih *IndexHandler) Show(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	fullTableName := c.Param("fullTableName")
	if fullTableName == "" {
		return response.BadRequestResponse(c, "Table name is required")
	}

	indexName := c.Param("indexName")

	index, err := ih.indexService.GetByName(indexName, fullTableName, request.ProjectUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, dto.GenericResource(&index))
}

// Store Index
//
// @Summary Create index
// @Description Add an index to a specified table within a project.
// @Tags Indexes
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param tableUUID path string true "Table UUID"
// @Param index body database.CreateIndexRequest true "Index details JSON"
//
// @Success 201 {object} response.Response{content=dto.GenericResponse} "Index created"
// @Failure 422 {object} response.UnprocessableErrorResponse "Unprocessable input response"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /tables/{tableUUID}/indexes [post]
func (ih *IndexHandler) Store(c echo.Context) error {
	var request databaseDto.CreateIndexRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	fullTableName := c.Param("fullTableName")
	if fullTableName == "" {
		return response.BadRequestResponse(c, "Table name is required")
	}

	index, err := ih.indexService.Create(fullTableName, databaseDto.ToCreateIndexInput(request), authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, dto.GenericResource(index))
}

// Delete Index
//
// @Summary Delete index
// @Description Remove an existing index from a given table.
// @Tags Indexes
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param tableUUID path string true "Table UUID"
// @Param index_name path string true "Index Name"
//
// @Success 204 "Index deleted successfully"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 404 "Index not found"
// @Failure 500 "Internal server error"
//
// @Router /tables/{tableUUID}/indexes/{indexName} [delete]
func (ih *IndexHandler) Delete(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	fullTableName := c.Param("fullTableName")
	if fullTableName == "" {
		return response.BadRequestResponse(c, "Table name is required")
	}

	indexName := c.Param("indexName")

	if _, err := ih.indexService.Delete(indexName, fullTableName, request.ProjectUUID, authUser); err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.DeletedResponse(c, nil)
}
