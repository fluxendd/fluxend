package repositories

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
	"myapp/models"
)

type TagRepository struct {
	db *sqlx.DB
}

func NewTagRepository(injector *do.Injector) (*TagRepository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)

	return &TagRepository{db: db}, nil
}

func (r *TagRepository) Create(tag *models.Tag) (*models.Tag, error) {
	query := "INSERT INTO tags (name) VALUES ($1) RETURNING id"
	err := r.db.QueryRowx(query, tag.Name).Scan(&tag.ID)
	if err != nil {
		return &models.Tag{}, fmt.Errorf("could not create row: %v", err)
	}

	return tag, nil
}
