package table

type Response struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Schema        string `json:"schema"`
	EstimatedRows int    `json:"estimatedRows"`
	TotalSize     string `json:"totalSize"`
}
