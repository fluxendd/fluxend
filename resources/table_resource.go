package resources

import (
	"fluxton/models"
	"fluxton/types"
	"github.com/google/uuid"
)

type TableResponse struct {
	ID        uuid.UUID           `json:"id"`
	ProjectID uuid.UUID           `json:"project_id"`
	CreatedBy uuid.UUID           `json:"created_by"`
	UpdatedBy uuid.UUID           `json:"updated_by"`
	Name      string              `json:"name"`
	Columns   []types.TableColumn `json:"columns"`
	CreatedAt string              `json:"created_at"`
	UpdatedAt string              `json:"updated_at"`
}

func TableResource(table *models.Table) TableResponse {
	return TableResponse{
		ID:        table.ID,
		ProjectID: table.ProjectID,
		CreatedBy: table.CreatedBy,
		UpdatedBy: table.UpdatedBy,
		Name:      table.Name,
		Columns:   table.Columns,
		CreatedAt: table.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: table.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func TableResourceCollection(tables []models.Table) []TableResponse {
	resourceNotes := make([]TableResponse, len(tables))
	for i, table := range tables {
		resourceNotes[i] = TableResource(&table)
	}

	return resourceNotes
}
