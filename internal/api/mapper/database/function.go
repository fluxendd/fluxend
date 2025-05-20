package database

import (
	databaseDto "fluxton/internal/api/dto/database"
	databaseDomain "fluxton/internal/domain/database"
)

func ToFunctionResource(function *databaseDomain.Function) databaseDto.FunctionResponse {
	return databaseDto.FunctionResponse{
		Name:       function.Name,
		Type:       function.Type,
		DataType:   function.DataType,
		Definition: function.Definition,
		Language:   function.Language,
	}
}

func ToFunctionResourceCollection(functions []databaseDomain.Function) []databaseDto.FunctionResponse {
	resourceFunctions := make([]databaseDto.FunctionResponse, len(functions))
	for i, function := range functions {
		resourceFunctions[i] = ToFunctionResource(&function)
	}

	return resourceFunctions
}
