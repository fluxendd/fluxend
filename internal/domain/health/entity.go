package health

type Health struct {
	// Service Status
	DatabaseStatus  string `json:"database_status"`
	AppStatus       string `json:"app_status"`
	PostgrestStatus string `json:"postgrest_status"`

	// System Resources
	DiskUsage     string `json:"disk_usage"`
	DiskAvailable string `json:"disk_available"`
	DiskTotal     string `json:"disk_total"`
	CPUUsage      string `json:"cpu_usage"`
	CPUCores      int    `json:"cpu_cores"`
}
