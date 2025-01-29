package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"myapp/requests"
	"myapp/resources"
	"myapp/responses"
	"myapp/services"
	"myapp/utils"
)

type NoteController struct {
	noteService services.NoteService
}

func NewNoteController(injector *do.Injector) (*NoteController, error) {
	noteService := do.MustInvoke[services.NoteService](injector)

	return &NoteController{noteService: noteService}, nil
}

func (nc *NoteController) List(c echo.Context) error {
	authenticatedUserId, _ := utils.NewAuth(c).Id()

	paginationParams := utils.ExtractPaginationParams(c)
	notes, err := nc.noteService.List(paginationParams, authenticatedUserId)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.NoteResourceCollection(notes))
}

func (nc *NoteController) Show(c echo.Context) error {
	authenticatedUserId, _ := utils.NewAuth(c).Id()

	id, err := utils.GetUintPathParam(c, "id", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	note, err := nc.noteService.GetByID(id, authenticatedUserId)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.NoteResource(&note))
}

func (nc *NoteController) Store(c echo.Context) error {
	var request requests.NoteCreateRequest
	authenticatedUserId, _ := utils.NewAuth(c).Id()

	if err := c.Bind(&request); err != nil {
		return responses.BadRequestResponse(c, "note.error.invalidPayload")
	}

	if err := request.Validate(); err != nil {
		return responses.UnprocessableResponse(c, err)
	}

	note, err := nc.noteService.Create(&request, authenticatedUserId)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.CreatedResponse(c, resources.NoteResource(&note))
}

func (nc *NoteController) Update(c echo.Context) error {
	var request requests.NoteCreateRequest
	authenticatedUserId, _ := utils.NewAuth(c).Id()

	id, err := utils.GetUintPathParam(c, "id", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if err := c.Bind(&request); err != nil {
		return responses.BadRequestResponse(c, "note.error.invalidPayload")
	}

	updatedNote, err := nc.noteService.Update(id, authenticatedUserId, &request)
	if err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.SuccessResponse(c, resources.NoteResource(updatedNote))
}

func (nc *NoteController) Delete(c echo.Context) error {
	authenticatedUserId, _ := utils.NewAuth(c).Id()
	id, err := utils.GetUintPathParam(c, "id", true)
	if err != nil {
		return responses.BadRequestResponse(c, err.Error())
	}

	if _, err := nc.noteService.Delete(id, authenticatedUserId); err != nil {
		return responses.ErrorResponse(c, err)
	}

	return responses.DeletedResponse(c, nil)
}
