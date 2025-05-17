package table

import (
	"fluxton/internal/domain/database/table"
)

type TableResponse struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Schema        string `json:"schema"`
	EstimatedRows int    `json:"estimatedRows"`
	TotalSize     string `json:"totalSize"`
}

func TableResource(table *table.Table) TableResponse {
	return TableResponse{
		Id:            table.Id,
		Name:          table.Name,
		Schema:        table.Schema,
		EstimatedRows: table.EstimatedRows,
		TotalSize:     table.TotalSize,
	}
}

func TableResourceCollection(tables []table.Table) []TableResponse {
	resourceTables := make([]TableResponse, len(tables))
	for i, table := range tables {
		resourceTables[i] = TableResource(&table)
	}

	return resourceTables
}
