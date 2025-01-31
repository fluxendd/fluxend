package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"myapp/requests"
	"myapp/resources"
	"myapp/responses"
	"myapp/services"
	"myapp/utils"
)

type ProjectController struct {
	projectService services.ProjectService
}

func NewProjectController(injector *do.Injector) (*ProjectController, error) {
	projectService := do.MustInvoke[services.ProjectService](injector)

	return &ProjectController{projectService: projectService}, nil
}

func (nc *ProjectController) List(c echo.Context) error {
	authenticatedUserId, _ := utils.NewAuth(c).Id()

	paginationParams := utils.ExtractPaginationParams(c)
	projects, err := nc.projectService.List(paginationParams, authenticatedUserId)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.ProjectResourceCollection(projects))
}

func (nc *ProjectController) Show(c echo.Context) error {
	authenticatedUser, _ := utils.NewAuth(c).User()

	id, err := utils.GetUintPathParam(c, "id", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	project, err := nc.projectService.GetByID(id, authenticatedUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.ProjectResource(&project))
}

func (nc *ProjectController) Store(c echo.Context) error {
	var request requests.ProjectCreateRequest
	authenticatedUser, _ := utils.NewAuth(c).User()

	if err := c.Bind(&request); err != nil {
		return responses.BadRequestResponse(c, "project.error.invalidPayload")
	}

	if err := request.Validate(); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	project, err := nc.projectService.Create(&request, authenticatedUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.ProjectResource(&project))
}

func (nc *ProjectController) Update(c echo.Context) error {
	var request requests.ProjectCreateRequest
	authenticatedUser, _ := utils.NewAuth(c).User()

	id, err := utils.GetUintPathParam(c, "id", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if err := c.Bind(&request); err != nil {
		return responses.BadRequestResponse(c, "project.error.invalidPayload")
	}

	updatedOrganization, err := nc.projectService.Update(id, authenticatedUser, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.ProjectResource(updatedOrganization))
}

func (nc *ProjectController) Delete(c echo.Context) error {
	authenticatedUser, _ := utils.NewAuth(c).User()

	id, err := utils.GetUintPathParam(c, "id", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if _, err := nc.projectService.Delete(id, authenticatedUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}
