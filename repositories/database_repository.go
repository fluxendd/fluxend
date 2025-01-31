package repositories

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
)

type DatabaseRepository struct {
	db *sqlx.DB
}

func NewDatabaseRepository(injector *do.Injector) (*DatabaseRepository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)

	return &DatabaseRepository{db: db}, nil
}

func (r *DatabaseRepository) Create(name string) error {
	_, err := r.db.Exec(fmt.Sprintf("CREATE DATABASE %s", name))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// Drop
// Recreate
// List
// GetSchema
// Exists
// Connect
