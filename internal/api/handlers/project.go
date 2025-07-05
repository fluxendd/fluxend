package handlers

import (
	"fluxend/internal/api/dto"
	projectDto "fluxend/internal/api/dto/project"
	"fluxend/internal/api/mapper"
	"fluxend/internal/api/response"
	"fluxend/internal/domain/openapi"
	"fluxend/internal/domain/project"
	"fluxend/pkg/auth"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"os"
)

type ProjectHandler struct {
	projectService project.Service
	openApiService openapi.Service
}

func NewProjectHandler(injector *do.Injector) (*ProjectHandler, error) {
	projectService := do.MustInvoke[project.Service](injector)
	openApiService := do.MustInvoke[openapi.Service](injector)

	return &ProjectHandler{
		projectService: projectService,
		openApiService: openApiService,
	}, nil
}

// List all projects
//
// @Summary List projects
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
// @Success 200 {object} response.Response{content=[]project.Response} "List of projects"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /projects [get]
func (ph *ProjectHandler) List(c echo.Context) error {
	var request dto.DefaultRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	organizationUUID, err := request.GetUUIDQueryParam(c, "organization_uuid", true)
	if err != nil {
		return response.BadRequestResponse(c, "Invalid organization ID")
	}

	paginationParams := request.ExtractPaginationParams(c)
	projects, err := ph.projectService.List(paginationParams, organizationUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToProjectResourceCollection(projects))
}

// Show details of a single project
//
// @Summary Retrieve project
// @Description Get details of a specific project
// @Tags Projects
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param projectUUID path string true "Project UUID"
//
// @Success 200 {object} response.Response{content=project.Response} "Project details"
// @Failure 422 {object} response.UnprocessableErrorResponse "Unprocessable input response"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /projects/{projectUUID} [get]
func (ph *ProjectHandler) Show(c echo.Context) error {
	var request dto.DefaultRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	projectUUID, err := request.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	fetchedProject, err := ph.projectService.GetByUUID(projectUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToProjectResource(&fetchedProject))
}

// Store creates a new project
//
// @Summary Create project
// @Description Create a new project for a specific organization
// @Tags Projects
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param organizationUUID query string true "Organization UUID"
// @Param name body project.CreateRequest true "Project name"
//
// @Success 201 {object} response.Response{content=project.Response} "Project details"
// @Failure 422 {object} response.UnprocessableErrorResponse "Unprocessable input response"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /projects [post]
func (ph *ProjectHandler) Store(c echo.Context) error {
	var request projectDto.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	updatedProject, err := ph.projectService.Create(projectDto.ToCreateProjectInput(&request), authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, mapper.ToProjectResource(&updatedProject))
}

// Update a project
//
// @Summary Update project
// @Description Update a project for a specific organization
// @Tags Projects
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param projectUUID path string true "Project UUID"
// @Param name body project.UpdateRequest true "Project name"
//
// @Success 200 {object} response.Response{content=project.Response} "Project details"
// @Failure 422 {object} response.UnprocessableErrorResponse "Unprocessable input response"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /projects/{projectUUID} [put]
func (ph *ProjectHandler) Update(c echo.Context) error {
	var request projectDto.UpdateRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	projectUUID, err := request.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	updatedOrganization, err := ph.projectService.Update(projectUUID, authUser, projectDto.ToUpdateProjectInput(&request))
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToProjectResource(updatedOrganization))
}

// Delete a project
//
// @Summary Delete project
// @Description Remove a project from the organization
// @Tags Projects
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param projectUUID path string true "Project UUID"
//
// @Success 200 {object} response.Response{} "Project deleted"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /projects/{projectUUID} [delete]
func (ph *ProjectHandler) Delete(c echo.Context) error {
	var request dto.DefaultRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	projectUUID, err := request.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	if _, err := ph.projectService.Delete(projectUUID, authUser); err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.DeletedResponse(c, nil)
}

// GenerateOpenAPI generate OpenAPI docs for project
//
// @Summary OpenAPI projects
// @Description Generate OpenAPI documentation for a project
// @Tags Projects
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param projectUUID path string true "Project UUID"
//
// @Param tables query string false "Comma-separated list of tables to include in OpenAPI"
//
// @Success 200 {object} response.Response "OpenAPI documentation response"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /projects/{projectUUID}/openapi [get]
func (ph *ProjectHandler) GenerateOpenAPI(c echo.Context) error {
	var request dto.DefaultRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	projectUUID, err := request.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	requestedTables := c.QueryParam("tables")

	openAPIResponse, err := ph.openApiService.Generate(projectUUID, requestedTables, authUser)
	fmt.Println(openAPIResponse)
	os.Exit(1)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, openAPIResponse)
}
