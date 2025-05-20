package response

type Response struct {
	Success bool        `json:"success"`
	Errors  []string    `json:"errors"`
	Content interface{} `json:"content"`
}
