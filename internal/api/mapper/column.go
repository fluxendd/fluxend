package mapper

import (
	databaseDto "fluxton/internal/api/dto/database"
	databaseDomain "fluxton/internal/domain/database"
)

func ToColumnResource(column *databaseDomain.Column) databaseDto.ColumnResponse {
	return databaseDto.ColumnResponse{
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

func ToColumnResourceCollection(columns []databaseDomain.Column) []databaseDto.ColumnResponse {
	resourceColumns := make([]databaseDto.ColumnResponse, len(columns))
	for i, currentColumn := range columns {
		resourceColumns[i] = ToColumnResource(&currentColumn)
	}

	return resourceColumns
}
