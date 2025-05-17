package function

type Function struct {
	Name          string `db:"routine_name" json:"name"`
	Type          string `db:"routine_type" json:"type"`
	DataType      string `db:"data_type" json:"dataType"`
	TypeUdtName   string `db:"type_udt_name" json:"typeUdtName"`
	Definition    string `db:"routine_definition" json:"definition"`
	Language      string `db:"external_language" json:"language"`
	SqlDataAccess string `db:"sql_data_access" json:"sqlDataAccess"`
}
