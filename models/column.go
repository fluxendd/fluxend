package models

type Column struct {
	Name           string `db:"column_name" json:"name"`
	Position       int    `db:"column_position" json:"position"`
	NotNull        bool   `db:"not_null" json:"notNull"`
	Type           string `db:"data_type" json:"type"`
	ConstraintType string `db:"constraint_type" json:"constraintType"`
	Primary        bool   `json:"primary"`
	Unique         bool   `json:"unique"`
	Default        string `db:"default_value" json:"defaultValue"`
}
