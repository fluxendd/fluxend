package table

import (
	"fluxton/internal/domain/database/column"
)

type Repository interface {
	Exists(name string) (bool, error)
	Create(name string, columns []column.Column) error
	Duplicate(existingTable string, newTable string) error
	List() ([]Table, error)
	GetByNameInSchema(schema, name string) (Table, error)
	DropIfExists(name string) error
	Rename(oldName string, newName string) error
}
