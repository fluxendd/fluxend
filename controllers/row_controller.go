package controllers

import (
	"fluxton/models"
	"fluxton/requests"
	"fluxton/resources"
	"fluxton/responses"
	"fluxton/services"
	"fluxton/utils"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type RowController struct {
	rowService services.RowService
}

func NewRowController(injector *do.Injector) (*RowController, error) {
	rowService := do.MustInvoke[services.RowService](injector)

	return &RowController{rowService: rowService}, nil
}

func (rc *RowController) Store(c echo.Context) error {
	var request requests.RowCreateRequest
	authenticatedUser, _ := utils.NewAuth(c).User()

	if err := request.BindAndValidate(c); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	projectID, err := utils.GetUintPathParam(c, "projectID", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	_, err = rc.rowService.Create(&request, projectID, c.Param("tableName"), authenticatedUser)
	if err != nil {
		return responses.UnprocessableResponse(c, []string{err.Error()})
	}

	return responses.CreatedResponse(c, resources.TableResource(&models.Table{}))
}
