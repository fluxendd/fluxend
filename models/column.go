package models

import (
	"database/sql"
)

type Column struct {
	Name     string `db:"name" json:"name"`
	Position int    `db:"position" json:"position"`
	NotNull  bool   `db:"not_null" json:"notNull"`
	Type     string `db:"type" json:"type"`
	Primary  bool   `db:"primary" json:"primary"`
	Unique   bool   `db:"unique" json:"unique"`
	Foreign  bool   `db:"foreign" json:"foreign"`
	Default  string `db:"default_value" json:"defaultValue"`

	// only required when constraint is FOREIGN KEY
	ReferenceTable  sql.NullString `db:"reference_table" json:"referenceTable,omitempty"`
	ReferenceColumn sql.NullString `db:"reference_column" json:"referenceColumn,omitempty"`
}
