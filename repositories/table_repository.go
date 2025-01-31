package repositories

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
)

type TableRepository struct {
	db *sqlx.DB
}

type Field struct {
	Name    string
	Type    string
	Primary bool
	Unique  bool
	NotNull bool
	Default string
}

func NewTableRepository(injector *do.Injector) (*TableRepository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)

	return &TableRepository{db: db}, nil
}

func (r *TableRepository) Create(name string, fields map[string]Field) error {
	query := "CREATE TABLE " + name + " ("

	for fieldName, field := range fields {
		query += fieldName + " " + field.Type

		if field.Primary {
			query += " PRIMARY KEY"
		}

		if field.Unique {
			query += " UNIQUE"
		}

		if field.NotNull {
			query += " NOT NULL"
		}

		if field.Default != "" {
			query += " DEFAULT " + field.Default
		}

		query += ", "
	}

	query = query[:len(query)-2] + ")"

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *TableRepository) List() ([]string, error) {
	var tables []string
	err := r.db.Select(&tables, "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		return []string{}, err
	}

	return tables, nil
}

func (r *TableRepository) Exists(name string) (bool, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1", name)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *TableRepository) Drop(name string) error {
	_, err := r.db.Exec("DROP TABLE IF EXISTS " + name)
	if err != nil {
		return err
	}

	return nil
}

func (r *TableRepository) Rename(oldName string, newName string) error {
	_, err := r.db.Exec(fmt.Sprintf("ALTER TABLE %s RENAME TO %s", oldName, newName))
	if err != nil {
		return err
	}

	return nil
}
