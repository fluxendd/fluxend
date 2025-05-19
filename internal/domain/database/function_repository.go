package database

type FunctionRepository interface {
	List(schema string) ([]Function, error)
	Create(functionSQL string) error
	GetByName(schema, functionName string) (Function, error)
	Delete(schema, functionName string) error
}
