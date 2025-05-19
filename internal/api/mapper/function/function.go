package function

import (
	functionDto "fluxton/internal/api/dto/database/function"
	functionDomain "fluxton/internal/domain/database"
)

func ToResource(function *functionDomain.Function) functionDto.Response {
	return functionDto.Response{
		Name:       function.Name,
		Type:       function.Type,
		DataType:   function.DataType,
		Definition: function.Definition,
		Language:   function.Language,
	}
}

func ToResourceCollection(functions []functionDomain.Function) []functionDto.Response {
	resourceFunctions := make([]functionDto.Response, len(functions))
	for i, function := range functions {
		resourceFunctions[i] = ToResource(&function)
	}

	return resourceFunctions
}
