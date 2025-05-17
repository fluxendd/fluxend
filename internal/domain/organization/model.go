package organization

import (
	"fluxton/internal/domain/shared"
	"github.com/google/uuid"
	"time"
)

type Organization struct {
	shared.BaseModel
	Uuid      uuid.UUID `db:"uuid"`
	Name      string    `db:"name"`
	CreatedBy uuid.UUID `db:"created_by"`
	UpdatedBy uuid.UUID `db:"updated_by"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u Organization) GetTableName() string {
	return "fluxton.organizations"
}
