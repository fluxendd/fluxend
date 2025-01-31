package controllers

import (
	"fluxton/requests"
	"fluxton/resources"
	"fluxton/responses"
	"fluxton/services"
	"fluxton/utils"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type TableController struct {
	tableService services.TableService
}

func NewTableController(injector *do.Injector) (*TableController, error) {
	tableService := do.MustInvoke[services.TableService](injector)

	return &TableController{tableService: tableService}, nil
}

/*func (pc *TableController) List(c echo.Context) error {
	var request requests.DefaultRequest
	authenticatedUserId, _ := utils.NewAuth(c).Id()

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	paginationParams := utils.ExtractPaginationParams(c)
	projects, err := pc.tableService.List(paginationParams, request.OrganizationID, authenticatedUserId)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.ProjectResourceCollection(projects))
}

func (pc *TableController) Show(c echo.Context) error {
	var request requests.DefaultRequest
	authenticatedUser, _ := utils.NewAuth(c).User()

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	id, err := utils.GetUintPathParam(c, "id", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	project, err := pc.tableService.GetByID(id, request.OrganizationID, authenticatedUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.ProjectResource(&project))
}*/

func (pc *TableController) Store(c echo.Context) error {
	var request requests.TableCreateRequest
	authenticatedUser, _ := utils.NewAuth(c).User()

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	projectID, err := utils.GetUintPathParam(c, "projectID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	table, err := pc.tableService.Create(&request, projectID, authenticatedUser)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.TableResource(&table))
}

/*func (pc *TableController) Update(c echo.Context) error {
	var request requests.ProjectCreateRequest
	authenticatedUser, _ := utils.NewAuth(c).User()

	id, err := utils.GetUintPathParam(c, "id", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	updatedOrganization, err := pc.tableService.Update(id, authenticatedUser, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.ProjectResource(updatedOrganization))
}

func (pc *TableController) Delete(c echo.Context) error {
	var request requests.DefaultRequest
	authenticatedUser, _ := utils.NewAuth(c).User()

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	id, err := utils.GetUintPathParam(c, "id", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if _, err := pc.tableService.Delete(id, authenticatedUser); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}
*/
