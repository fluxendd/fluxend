package database

import (
	"github.com/guregu/null/v6"
)

type ColumnResponse struct {
	Name            string      `json:"name"`
	Position        int         `json:"position"`
	NotNull         bool        `json:"notNull"`
	Type            string      `json:"type"`
	Default         string      `json:"defaultValue"`
	Primary         bool        `json:"primary"`
	Unique          bool        `json:"unique"`
	Foreign         bool        `json:"foreign"`
	ReferenceTable  null.String `json:"referenceTable" swaggertype:"string"`
	ReferenceColumn null.String `json:"referenceColumn" swaggertype:"string"`
}
