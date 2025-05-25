package logging

import (
	"fluxend/internal/domain/logging"
)

func ToLogListInput(request *ListRequest) *logging.ListInput {
	return &logging.ListInput{
		UserUuid:  request.UserUuid,
		Status:    request.Status,
		Method:    request.Method,
		Endpoint:  request.Endpoint,
		IPAddress: request.IPAddress,
	}
}
