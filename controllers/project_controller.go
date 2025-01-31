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

func (pc *ProjectController) List(c echo.Context) error {
	var request requests.ProjectDefaultRequest
	authenticatedUserId, _ := utils.NewAuth(c).Id()

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	paginationParams := utils.ExtractPaginationParams(c)
	projects, err := pc.projectService.List(paginationParams, request.OrganizationID, authenticatedUserId)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.ProjectResourceCollection(projects))
}

func (pc *ProjectController) Show(c echo.Context) error {
	var request requests.ProjectDefaultRequest
	authenticatedUser, _ := utils.NewAuth(c).User()

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	id, err := utils.GetUintPathParam(c, "id", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	project, err := pc.projectService.GetByID(id, request.OrganizationID, authenticatedUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.ProjectResource(&project))
}

func (pc *ProjectController) Store(c echo.Context) error {
	var request requests.ProjectCreateRequest
	authenticatedUser, _ := utils.NewAuth(c).User()

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	project, err := pc.projectService.Create(&request, authenticatedUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.ProjectResource(&project))
}

func (pc *ProjectController) Update(c echo.Context) error {
	var request requests.ProjectCreateRequest
	authenticatedUser, _ := utils.NewAuth(c).User()

	id, err := utils.GetUintPathParam(c, "id", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	updatedOrganization, err := pc.projectService.Update(id, authenticatedUser, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.ProjectResource(updatedOrganization))
}

func (pc *ProjectController) Delete(c echo.Context) error {
	var request requests.ProjectDefaultRequest
	authenticatedUser, _ := utils.NewAuth(c).User()

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	id, err := utils.GetUintPathParam(c, "id", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if _, err := pc.projectService.Delete(id, authenticatedUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}
