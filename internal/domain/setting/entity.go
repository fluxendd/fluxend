package setting

import (
	"fluxend/internal/domain/shared"
	"time"
)

type Setting struct {
	shared.BaseEntity
	ID           int       `db:"id"`
	Name         string    `db:"name"`
	Value        string    `db:"value"`
	DefaultValue string    `db:"default_value"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
