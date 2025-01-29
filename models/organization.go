package models

import (
	"fmt"
	"time"
)

type Organization struct {
	ID        uint      `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u Organization) GetFields() string {
	return fmt.Sprintf("id, name, created_at, updated_at")
}
