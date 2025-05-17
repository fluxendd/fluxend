package column

import (
	columnDto "fluxton/internal/api/dto/database/column"
	columnDomain "fluxton/internal/domain/database/column"
)

func ToResource(column *columnDomain.Column) columnDto.Response {
	return columnDto.Response{
		Name:            column.Name,
		Position:        column.Position,
		NotNull:         column.NotNull,
		Type:            column.Type,
		Default:         column.Default,
		Primary:         column.Primary,
		Unique:          column.Unique,
		Foreign:         column.Foreign,
		ReferenceTable:  column.ReferenceTable,
		ReferenceColumn: column.ReferenceColumn,
	}
}

func ToResourceCollection(columns []columnDomain.Column) []columnDto.Response {
	resourceColumns := make([]columnDto.Response, len(columns))
	for i, currentColumn := range columns {
		resourceColumns[i] = ToResource(&currentColumn)
	}

	return resourceColumns
}
