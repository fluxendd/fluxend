package logging

import (
	"fluxton/internal/domain/shared"
)

type Repository interface {
	List(input *ListInput, paginationParams shared.PaginationParams) ([]RequestLog, error)
	Create(requestLog *RequestLog) (*RequestLog, error)
}
