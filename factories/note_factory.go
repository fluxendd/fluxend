package factories

import (
	"github.com/samber/do"
	"myapp/utils"
	"time"

	"myapp/models"
	"myapp/repositories"
)

type NoteOption func(*models.Note)

type NoteFactory struct {
	repo *repositories.NoteRepository
}

func NewNoteFactory(injector *do.Injector) (*NoteFactory, error) {
	repo := do.MustInvoke[*repositories.NoteRepository](injector)

	return &NoteFactory{repo: repo}, nil
}

// Create a note with options
func (f *NoteFactory) Create(opts ...NoteOption) (*models.Note, error) {
	note := &models.Note{
		Title:     utils.Faker.Lorem().Sentence(5),
		Content:   utils.Faker.Lorem().Paragraph(2),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for _, opt := range opts {
		opt(note)
	}

	createdNote, err := f.repo.Create(note)
	if err != nil {
		return nil, err
	}

	return createdNote, nil
}

func (f *NoteFactory) CreateMany(count int, opts ...NoteOption) ([]*models.Note, error) {
	var notes []*models.Note
	for i := 0; i < count; i++ {
		note, err := f.Create(opts...)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

func (f *NoteFactory) WithTitle(title string) NoteOption {
	return func(note *models.Note) {
		note.Title = title
	}
}

func (f *NoteFactory) WithUserId(userId uint) NoteOption {
	return func(note *models.Note) {
		note.UserId = userId
	}
}
