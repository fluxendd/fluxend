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
	"strings"
)

type ProjectController struct {
	projectService services.ProjectService
}

func NewProjectController(injector *do.Injector) (*ProjectController, error) {
	projectService := do.MustInvoke[services.ProjectService](injector)

	return &ProjectController{projectService: projectService}, nil
}

func (pc *ProjectController) List(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	organizationUUID, err := uuid.Parse(strings.TrimSpace(c.QueryParam("organization_id")))
	if err != nil {
		return responses.BadRequestResponse(c, "Invalid organization ID")
	}

	paginationParams := utils.ExtractPaginationParams(c)
	projects, err := pc.projectService.List(paginationParams, organizationUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.ProjectResourceCollection(projects))
}

func (pc *ProjectController) Show(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	projectID, err := utils.GetUUIDPathParam(c, "projectID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	project, err := pc.projectService.GetByID(projectID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.ProjectResource(&project))
}

func (pc *ProjectController) Store(c echo.Context) error {
	var request requests.ProjectCreateRequest
	authUser, _ := utils.NewAuth(c).User()

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	project, err := pc.projectService.Create(&request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.ProjectResource(&project))
}

func (pc *ProjectController) Update(c echo.Context) error {
	var request requests.ProjectUpdateRequest
	authUser, _ := utils.NewAuth(c).User()

	projectID, err := utils.GetUUIDPathParam(c, "projectID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	updatedOrganization, err := pc.projectService.Update(projectID, authUser, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.ProjectResource(updatedOrganization))
}

func (pc *ProjectController) Delete(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	projectID, err := utils.GetUUIDPathParam(c, "projectID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if _, err := pc.projectService.Delete(projectID, authUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}
