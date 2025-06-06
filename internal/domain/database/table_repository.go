package database

type TableRepository interface {
	Exists(name string) (bool, error)
	Create(name string, columns []Column) error
	Duplicate(existingTable string, newTable string) error
	List() ([]Table, error)
	GetByNameInSchema(schema, name string) (Table, error)
	DropIfExists(name string) error
	Rename(oldName string, newName string) error
}
