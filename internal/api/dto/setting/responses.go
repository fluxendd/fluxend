package setting

type Response struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Value        string `json:"value"`
	DefaultValue string `json:"defaultValue"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}
