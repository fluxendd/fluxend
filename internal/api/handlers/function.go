package handlers

import (
	"fluxend/internal/api/dto"
	databaseDto "fluxend/internal/api/dto/database"
	"fluxend/internal/api/mapper"
	"fluxend/internal/api/response"
	"fluxend/internal/domain/database"
	"fluxend/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type FunctionHandler struct {
	functionService database.FunctionService
}

func NewFunctionHandler(injector *do.Injector) (*FunctionHandler, error) {
	functionService := do.MustInvoke[database.FunctionService](injector)

	return &FunctionHandler{functionService: functionService}, nil
}

// List retrieves all functions for a schema
//
// @Summary List functions
// @Description Retrieve a list of all functions for the specified schema
// @Tags Functions
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @param Header X-Project header string true "Project UUID"
//
// @Param projectUUID path string true "Project UUID"
// @Param schema path string true "Schema to search under"
//
// @Success 200 {array} response.Response{content=[]database.FunctionResponse} "List of functions"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /functions/{schema} [get]
func (fh *FunctionHandler) List(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	schema := c.Param("schema")
	if schema == "" {
		return response.BadRequestResponse(c, "Schema is required")
	}

	functions, err := fh.functionService.List(schema, request.ProjectUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToFunctionResourceCollection(functions))
}

// Show retrieves details of a specific function
//
// @Summary Retrieve function
// @Description Get details of a specific function
// @Tags Functions
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param Header X-Project header string true "Project UUID"
//
// @Param schema path string true "Schema name"
// @Param functionName path string true "Function name"
//
// @Success 200 {object} response.Response{content=database.FunctionResponse} "Function details"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /functions/{schema}/{functionName} [get]
func (fh *FunctionHandler) Show(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	schema := c.Param("schema")
	if schema == "" {
		return response.BadRequestResponse(c, "Schema is required")
	}

	functionName := c.Param("functionName")
	if functionName == "" {
		return response.BadRequestResponse(c, "Function name is required")
	}

	fetchedFunction, err := fh.functionService.GetByName(functionName, schema, request.ProjectUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, mapper.ToFunctionResource(&fetchedFunction))
}

// Store creates a new function
//
// @Summary Create function
// @Description Add a new function for specific schema
// @Tags Functions
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param Header X-Project header string true "Project UUID"
//
// @Param function body database.CreateFunctionRequest true "Function details"
//
// @Success 201 {object} response.Response{content=database.FunctionResponse} "Function created"
// @Failure 422 {object} response.UnprocessableErrorResponse "Unprocessable input response"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /functions/{schema} [post]
func (fh *FunctionHandler) Store(c echo.Context) error {
	var request databaseDto.CreateFunctionRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	schema := c.Param("schema")
	if schema == "" {
		return response.BadRequestResponse(c, "Schema is required")
	}

	createdFunction, err := fh.functionService.Create(schema, databaseDto.ToCreateFunctionInput(request), authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, mapper.ToFunctionResource(&createdFunction))
}

// Delete removes a function
//
// @Summary Delete function
// @Description Remove a function from the schema
// @Tags Functions
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param Header X-Project header string true "Project UUID"
//
// @Param projectUUID path string true "Project UUID"
// @Param schema path string true "Schema name"
// @Param functionName path string true "Function name"
//
// @Success 204 "Form deleted"
// @Failure 400 {object} response.BadRequestErrorResponse "Bad request response"
// @Failure 401 {object} response.UnauthorizedErrorResponse "Unauthorized response"
// @Failure 500 {object} response.InternalServerErrorResponse "Internal server error response"
//
// @Router /functions/{schema}/{functionName} [delete]
func (fh *FunctionHandler) Delete(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	schema := c.Param("schema")
	if schema == "" {
		return response.BadRequestResponse(c, "Schema is required")
	}

	functionName := c.Param("functionName")
	if functionName == "" {
		return response.BadRequestResponse(c, "Function name is required")
	}

	if _, err := fh.functionService.Delete(schema, functionName, request.ProjectUUID, authUser); err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.DeletedResponse(c, nil)
}
