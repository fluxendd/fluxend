package models

import (
	"github.com/guregu/null/v6"
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
	ReferenceTable  null.String `db:"reference_table" json:"referenceTable,omitempty" swaggertype:"string"`
	ReferenceColumn null.String `db:"reference_column" json:"referenceColumn,omitempty" swaggertype:"string"`
}
