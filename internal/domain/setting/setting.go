package setting

import (
	"fluxton/internal/domain/shared"
	"time"
)

type Setting struct {
	shared.BaseModel
	ID           int       `db:"id"`
	Name         string    `db:"name"`
	Value        string    `db:"value"`
	DefaultValue string    `db:"default_value"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (s Setting) GetTableName() string {
	return "settings"
}
