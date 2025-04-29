package models

type Column struct {
	Name     string `json:"name"`
	Position int    `json:"position"`
	NotNull  bool   `json:"notNull"`
	Type     string `json:"type"`
	Primary  bool   `json:"primary"`
	Unique   bool   `json:"unique"`
	Foreign  bool   `json:"foreign"`
	Default  string `json:"defaultValue"`

	// only required when constraint is FOREIGN KEY
	ReferenceTable  string `json:"referenceTable,omitempty"`
	ReferenceColumn string `json:"referenceColumn,omitempty"`
}
