package table

import (
	tableDto "fluxton/internal/api/dto/database/table"
	tableDomain "fluxton/internal/domain/database"
)

func ToResource(table *tableDomain.Table) tableDto.Response {
	return tableDto.Response{
		Id:            table.Id,
		Name:          table.Name,
		Schema:        table.Schema,
		EstimatedRows: table.EstimatedRows,
		TotalSize:     table.TotalSize,
	}
}

func ToResourceCollection(tables []tableDomain.Table) []tableDto.Response {
	resourceTables := make([]tableDto.Response, len(tables))
	for i, currentTable := range tables {
		resourceTables[i] = ToResource(&currentTable)
	}

	return resourceTables
}
