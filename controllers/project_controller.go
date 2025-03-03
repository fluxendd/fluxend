package controllers

import (
	"fluxton/requests/project_requests"
	"fluxton/resources"
	"fluxton/responses"
	"fluxton/services"
	"fluxton/utils"
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
// @Tags projects
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param organizationUUID query string true "Organization UUID"
//
// @Success 200 {object} responses.Response{content=[]resources.ProjectResponse} "List of projects"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /projects [get]
func (pc *ProjectController) List(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	organizationUUID, err := utils.GetUUIDQueryParam(c, "organization_uuid", true)
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

// Show details of a single project
//
// @Summary Show details of a single project
// @Description Get details of a specific project
// @Tags projects
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param projectID path string true "Project ID"
//
// @Success 200 {object} responses.Response{content=resources.ProjectResponse} "Project details"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /projects/{projectID} [get]
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

// Store creates a new project
//
// @Summary Create a new project
// @Description Create a new project for a specific organization
// @Tags projects
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
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	project, err := pc.projectService.Create(&request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.ProjectResource(&project))
}

// Update a project
//
// @Summary Update a project
// @Description Update a project for a specific organization
// @Tags projects
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param projectID path string true "Project ID"
// @Param name body project_requests.UpdateRequest true "Project name"
//
// @Success 200 {object} responses.Response{content=resources.ProjectResponse} "Project details"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /projects/{projectID} [put]
func (pc *ProjectController) Update(c echo.Context) error {
	var request project_requests.UpdateRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	projectID, err := utils.GetUUIDPathParam(c, "projectID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	updatedOrganization, err := pc.projectService.Update(projectID, authUser, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.ProjectResource(updatedOrganization))
}

// Delete a project
//
// @Summary Delete a project
// @Description Remove a project from the organization
// @Tags projects
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param projectID path string true "Project ID"
//
// @Success 200 {object} responses.Response{} "Project deleted"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /projects/{projectID} [delete]
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
