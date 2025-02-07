package repositories

import (
	"fluxton/types"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ClientTableRepository struct {
	connection *sqlx.DB
}

func NewClientTableRepository(connection *sqlx.DB) (*ClientTableRepository, error) {
	return &ClientTableRepository{connection: connection}, nil
}

func (r *ClientTableRepository) Create(name string, columns []types.TableColumn) error {
	query := "CREATE TABLE " + name + " ("

	for _, column := range columns {
		query += column.Name + " " + column.Type

		if column.Primary {
			query += " PRIMARY KEY"
		}

		if column.Unique {
			query += " UNIQUE"
		}

		if column.NotNull {
			query += " NOT NULL"
		}

		if column.Default != "" {
			query += " DEFAULT " + column.Default
		}

		query += ", "
	}

	query = query[:len(query)-2] + ")"

	_, err := r.connection.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *ClientTableRepository) Duplicate(oldName string, newName string) error {
	_, err := r.connection.Exec(fmt.Sprintf("CREATE TABLE %s AS TABLE %s", newName, oldName))
	if err != nil {
		return err
	}

	return nil
}

func (r *ClientTableRepository) List() ([]string, error) {
	var tables []string
	err := r.connection.Select(&tables, "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		return []string{}, err
	}

	return tables, nil
}

func (r *ClientTableRepository) Exists(name string) (bool, error) {
	var count int
	err := r.connection.Get(&count, "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1", name)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ClientTableRepository) DropIfExists(name string) error {
	_, err := r.connection.Exec("DROP TABLE IF EXISTS " + name)
	if err != nil {
		return err
	}

	return nil
}

func (r *ClientTableRepository) Rename(oldName string, newName string) error {
	_, err := r.connection.Exec(fmt.Sprintf("ALTER TABLE %s RENAME TO %s", oldName, newName))
	if err != nil {
		return err
	}

	return nil
}
