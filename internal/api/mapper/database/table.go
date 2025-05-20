package database

import (
	databaseDto "fluxton/internal/api/dto/database"
	databaseDomain "fluxton/internal/domain/database"
)

func ToTableResource(table *databaseDomain.Table) databaseDto.TableResponse {
	return databaseDto.TableResponse{
		Id:            table.Id,
		Name:          table.Name,
		Schema:        table.Schema,
		EstimatedRows: table.EstimatedRows,
		TotalSize:     table.TotalSize,
	}
}

func ToTableResourceCollection(tables []databaseDomain.Table) []databaseDto.TableResponse {
	resourceTables := make([]databaseDto.TableResponse, len(tables))
	for i, currentTable := range tables {
		resourceTables[i] = ToTableResource(&currentTable)
	}

	return resourceTables
}
