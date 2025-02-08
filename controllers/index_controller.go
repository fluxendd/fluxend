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

func (pc *IndexController) List(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	projectID, tableID, err := pc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	indexes, err := pc.indexService.List(tableID, projectID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.GenericResourceCollection(indexes))
}

func (pc *IndexController) Show(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	projectID, tableID, err := pc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	indexName := c.Param("indexName")

	index, err := pc.indexService.GetByName(indexName, tableID, projectID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.GenericResource(&index))
}

func (pc *IndexController) Store(c echo.Context) error {
	var request requests.IndexCreateRequest
	authUser, _ := utils.NewAuth(c).User()

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	projectID, tableID, err := pc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	index, err := pc.indexService.Create(projectID, tableID, &request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.GenericResource(index))
}

func (pc *IndexController) Delete(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	projectID, tableID, err := pc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	indexName := c.Param("indexName")

	if _, err := pc.indexService.Delete(indexName, tableID, projectID, authUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}

func (pc *IndexController) parseRequest(c echo.Context) (uuid.UUID, uuid.UUID, error) {
	projectID, err := utils.GetUUIDPathParam(c, "projectID", true)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	tableID, err := utils.GetUUIDPathParam(c, "tableID", true)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return projectID, tableID, nil
}
