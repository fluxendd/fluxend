package database

type RowRepository interface {
	CreateMany(tableName string, columns []Column, values [][]string) error
}
