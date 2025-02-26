package resources

import (
	"fluxton/models"
	"fluxton/types"
	"github.com/google/uuid"
)

type TableResponse struct {
	Uuid        uuid.UUID           `json:"uuid"`
	ProjectUuid uuid.UUID           `json:"project_uuid"`
	CreatedBy   uuid.UUID           `json:"created_by"`
	UpdatedBy   uuid.UUID           `json:"updated_by"`
	Name        string              `json:"name"`
	Columns     []types.TableColumn `json:"columns"`
	CreatedAt   string              `json:"created_at"`
	UpdatedAt   string              `json:"updated_at"`
}

func TableResource(table *models.Table) TableResponse {
	return TableResponse{
		Uuid:        table.Uuid,
		ProjectUuid: table.ProjectUuid,
		CreatedBy:   table.CreatedBy,
		UpdatedBy:   table.UpdatedBy,
		Name:        table.Name,
		Columns:     table.Columns,
		CreatedAt:   table.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   table.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func TableResourceCollection(tables []models.Table) []TableResponse {
	resourceNotes := make([]TableResponse, len(tables))
	for i, table := range tables {
		resourceNotes[i] = TableResource(&table)
	}

	return resourceNotes
}
