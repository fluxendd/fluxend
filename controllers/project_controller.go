package controllers

import (
	"fluxton/internal/api/response"
	"fluxton/pkg/auth"
	"fluxton/requests"
	"fluxton/requests/project_requests"
	"fluxton/resources"
	"fluxton/services"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type ProjectController struct {
	projectService services.ProjectService
}

func NewProjectController(injector *do.Injector) (*ProjectController, error) {
	projectService := do.MustInvoke[services.ProjectService](injector)

	return &ProjectController{projectService: projectService}, nil
}

// List all projects
//
// @Summary List all projects
// @Description Get all projects for a specific organization
// @Tags Projects
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param organizationUUID query string true "Organization UUID"
//
// @Param page query string false "Page number for pagination"
// @Param limit query string false "Number of items per page"
// @Param sort query string false "Field to sort by"
// @Param order query string false "Sort order (asc or desc)"
//
// @Success 200 {object} responses.Response{content=[]resources.ProjectResponse} "List of projects"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /projects [get]
func (pc *ProjectController) List(c echo.Context) error {
	var request requests.DefaultRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	organizationUUID, err := request.GetUUIDQueryParam(c, "organization_uuid", true)
	if err != nil {
		return response.BadRequestResponse(c, "Invalid organization ID")
	}

	paginationParams := request.ExtractPaginationParams(c)
	projects, err := pc.projectService.List(paginationParams, organizationUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, resources.ProjectResourceCollection(projects))
}

// Show details of a single project
//
// @Summary Show details of a single project
// @Description Get details of a specific project
// @Tags Projects
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param projectUUID path string true "Project UUID"
//
// @Success 200 {object} responses.Response{content=resources.ProjectResponse} "Project details"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /projects/{projectUUID} [get]
func (pc *ProjectController) Show(c echo.Context) error {
	var request requests.DefaultRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	projectUUID, err := request.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	project, err := pc.projectService.GetByUUID(projectUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, resources.ProjectResource(&project))
}

// Store creates a new project
//
// @Summary Create a new project
// @Description Create a new project for a specific organization
// @Tags Projects
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param organizationUUID query string true "Organization UUID"
// @Param name body project_requests.CreateRequest true "Project name"
//
// @Success 201 {object} responses.Response{content=resources.ProjectResponse} "Project details"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /projects [post]
func (pc *ProjectController) Store(c echo.Context) error {
	var request project_requests.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	project, err := pc.projectService.Create(&request, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, resources.ProjectResource(&project))
}

// Update a project
//
// @Summary Update a project
// @Description Update a project for a specific organization
// @Tags Projects
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param projectUUID path string true "Project UUID"
// @Param name body project_requests.UpdateRequest true "Project name"
//
// @Success 200 {object} responses.Response{content=resources.ProjectResponse} "Project details"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /projects/{projectUUID} [put]
func (pc *ProjectController) Update(c echo.Context) error {
	var request project_requests.UpdateRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	projectUUID, err := request.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	updatedOrganization, err := pc.projectService.Update(projectUUID, authUser, &request)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, resources.ProjectResource(updatedOrganization))
}

// Delete a project
//
// @Summary Delete a project
// @Description Remove a project from the organization
// @Tags Projects
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param projectUUID path string true "Project UUID"
//
// @Success 200 {object} responses.Response{} "Project deleted"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /projects/{projectUUID} [delete]
func (pc *ProjectController) Delete(c echo.Context) error {
	var request requests.DefaultRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	projectUUID, err := request.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	if _, err := pc.projectService.Delete(projectUUID, authUser); err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.DeletedResponse(c, nil)
}
