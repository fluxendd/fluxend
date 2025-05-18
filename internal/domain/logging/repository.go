package logging

import (
	"fluxton/internal/api/dto"
)

type Repository interface {
	List(paginationParams dto.PaginationParams) ([]RequestLog, error)
	Create(requestLog *RequestLog) (*RequestLog, error)
}
