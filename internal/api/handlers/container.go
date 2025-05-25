package handlers

import (
	"fluxend/internal/api/dto"
	containerDto "fluxend/internal/api/dto/storage/container"
	"fluxend/internal/api/mapper"
	"fluxend/internal/api/response"
	"fluxend/internal/domain/storage/container"
	"fluxend/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type ContainerHandler struct {
	containerService container.Service
}

func NewContainerHandler(injector *do.Injector) (*ContainerHandler, error) {
	containerService := do.MustInvoke[container.Service](injector)

	return &ContainerHandler{containerService: containerService}, nil
}

// List retrieves all container
//
// @Summary List all container
// @Description Retrieve a list of container in a specified project.
// @Tags Containers
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param page query string false "Page number for pagination"
// @Param limit query string false "Number of items per page"
// @Param sort query string false "Field to sort by"
// @Param order query string false "Sort order (asc or desc)"
//
// @Success 200 {object} response.Response{content=[]container.Response} "List of container"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /storage [get]
func (ch *ContainerHandler) List(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	paginationParams := request.ExtractPaginationParams(c)
	containers, err := ch.containerService.List(paginationParams, request.ProjectUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToContainerResourceCollection(containers))
}

// Show retrieves details of a specific container.
//
// @Summary Show details of a single container
// @Description Get details of a specific container
// @Tags Containers
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param containerUUID path string true "Container UUID"
//
// @Success 200 {object} response.Response{content=container.Response} "Container details"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /storage/containers/{containerUUID} [get]
func (ch *ContainerHandler) Show(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader

	authUser, _ := auth.NewAuth(c).User()

	containerUUID, err := request.GetUUIDPathParam(c, "containerUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	fetchedContainer, err := ch.containerService.GetByUUID(containerUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToContainerResource(&fetchedContainer))
}

// Store creates a new container
//
// @Summary Create a new container
// @Description Add a new container to a project
// @Tags Containers
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
// @Param container body container.CreateRequest true "Container details"
//
// @Success 201 {object} response.Response{content=container.Response} "Container created"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /storage [post]
func (ch *ContainerHandler) Store(c echo.Context) error {
	var request containerDto.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	fetchedContainer, err := ch.containerService.Create(containerDto.ToCreateContainerInput(&request), authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, mapper.ToContainerResource(&fetchedContainer))
}

// Update a container
//
// @Summary Update a container
// @Description Modify an existing container's details
// @Tags Containers
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param containerUUID path string true "Container UUID"
// @Param container body container.CreateRequest true "Container details"
//
// @Success 200 {object} response.Response{content=container.Response} "Container updated"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /storage/containers/{containerUUID} [put]
func (ch *ContainerHandler) Update(c echo.Context) error {
	var request containerDto.CreateRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	containerUUID, err := request.GetUUIDPathParam(c, "containerUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	updatedContainer, err := ch.containerService.Update(containerUUID, authUser, containerDto.ToCreateContainerInput(&request))
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToContainerResource(updatedContainer))
}

// Delete a container
//
// @Summary Delete a container
// @Description Remove a container from a project
// @Tags Containers
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param X-Project header string true "Project UUID"
//
// @Param containerUUID path string true "Container UUID"
//
// @Success 204 "Container deleted"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /storage/containers/{containerUUID} [delete]
func (ch *ContainerHandler) Delete(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	containerUUID, err := request.GetUUIDPathParam(c, "containerUUID", true)
	if err != nil {
		return response.BadRequestResponse(c, err.Error())
	}

	if _, err := ch.containerService.Delete(containerUUID, authUser); err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.DeletedResponse(c, nil)
}
