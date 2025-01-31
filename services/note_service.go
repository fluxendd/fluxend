package services

import (
	"fluxton/errs"
	"fluxton/models"
	"fluxton/policies"
	"fluxton/repositories"
	"fluxton/requests"
	"fluxton/utils"
	"github.com/samber/do"
)

type NoteService interface {
	List(paginationParams utils.PaginationParams, authenticatedUserId uint) ([]models.Note, error)
	GetByID(id, authenticatedUserId uint) (models.Note, error)
	Create(request *requests.NoteCreateRequest, authenticatedUserId uint) (models.Note, error)
	Update(noteId, authenticatedUserId uint, request *requests.NoteCreateRequest) (*models.Note, error)
	Delete(noteId, authenticatedUserId uint) (bool, error)
}

type NoteServiceImpl struct {
	notePolicy policies.NotePolicy
	noteRepo   *repositories.NoteRepository
}

func NewNoteService(injector *do.Injector) (NoteService, error) {
	policy := policies.NewNotePolicy()
	repo := do.MustInvoke[*repositories.NoteRepository](injector)

	return &NoteServiceImpl{
		notePolicy: policy,
		noteRepo:   repo,
	}, nil
}

func (s *NoteServiceImpl) List(paginationParams utils.PaginationParams, authenticatedUserId uint) ([]models.Note, error) {
	return s.noteRepo.List(paginationParams, authenticatedUserId)
}

func (s *NoteServiceImpl) GetByID(id, authenticatedUserId uint) (models.Note, error) {
	return s.noteRepo.GetByID(id, authenticatedUserId)
}

func (s *NoteServiceImpl) Create(request *requests.NoteCreateRequest, authenticatedUserId uint) (models.Note, error) {
	if !s.notePolicy.CanCreate(authenticatedUserId) {
		return models.Note{}, errs.NewForbiddenError("note.error.createForbidden")
	}

	note := models.Note{
		Title:   request.Title,
		Content: request.Content,
		UserId:  authenticatedUserId,
	}

	_, err := s.noteRepo.Create(&note)
	if err != nil {
		return models.Note{}, err
	}

	return note, nil
}

func (s *NoteServiceImpl) Update(noteId, authenticatedUserId uint, request *requests.NoteCreateRequest) (*models.Note, error) {
	note, err := s.noteRepo.GetByID(noteId, authenticatedUserId)
	if err != nil {
		return nil, err
	}

	if !s.notePolicy.CanUpdate(note.UserId, authenticatedUserId) {
		return &models.Note{}, errs.NewForbiddenError("note.error.updateForbidden")
	}

	err = utils.PopulateModel(&note, request)
	if err != nil {
		return nil, err
	}

	return s.noteRepo.Update(noteId, &note)
}

func (s *NoteServiceImpl) Delete(noteId, authenticatedUserId uint) (bool, error) {
	note, err := s.noteRepo.GetByID(noteId, authenticatedUserId)
	if err != nil {
		return false, err
	}

	if !s.notePolicy.CanUpdate(note.UserId, authenticatedUserId) {
		return false, errs.NewForbiddenError("note.error.deleteForbidden")
	}

	return s.noteRepo.Delete(noteId)
}
