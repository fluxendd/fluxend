package types

type TableColumn struct {
	Name     string `db:"column_name" json:"name"`
	Position int    `db:"column_position" json:"position"`
	NotNull  bool   `db:"not_null" json:"notNull"`
	Type     string `db:"data_type" json:"type"`
	Default  string `db:"default_value" json:"defaultValue"`
	Primary  bool   `json:"primary"`
	Unique   bool   `json:"unique"`
}
