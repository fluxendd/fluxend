package repositories

import (
	"fluxton/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ClientFunctionRepository struct {
	connection *sqlx.DB
}

func NewClientFunctionRepository(connection *sqlx.DB) (*ClientFunctionRepository, error) {
	return &ClientFunctionRepository{connection: connection}, nil
}

func (r *ClientFunctionRepository) ListFunctions(schema string) ([]models.Function, error) {
	var functions []models.Function
	query := `
		SELECT routine_name, routine_type, data_type, type_udt_name, routine_definition, external_language, sql_data_access
		FROM information_schema.routines 
		WHERE routine_type = 'FUNCTION' AND specific_schema = $1`
	err := r.connection.Select(&functions, query, schema)
	if err != nil {
		return nil, err
	}

	return functions, nil
}

func (r *ClientFunctionRepository) Create(functionSQL string) error {
	_, err := r.connection.Exec(functionSQL)

	return err
}

func (r *ClientFunctionRepository) GetByName(schema, functionName string) (models.Function, error) {
	var function models.Function
	query := `
		SELECT 
			r.routine_name, 
			r.routine_type, 
			r.data_type, 
			r.type_udt_name, 
			r.external_language, 
			r.sql_data_access,
			pg_get_functiondef(p.oid) AS routine_definition
		FROM information_schema.routines r
		JOIN pg_proc p ON r.routine_name = p.proname
		JOIN pg_namespace n ON p.pronamespace = n.oid
		WHERE r.specific_schema = $1 AND r.routine_name = $2 AND n.nspname = $1`

	err := r.connection.Get(&function, query, schema, functionName)
	if err != nil {
		return models.Function{}, err
	}
	return function, nil
}

func (r *ClientFunctionRepository) Update(schema, functionName, newFunctionSQL string) error {
	err := r.DeleteFunction(schema, functionName)
	if err != nil {
		return err
	}

	// Recreate function
	return r.Create(newFunctionSQL)
}

func (r *ClientFunctionRepository) DeleteFunction(schema, functionName string) error {
	query := fmt.Sprintf(`DROP FUNCTION IF EXISTS %s.%s CASCADE`, schema, functionName)
	_, err := r.connection.Exec(query)
	if err != nil {
		return err
	}

	return err
}
