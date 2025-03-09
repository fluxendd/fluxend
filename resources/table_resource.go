package resources

import (
	"fluxton/models"
)

type TableResponse struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Schema string `json:"schema"`
}

func TableResource(table *models.Table) TableResponse {
	return TableResponse{
		Id:     table.Id,
		Name:   table.Name,
		Schema: table.Schema,
	}
}

func TableResourceCollection(tables []models.Table) []TableResponse {
	resourceTables := make([]TableResponse, len(tables))
	for i, table := range tables {
		resourceTables[i] = TableResource(&table)
	}

	return resourceTables
}
