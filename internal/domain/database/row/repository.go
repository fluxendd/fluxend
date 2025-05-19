package row

import (
	"fluxton/internal/domain/database/column"
)

type Repository interface {
	CreateMany(tableName string, columns []column.Column, values [][]string) error
}
