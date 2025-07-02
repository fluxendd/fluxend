package shared

type PaginationParams struct {
	Page  int
	Limit int
	Sort  string
	Order string
}

type PaginationDetails struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}
