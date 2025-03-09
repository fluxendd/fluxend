package resources

import (
	"fluxton/models"
)

type TableResponse struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Schema        string `json:"schema"`
	EstimatedRows int    `json:"estimatedRows"`
	TotalSize     string `json:"totalSize"`
}

func TableResource(table *models.Table) TableResponse {
	return TableResponse{
		Id:            table.Id,
		Name:          table.Name,
		Schema:        table.Schema,
		EstimatedRows: table.EstimatedRows,
		TotalSize:     table.TotalSize,
	}
}

func TableResourceCollection(tables []models.Table) []TableResponse {
	resourceTables := make([]TableResponse, len(tables))
	for i, table := range tables {
		resourceTables[i] = TableResource(&table)
	}

	return resourceTables
}
