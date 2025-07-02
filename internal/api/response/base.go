package response

type Response struct {
	Success  bool        `json:"success"`
	Errors   []string    `json:"errors,omitempty"`
	Content  interface{} `json:"content,omitempty"`
	Metadata interface{} `json:"metadata,omitempty"`
}

type BadRequestErrorResponse struct {
	Success bool     `json:"success" example:"false"`
	Errors  []string `json:"errors" example:"Invalid input parameter"`
	Content *string  `json:"content" example:"null"`
}

type UnauthorizedErrorResponse struct {
	Success bool     `json:"success" example:"false"`
	Errors  []string `json:"errors" example:"Unauthorized access"`
	Content *string  `json:"content" example:"null"`
}

type UnprocessableErrorResponse struct {
	Success bool     `json:"success" example:"false"`
	Errors  []string `json:"errors" example:"Validation failed"`
	Content *string  `json:"content" example:"null"`
}

type InternalServerErrorResponse struct {
	Success bool     `json:"success" example:"false"`
	Errors  []string `json:"errors" example:"Internal server error"`
	Content *string  `json:"content" example:"null"`
}

type NotFoundErrorResponse struct {
	Success bool     `json:"success" example:"false"`
	Errors  []string `json:"errors" example:"Resource not found"`
	Content *string  `json:"content" example:"null"`
}

type ForbiddenErrorResponse struct {
	Success bool     `json:"success" example:"false"`
	Errors  []string `json:"errors" example:"Forbidden access"`
	Content *string  `json:"content" example:"null"`
}
