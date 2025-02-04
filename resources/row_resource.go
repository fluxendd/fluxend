package resources

import (
	"fluxton/models"
)

type RowResponse map[string]interface{}

func RowResource(row models.Row) RowResponse {
	resource := RowResponse{}
	for key, value := range row {
		resource[key] = value
	}
	return resource
}

func RowResourceCollection(rows []models.Row) []RowResponse {
	resourceRows := make([]RowResponse, len(rows))
	for i, row := range rows {
		resourceRows[i] = RowResource(row)
	}
	return resourceRows
}
