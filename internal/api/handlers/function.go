package handlers

import (
	"fluxton/internal/api/dto"
	"fluxton/internal/api/dto/database"
	functionMapper "fluxton/internal/api/mapper/database"
	"fluxton/internal/api/response"
	functionDomain "fluxton/internal/domain/database"
	"fluxton/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type FunctionHandler struct {
	functionService functionDomain.IndexService
}

func NewFunctionHandler(injector *do.Injector) (*FunctionHandler, error) {
	functionService := do.MustInvoke[functionDomain.IndexService](injector)

	return &FunctionHandler{functionService: functionService}, nil
}

// List retrieves all functions for a schema
//
// @Summary List all functions
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
// @Success 200 {array} responses.Response{content=[]resources.FunctionResponse} "List of functions"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /functions/{schema} [get]
func (fc *FunctionHandler) List(c echo.Context) error {
	var request dto.DefaultRequestWithProjectHeader
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	schema := c.Param("schema")
	if schema == "" {
		return response.BadRequestResponse(c, "Schema is required")
	}

	functions, err := fc.functionService.List(schema, request.ProjectUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, functionMapper.ToResourceCollection(functions))
}

// Show retrieves details of a specific function
//
// @Summary Show details of a single function
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
// @Success 200 {object} responses.Response{content=resources.FunctionResponse} "Function details"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /functions/{schema}/{functionName} [get]
func (fc *FunctionHandler) Show(c echo.Context) error {
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

	fetchedFunction, err := fc.functionService.GetByName(functionName, schema, request.ProjectUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, functionMapper.ToResource(&fetchedFunction))
}

// Store creates a new function
//
// @Summary Create a new function
// @Description Add a new function for specific schema
// @Tags Functions
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param Header X-Project header string true "Project UUID"
//
// @Param form body requests.CreateFunctionRequest true "Function details"
//
// @Success 201 {object} responses.Response{content=resources.FunctionResponse} "Function created"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /functions/{schema} [post]
func (fc *FunctionHandler) Store(c echo.Context) error {
	var request database.CreateFunctionRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	schema := c.Param("schema")
	if schema == "" {
		return response.BadRequestResponse(c, "Schema is required")
	}

	createdFunction, err := fc.functionService.Create(schema, &request, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, functionMapper.ToResource(&createdFunction))
}

// Delete removes a function
//
// @Summary Delete a function
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
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /functions/{schema}/{functionName} [delete]
func (fc *FunctionHandler) Delete(c echo.Context) error {
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

	if _, err := fc.functionService.Delete(schema, functionName, request.ProjectUUID, authUser); err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.DeletedResponse(c, nil)
}
