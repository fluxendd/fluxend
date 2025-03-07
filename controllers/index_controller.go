package controllers

import (
	"fluxton/requests"
	"fluxton/resources"
	"fluxton/responses"
	"fluxton/services"
	"fluxton/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type IndexController struct {
	indexService services.IndexService
}

func NewIndexController(injector *do.Injector) (*IndexController, error) {
	indexService := do.MustInvoke[services.IndexService](injector)

	return &IndexController{indexService: indexService}, nil
}

// List Indexes
//
// @Summary List indexes for a table
// @Description Retrieve a list of indexes for a given table.
// @Tags Indexes
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param project_id path string true "Project ID"
// @Param table_id path string true "Table ID"
//
// @Success 200 {object} responses.Response{content=[]resources.GenericResponse} "List of indexes"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /tables/{tableUUID}/indexes [get]
func (ic *IndexController) List(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	projectUUID, tableUUID, err := ic.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	indexes, err := ic.indexService.List(tableUUID, projectUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.GenericResourceCollection(indexes))
}

// Show Index
//
// @Summary Show details of a specific index
// @Description Retrieve details for a specific index in a table.
// @Tags Indexes
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param project_id path string true "Project ID"
// @Param table_id path string true "Table ID"
// @Param index_name path string true "Index Name"
//
// @Success 200 {object} responses.Response{content=resources.GenericResponse} "Index details"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 404 "Index not found"
// @Failure 500 "Internal server error"
//
// @Router /tables/{tableUUID}/indexes/{index_name} [get]
func (ic *IndexController) Show(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	projectUUID, tableUUID, err := ic.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	indexName := c.Param("indexName")

	index, err := ic.indexService.GetByName(indexName, tableUUID, projectUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.GenericResource(&index))
}

// Store Index
//
// @Summary Create a new index
// @Description Add an index to a specified table within a project.
// @Tags Indexes
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param project_id path string true "Project ID"
// @Param table_id path string true "Table ID"
// @Param index body requests.IndexCreateRequest true "Index details JSON"
//
// @Success 201 {object} responses.Response{content=resources.GenericResponse} "Index created"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 422 "Unprocessable entity"
// @Failure 500 "Internal server error"
//
// @Router /tables/{tableUUID}/indexes [post]
func (ic *IndexController) Store(c echo.Context) error {
	var request requests.IndexCreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	projectUUID, tableUUID, err := ic.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	index, err := ic.indexService.Create(projectUUID, tableUUID, &request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.GenericResource(index))
}

// Delete Index
//
// @Summary Delete an index from a table
// @Description Remove an existing index from a given table.
// @Tags Indexes
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param project_id path string true "Project ID"
// @Param table_id path string true "Table ID"
// @Param index_name path string true "Index Name"
//
// @Success 204 "Index deleted successfully"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 404 "Index not found"
// @Failure 500 "Internal server error"
//
// @Router /tables/{tableUUID}/indexes/{index_name} [delete]
func (ic *IndexController) Delete(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	projectUUID, tableUUID, err := ic.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	indexName := c.Param("indexName")

	if _, err := ic.indexService.Delete(indexName, tableUUID, projectUUID, authUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}

func (ic *IndexController) parseRequest(c echo.Context) (uuid.UUID, uuid.UUID, error) {
	projectUUID, err := utils.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	tableUUID, err := utils.GetUUIDPathParam(c, "tableUUID", true)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return projectUUID, tableUUID, nil
}
