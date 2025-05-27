package health

type Health struct {
	DatabaseStatus  string `json:"database_status"`
	AppStatus       string `json:"app_status"`
	PostgrestStatus string `json:"postgrest_status"`
}
