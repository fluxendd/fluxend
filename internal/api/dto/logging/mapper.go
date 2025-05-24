package logging

import (
	"fluxton/internal/domain/logging"
	"github.com/guregu/null/v6"
)

func ToLogListInput(request *ListRequest) *logging.ListInput {
	return &logging.ListInput{
		UserUuid:  null.StringFrom(request.UserUuid),
		Status:    null.StringFrom(request.Status),
		Method:    null.StringFrom(request.Method),
		Endpoint:  null.StringFrom(request.Endpoint),
		IPAddress: null.StringFrom(request.IPAddress),
	}
}
