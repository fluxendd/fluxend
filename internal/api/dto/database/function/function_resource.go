package function

import (
	"fluxton/internal/domain/database/function"
)

type FunctionResponse struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	DataType   string `json:"dataType"`
	Definition string `json:"definition"`
	Language   string `json:"language"`
}

func FunctionResource(function *function.Function) FunctionResponse {
	return FunctionResponse{
		Name:       function.Name,
		Type:       function.Type,
		DataType:   function.DataType,
		Definition: function.Definition,
		Language:   function.Language,
	}
}

func FunctionResourceCollection(functions []function.Function) []FunctionResponse {
	resourceFunctions := make([]FunctionResponse, len(functions))
	for i, function := range functions {
		resourceFunctions[i] = FunctionResource(&function)
	}

	return resourceFunctions
}
