package controllers

import (
	"fluxton/internal/api/response"
	"fluxton/pkg/auth"
	"fluxton/requests"
	"fluxton/resources"
	"fluxton/services"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type FunctionController struct {
	functionService services.FunctionService
}

func NewFunctionController(injector *do.Injector) (*FunctionController, error) {
	functionService := do.MustInvoke[services.FunctionService](injector)

	return &FunctionController{functionService: functionService}, nil
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
func (fc *FunctionController) List(c echo.Context) error {
	var request requests.DefaultRequestWithProjectHeader
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

	return response.SuccessResponse(c, resources.FunctionResourceCollection(functions))
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
func (fc *FunctionController) Show(c echo.Context) error {
	var request requests.DefaultRequestWithProjectHeader
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

	function, err := fc.functionService.GetByName(functionName, schema, request.ProjectUUID, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.SuccessResponse(c, resources.FunctionResource(&function))
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
func (fc *FunctionController) Store(c echo.Context) error {
	var request requests.CreateFunctionRequest
	if err := request.BindAndValidate(c); err != nil {
		return response.UnprocessableResponse(c, err)
	}

	authUser, _ := auth.NewAuth(c).User()

	schema := c.Param("schema")
	if schema == "" {
		return response.BadRequestResponse(c, "Schema is required")
	}

	function, err := fc.functionService.Create(schema, &request, authUser)
	if err != nil {
		return response.ErrorResponse(c, err)
	}

	return response.CreatedResponse(c, resources.FunctionResource(&function))
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
func (fc *FunctionController) Delete(c echo.Context) error {
	var request requests.DefaultRequestWithProjectHeader
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
