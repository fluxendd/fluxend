package factories

import (
	"github.com/samber/do"
	"myapp/utils"
	"time"

	"myapp/models"
	"myapp/repositories"
)

type TagOption func(*models.Tag)

type TagFactory struct {
	repo *repositories.TagRepository
}

func NewTagFactory(injector *do.Injector) (*TagFactory, error) {
	repo := do.MustInvoke[*repositories.TagRepository](injector)

	return &TagFactory{repo: repo}, nil
}

// Create a tag with options
func (f *TagFactory) Create(opts ...TagOption) (*models.Tag, error) {
	tag := &models.Tag{
		Name:      utils.Faker.Lorem().Word(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for _, opt := range opts {
		opt(tag)
	}

	createdTag, err := f.repo.Create(tag)
	if err != nil {
		return nil, err
	}

	return createdTag, nil
}

func (f *TagFactory) CreateMany(count int, opts ...TagOption) ([]*models.Tag, error) {
	var tags []*models.Tag
	for i := 0; i < count; i++ {
		tag, err := f.Create(opts...)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (f *TagFactory) CreateWithName(name string) (*models.Tag, error) {
	opts := func(tag *models.Tag) {
		tag.Name = name
	}

	return f.Create(opts)
}
