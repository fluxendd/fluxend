package database

type Repository interface {
	CreateMany(tableName string, columns []Column, values [][]string) error
}
