package database

type ColumnRepository interface {
	List(tableName string) ([]Column, error)
	Has(tableName, columnName string) (bool, error)
	HasAny(tableName string, columns []Column) (bool, error)
	HasAll(tableName string, columns []Column) (bool, error)
	CreateOne(tableName string, column Column) error
	CreateMany(tableName string, fields []Column) error
	AlterOne(tableName string, columns []Column) error
	AlterMany(tableName string, fields []Column) error
	Rename(tableName, oldColumnName, newColumnName string) error
	Drop(tableName, columnName string) error
	mapColumnsToNames(columns []Column) []string
	BuildColumnDefinition(column Column) string
	BuildForeignKeyConstraint(tableName string, column Column) (string, bool)
}
