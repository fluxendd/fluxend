package models

type Function struct {
	RoutineName       string `db:"routine_name" json:"routineName"`
	RoutineType       string `db:"routine_type" json:"routineType"`
	DataType          string `db:"data_type" json:"dataType"`
	TypeUdtName       string `db:"type_udt_name" json:"typeUdtName"`
	RoutineDefinition string `db:"routine_definition" json:"routineDefinition"`
	ExternalLanguage  string `db:"external_language" json:"externalLanguage"`
	SqlDataAccess     string `db:"sql_data_access" json:"sqlDataAccess"`
}
