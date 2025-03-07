package controllers

import (
	"fluxton/errs"
	"fluxton/requests"
	"fluxton/resources"
	"fluxton/responses"
	"fluxton/services"
	"fluxton/utils"
	"github.com/google/uuid"
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
// @Tags Schema
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param projectUUID path string true "Project UUID"
// @Param schema path string true "Schema to search under"
//
// @Success 200 {array} responses.Response{content=[]resources.FunctionResponse} "List of functions"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /projects/{projectUUID}/functions/{schema} [get]
func (fc *FunctionController) List(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	projectUUID, schema, err := fc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	functions, err := fc.functionService.List(schema, projectUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FunctionResourceCollection(functions))
}

// Show retrieves details of a specific function
//
// @Summary Show details of a single function
// @Description Get details of a specific function
// @Tags forms
//
// @Accept json
// @Produce json
//
// @Param Authorization header string true "Bearer Token"
// @Param projectUUID path string true "Project UUID"
// @Param schema path string true "Schema name"
// @Param functionName path string true "Function name"
//
// @Success 200 {object} responses.Response{content=resources.FunctionResponse} "Function details"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /projects/{projectUUID}/functions/{schema}/{functionName} [get]
func (fc *FunctionController) Show(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	projectUUID, schema, err := fc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	functionName := c.Param("functionName")
	if functionName == "" {
		return responses.BadRequestResponse(c, "Function name is required")
	}

	function, err := fc.functionService.GetByName(functionName, schema, projectUUID, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.FunctionResource(&function))
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
// @Param form body requests.CreateFunctionRequest true "Function details"
//
// @Success 201 {object} responses.Response{content=resources.FunctionResponse} "Function created"
// @Failure 422 "Unprocessable entity"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /projects/{projectUUID}/functions/{schema} [post]
func (fc *FunctionController) Store(c echo.Context) error {
	var request requests.CreateFunctionRequest
	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	authUser, _ := utils.NewAuth(c).User()

	projectUUID, schema, err := fc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	function, err := fc.functionService.Create(projectUUID, schema, &request, authUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.FunctionResource(&function))
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
// @Param projectUUID path string true "Project UUID"
// @Param schema path string true "Schema name"
// @Param functionName path string true "Function name"
//
// @Success 204 "Form deleted"
// @Failure 400 "Invalid input"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
//
// @Router /projects/{projectUUID}/functions/{schema}/{functionName} [delete]
func (fc *FunctionController) Delete(c echo.Context) error {
	authUser, _ := utils.NewAuth(c).User()

	projectUUID, schema, err := fc.parseRequest(c)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	functionName := c.Param("functionName")
	if functionName == "" {
		return responses.BadRequestResponse(c, "Function name is required")
	}

	if _, err := fc.functionService.Delete(functionName, schema, projectUUID, authUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}

func (fc *FunctionController) parseRequest(c echo.Context) (uuid.UUID, string, error) {
	projectUUID, err := utils.GetUUIDPathParam(c, "projectUUID", true)
	if err != nil {
		return uuid.UUID{}, "", err
	}

	schema := c.QueryParam("schema")
	if schema == "" {
		return uuid.UUID{}, "", errs.NewBadRequestError("Schema is required")
	}

	return projectUUID, schema, nil
}
