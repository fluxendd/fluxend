package logging

import (
	"fluxend/internal/domain/shared"
)

type Repository interface {
	List(input *ListInput, paginationParams shared.PaginationParams) ([]RequestLog, shared.PaginationDetails, error)
	Create(requestLog *RequestLog) (*RequestLog, error)
}
