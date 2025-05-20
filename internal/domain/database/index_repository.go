package database

type IndexRepository interface {
	GetByName(tableName, indexName string) (string, error)
	Has(tableName, indexName string) (bool, error)
	List(tableName string) ([]string, error)
	Create(tableName string, indexName string, columns []string, isUnique bool) (bool, error)
	DropIfExists(indexName string) (bool, error)
}
