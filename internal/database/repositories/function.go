package repositories

import (
	"fluxend/internal/domain/database"
	"fluxend/internal/domain/shared"
	"fmt"
	"github.com/lib/pq"
	"github.com/samber/do"
)

type FunctionRepository struct {
	db shared.DB
}

func NewFunctionRepository(injector *do.Injector) (*FunctionRepository, error) {
	db := do.MustInvoke[shared.DB](injector)
	return &FunctionRepository{db: db}, nil
}

func (r *FunctionRepository) List(schema string) ([]database.Function, error) {
	var functions []database.Function
	query := `
       SELECT routine_name, routine_type, data_type, type_udt_name, routine_definition, external_language, sql_data_access
       FROM information_schema.routines 
       WHERE routine_type = 'FUNCTION' AND specific_schema = $1`

	return functions, r.db.SelectList(&functions, query, schema)
}

func (r *FunctionRepository) Create(functionSQL string) error {
	_, err := r.db.ExecWithRowsAffected(functionSQL)
	return err
}

func (r *FunctionRepository) GetByName(schema, functionName string) (database.Function, error) {
	var function database.Function
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

	return function, r.db.Get(&function, query, schema, functionName)
}

func (r *FunctionRepository) Delete(schema, functionName string) error {
	query := fmt.Sprintf(`DROP FUNCTION IF EXISTS %s.%s CASCADE`, pq.QuoteIdentifier(schema), pq.QuoteIdentifier(functionName))
	_, err := r.db.ExecWithRowsAffected(query)
	return err
}
